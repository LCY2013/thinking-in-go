package main

import (
	"log"
	"net"
	"time"
)

/*
1、向服务端建立 TCP 连接

一旦服务端按照上面的Listen + Accept结构成功启动，客户端便可以使用net.Dial或net.DialTimeout向服务端发起连接建立的请求。
conn, err := net.Dial("tcp", "localhost:8888")
conn, err := net.DialTimeout("tcp", "localhost:8888", 2 * time.Second)

Dial 函数向服务端发起 TCP 连接，这个函数会一直阻塞，直到连接成功或失败后，才会返回。
而 DialTimeout 带有超时机制，如果连接耗时大于超时时间，这个函数会返回超时错误。
对于客户端来说，连接的建立还可能会遇到几种特殊情形。

第一种情况：网络不可达或对方服务未启动。
如果传给Dial的服务端地址是网络不可达的，或者服务地址中端口对应的服务并没有启动，端口未被监听（Listen），Dial几乎会立即返回类似这样的错误。

第二种情况：对方服务的 listen backlog 队列满。
当对方服务器很忙，瞬间有大量客户端尝试向服务端建立连接时，服务端可能会出现 listen backlog 队列满，接收连接（accept）不及时的情况，这就会导致客户端的Dial调用阻塞，直到服务端进行一次 accept，从 backlog 队列中腾出一个槽位，客户端的 Dial 才会返回成功。
而且，不同操作系统下 backlog 队列的长度是不同的，在 macOS 下，这个默认值如下：
$sysctl -a|grep kern.ipc.somaxconn
kern.ipc.somaxconn: 128

在 Ubuntu Linux 下，backlog 队列的长度值与系统中net.ipv4.tcp_max_syn_backlog的设置有关。
那么，极端情况下，如果服务端一直不执行accept操作，那么客户端会一直阻塞吗？
答案是不会！
如果服务端运行在 macOS 下，那么客户端会阻塞大约 1 分多钟，才会返回超时错误：
而如果服务端运行在 Ubuntu 上，客户端的Dial调用大约在 2 分多钟后提示超时错误，这个结果也和 Linux 的系统设置有关。

第三种情况：若网络延迟较大，Dial 将阻塞并超时。
如果网络延迟较大，TCP 连接的建立过程（三次握手）将更加艰难坎坷，会经历各种丢包，时间消耗自然也会更长，这种情况下，Dial函数会阻塞。
如果经过长时间阻塞后依旧无法建立连接，那么Dial也会返回类似getsockopt: operation timed out的错误。

在连接建立阶段，多数情况下Dial是可以满足需求的，即便是阻塞一小会儿也没事。
但对于那些需要有严格的连接时间限定的 Go 应用，如果一定时间内没能成功建立连接，程序可能会需要执行一段“错误”处理逻辑，所以，这种情况下使用DialTimeout函数更适合。

2、Socket 读操作
连接建立起来后，就要在连接上进行读写以完成业务逻辑。
Go 运行时隐藏了 I/O 多路复用的复杂性。Go 语言使用者只需采用 Goroutine+ 阻塞 I/O 模型，就可以满足大部分场景需求。
Dial 连接成功后，会返回一个 net.Conn 接口类型的变量值，这个接口变量的底层类型为一个 *TCPConn：
//$GOROOT/src/net/tcpsock.go
type TCPConn struct {
    conn
}
TCPConn 内嵌了一个非导出类型：conn（封装了底层的 socket），因此，TCPConn“继承”了conn类型的Read和Write方法，后续通过Dial函数返回值调用的Read和Write方法都是 net.conn 的方法，它们分别代表了对 socket 的读和写。

首先是 Socket 中无数据的场景。
连接建立后，如果客户端未发送数据，服务端会阻塞在 Socket 的读操作上，这种“阻塞 I/O 模型”的行为模式是一致的。
执行该这个操作的 Goroutine 也会被挂起。
Go 运行时会监视这个 Socket，直到它有数据读事件，才会重新调度这个 Socket 对应的 Goroutine 完成读操作。

第二种情况是 Socket 中有部分数据。
如果 Socket 中有部分数据就绪，且数据数量小于一次读操作期望读出的数据长度，那么读操作将会成功读出这部分数据，并返回，而不是等待期望长度数据全部读取后，再返回。

举个例子，服务端创建一个长度为 10 的切片作为接收数据的缓冲区，等待 Read 操作将读取的数据放入切片。
当客户端在已经建立成功的连接上，成功写入两个字节的数据（比如：hi）后，服务端的 Read 方法将成功读取数据，并返回n=2，err=nil，而不是等收满 10 个字节后才返回。

第三种情况是 Socket 中有足够数据。
如果连接上有数据，且数据长度大于等于一次Read操作期望读出的数据长度，那么Read将会成功读出这部分数据，并返回。这个情景是最符合我们对Read的期待的了。

我们以上面的例子为例，当客户端在已经建立成功的连接上，成功写入 15 个字节的数据后，服务端进行第一次Read时，会用连接上的数据将我们传入的切片缓冲区（长度为 10）填满后返回：n = 10, err = nil。这个时候，内核缓冲区中还剩 5 个字节数据，当服务端再次调用Read方法时，就会把剩余数据全部读出。

最后一种情况是设置读操作超时。
有些场合，对 socket 的读操作的阻塞时间有严格限制的，但由于 Go 使用的是阻塞 I/O 模型，如果没有可读数据，Read 操作会一直阻塞在对 Socket 的读操作上。

这时，我们可以通过 net.Conn 提供的 SetReadDeadline 方法，设置读操作的超时时间，当超时后仍然没有数据可读的情况下，Read 操作会解除阻塞并返回超时错误，这就给 Read 方法的调用者提供了进行其他业务处理逻辑的机会。

SetReadDeadline 方法接受一个绝对时间作为超时的 deadline，一旦通过这个方法设置了某个 socket 的 Read deadline，那么无论后续的 Read 操作是否超时，只要我们不重新设置 Deadline，那么后面与这个 socket 有关的所有读操作，都会返回超时失败错误。
func handleConn(c net.Conn) {
    defer c.Close()
    for {
        // read from the connection
        var buf = make([]byte, 128)
        c.SetReadDeadline(time.Now().Add(time.Second))
        n, err := c.Read(buf)
        if err != nil {
            log.Printf("conn read %d bytes,  error: %s", n, err)
            if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
                // 进行其他业务逻辑的处理
                continue
            }
            return
        }
        log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
    }
}

3、Socket 写操作
通过 net.Conn 实例的 Write 方法，可以将数据写入 Socket。当 Write 调用的返回值 n 的值，与预期要写入的数据长度相等，且 err = nil 时，就执行了一次成功的 Socket 写操作，这是在调用 Write 时遇到的最常见的情形。

和 Socket 的读操作一些特殊情形相比，Socket 写操作遇到的特殊情形同样不少，也逐一看一下。

第一种情况：写阻塞。
TCP 协议通信两方的操作系统内核，都会为这个连接保留数据缓冲区，调用 Write 向 Socket 写入数据，实际上是将数据写入到操作系统协议栈的数据缓冲区中。
TCP 是全双工通信，因此每个方向都有独立的数据缓冲。
当发送方将对方的接收缓冲区，以及自身的发送缓冲区都写满后，再调用 Write 方法就会出现阻塞的情况。

来看一个具体例子，这个例子的客户端代码如下：
action/16_action/network/v01/client/write/main.go

客户端每次调用 Write 方法向服务端写入 65536 个字节，并在 Write 方法返回后，输出此次 Write 的写入字节数和程序启动后写入的总字节数量。

服务端的处理程序逻辑，也摘录了主要部分，可以看一下：
action/16_action/network/v01/server/write/write.go

第二种情况：写入部分数据。
Write 操作存在写入部分数据的情况，比如上面例子中，当客户端输出日志停留在“write 65536 bytes this time, 655360 bytes in total”时，杀掉服务端。

显然，Write并不是在 655360 这个地方阻塞的，而是后续又写入 24108 个字节后发生了阻塞，服务端 Socket 关闭后，看到客户端又写入 24108 字节后，才返回的broken pipe错误。
由于这 24108 字节数据并未真正被服务端接收到，程序需要考虑妥善处理这些数据，以防数据丢失。

第三种情况：写入超时。
如果非要给 Write 操作增加一个期限，可以调用 SetWriteDeadline 方法。
比如，可以将上面例子中的客户端源码拷贝一份，然后在新客户端源码中的 Write 调用之前，增加一行超时时间设置代码：
conn.SetWriteDeadline(time.Now().Add(time.Microsecond * 10))

可以看到，在 Write 方法写入超时时，依旧存在数据部分写入（仅写入 24108 个字节）的情况。
另外，和 SetReadDeadline 一样，只要我们通过 SetWriteDeadline 设置了写超时，那无论后续 Write 方法是否成功，
如果不重新设置写超时或取消写超时，后续对 Socket 的写操作都将以超时失败告终。

综合上面这些例子，虽然 Go 提供了阻塞 I/O 的便利，但在调用Read和Write时，依旧要综合函数返回的n和err的结果以做出正确处理。

不过，前面说的 Socket 读与写都是限于单 Goroutine 下的操作，如果多个 Goroutine 并发读或写一个 socket 会发生什么呢？

4、并发 Socket 读写
Goroutine 的网络编程模型，决定了存在着不同 Goroutine 间共享conn的情况，那么conn的读写是否是 Goroutine 并发安全的呢？
不过，在深入这个问题之前，先从应用的角度上，看看并发 read 操作和 write 操作的 Goroutine 安全的必要性。

对于 Read 操作而言，由于 TCP 是面向字节流，conn.Read无法正确区分数据的业务边界，因此，多个 Goroutine 对同一个 conn 进行 read 的意义不大，Goroutine 读到不完整的业务包，反倒增加了业务处理的难度。

但对与 Write 操作而言，倒是有多个 Goroutine 并发写的情况。
不过 conn 读写是否是 Goroutine 安全的测试并不是很好做，先深入一下运行时代码，从理论上给这个问题定个性。

首先，net.conn只是*netFD 的外层包裹结构，最终 Write 和 Read 都会落在其中的fd字段上：
//$GOROOT/src/net/net.go
type conn struct {
    fd *netFD
}

另外，netFD 在不同平台上有着不同的实现，以net/fd_posix.go中的netFD为例看看：
// $GOROOT/src/net/fd_unix.go
// Network file descriptor.
type netFD struct {
    pfd poll.FD

    // immutable until Close
    family      int
    sotype      int
    isConnected bool // handshake completed or use of association with peer
    net         string
    laddr       Addr
    raddr       Addr
}

netFD 中最重要的字段是 poll.FD 类型的 pfd，它用于表示一个网络连接。我也把它的结构摘录了一部分：
// $GOROOT/src/internal/poll/fd_unix.go
// FD is a file descriptor. The net and os packages use this type as a
// field of a larger type representing a network connection or OS file.
type FD struct {
    // Lock sysfd and serialize access to Read and Write methods.
    fdmu fdMutex

    // System file descriptor. Immutable until Close.
    Sysfd int

    // I/O poller.
    pd pollDesc
    // Writev cache.
    iovecs *[]syscall.Iovec
    ... ...
}

FD类型中包含了一个运行时实现的fdMutex类型字段。
从它的注释来看，这个fdMutex用来串行化对字段Sysfd的 Write 和 Read 操作。
也就是说，所有对这个 FD 所代表的连接的 Read 和 Write 操作，都是由fdMutex来同步的。
从FD的 Read 和 Write 方法的实现，也证实了这一点：
// $GOROOT/src/internal/poll/fd_unix.go
func (fd *FD) Read(p []byte) (int, error) {
    if err := fd.readLock(); err != nil {
        return 0, err
    }
    defer fd.readUnlock()
    if len(p) == 0 {
        // If the caller wanted a zero byte read, return immediately
        // without trying (but after acquiring the readLock).
        // Otherwise syscall.Read returns 0, nil which looks like
        // io.EOF.
        // TODO(bradfitz): make it wait for readability? (Issue 15735)
        return 0, nil
    }
    if err := fd.pd.prepareRead(fd.isFile); err != nil {
        return 0, err
    }
    if fd.IsStream && len(p) > maxRW {
        p = p[:maxRW]
    }
    for {
        n, err := ignoringEINTRIO(syscall.Read, fd.Sysfd, p)
        if err != nil {
            n = 0
            if err == syscall.EAGAIN && fd.pd.pollable() {
                if err = fd.pd.waitRead(fd.isFile); err == nil {
                    continue
                }
            }
        }
        err = fd.eofError(n, err)
        return n, err
    }
}
func (fd *FD) Write(p []byte) (int, error) {
    if err := fd.writeLock(); err != nil {
        return 0, err
    }
    defer fd.writeUnlock()
    if err := fd.pd.prepareWrite(fd.isFile); err != nil {
        return 0, err
    }
    var nn int
    for {
        max := len(p)
        if fd.IsStream && max-nn > maxRW {
            max = nn + maxRW
        }
        n, err := ignoringEINTRIO(syscall.Write, fd.Sysfd, p[nn:max])
        if n > 0 {
            nn += n
        }
        if nn == len(p) {
            return nn, err
        }
        if err == syscall.EAGAIN && fd.pd.pollable() {
            if err = fd.pd.waitWrite(fd.isFile); err == nil {
                continue
            }
        }
        if err != nil {
            return nn, err
        }
        if n == 0 {
            return nn, io.ErrUnexpectedEOF
        }
    }
}

你看，每次 Write 操作都是受 lock 保护，直到这次数据全部写完才会解锁。
因此，在应用层面，要想保证多个 Goroutine 在一个conn上 write 操作是安全的，需要一次 write 操作完整地写入一个“业务包”。
一旦将业务包的写入拆分为多次 write，那也无法保证某个 Goroutine 的某“业务包”数据在conn发送的连续性。

同时，我们也可以看出即便是 Read 操作，也是有 lock 保护的。
多个 Goroutine 对同一conn的并发读，不会出现读出内容重叠的情况，但就像前面讲并发读的必要性时说的那样，一旦采用了不恰当长度的切片作为 buf，很可能读出不完整的业务包，这反倒会带来业务上的处理难度。

比如一个完整数据包：world，当 Goroutine 的读缓冲区长度 < 5 时，就存在这样一种可能：一个 Goroutine 读出了“worl”，而另外一个 Goroutine 读出了"d"。

5、Socket 关闭
通常情况下，当客户端需要断开与服务端的连接时，客户端会调用 net.Conn 的 Close 方法关闭与服务端通信的 Socket。
如果客户端主动关闭了 Socket，那么服务端的Read调用将会读到什么呢？这里要分“有数据关闭”和“无数据关闭”两种情况。

“有数据关闭”是指在客户端关闭连接（Socket）时，Socket 中还有服务端尚未读取的数据。
在这种情况下，服务端的 Read 会成功将剩余数据读取出来，最后一次 Read 操作将得到io.EOF错误码，表示客户端已经断开了连接。
如果是在“无数据关闭”情形下，服务端调用的 Read 方法将直接返回io.EOF。

不过因为 Socket 是全双工的，客户端关闭 Socket 后，如果服务端 Socket 尚未关闭，这个时候服务端向 Socket 的写入操作依然可能会成功，因为数据会成功写入己方的内核 socket 缓冲区中，即便最终发不到对方 socket 缓冲区也会这样。
因此，当发现对方 socket 关闭后，己方应该正确合理处理自己的 socket，再继续 write 已经没有任何意义了。







*/

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// read from the connection
		var buf = make([]byte, 128)
		err := c.SetReadDeadline(time.Now().Add(time.Second))
		if err != nil {
			log.Printf("set read deadline error: %+v", err)
		}
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes,  error: %s", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				// 进行其他业务逻辑的处理
				continue
			}
			return
		}
		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
}

func main() {
	//conn, err := net.Dial("tcp", "localhost:8888")
}
