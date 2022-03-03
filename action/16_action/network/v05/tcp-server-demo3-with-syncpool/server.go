package tcp_server_demo3

/*
1、建立对协议的抽象
程序是对现实世界的抽象。对于现实世界的自定义应用协议规范，需要在程序世界建立起对这份协议的抽象。
在进行抽象之前，先建立这次实现要用的源码项目 tcp-server-demo1，建立的步骤如下：
$mkdir tcp-server-demo1
$cd tcp-server-demo1
$go mod init github.com/bigwhite/tcp-server-demo1

action/16_action/network/v02/tcp-server-demo1/自定义协议规范.png

1）深入协议字段
这是一个高度简化的、基于二进制模式定义的协议。二进制模式定义的特点，就是采用长度字段标识独立数据包的边界。

在这个协议规范中看到：请求包和应答包的第一个字段（totalLength）都是包的总长度，它就是用来标识包边界的那个字段，也是在应用层用于“分割包”的最重要字段。

请求包与应答包的第二个字段也一样，都是 commandID，这个字段用于标识包类型，这里我们定义四种包类型：
连接请求包（值为 0x01）
消息请求包（值为 0x02）
连接响应包（值为 0x81）
消息响应包（值为 0x82）
换为对应的代码就是：
const (
    CommandConn   = iota + 0x01 // 0x01，连接请求包
    CommandSubmit               // 0x02，消息请求包
)
const (
    CommandConnAck   = iota + 0x80 // 0x81，连接请求的响应包
    CommandSubmitAck               // 0x82，消息请求的响应包
)

请求包与应答包的第三个字段都是 ID，ID 是每个连接上请求包的消息流水号，顺序累加，步长为 1，循环使用，多用来请求发送方后匹配响应包，所以要求一对请求与响应消息的流水号必须相同。

请求包与响应包唯一的不同之处，就在于最后一个字段：请求包定义了有效载荷（payload），这个字段承载了应用层需要的业务数据；而响应包则定义了请求包的响应状态字段（result），这里其实简化了响应状态字段的取值，成功的响应用 0 表示，如果是失败的响应，无论失败原因是什么，我们都用 1 来表示。

明确了应用层协议的各个字段定义之后，接下来就看看如何建立起对这个协议的抽象。

2）建立 Frame 和 Packet 抽象
首先要知道，TCP 连接上的数据是一个没有边界的字节流，但在业务层眼中，没有字节流，只有各种协议消息。
因此，无论是从客户端到服务端，还是从服务端到客户端，业务层在连接上看到的都应该是一个挨着一个的协议消息流。

现在建立第一个抽象：Frame。
每个 Frame 表示一个协议消息，这样在业务层眼中，连接上的字节流就是由一个接着一个 Frame 组成的，如下图所示：
action/16_action/network/v02/tcp-server-demo1/Frame表示一个协议消息.png

自定义协议就封装在这一个个的 Frame 中。
协议规定了将 Frame 分割开来的方法，那就是利用每个 Frame 开始处的 totalLength，每个 Frame 由一个 totalLength 和 Frame 的负载（payload）构成，
比如可以看看下图中左侧的 Frame 结构：
action/16_action/network/v02/tcp-server-demo1/Frame表示一个协议消息.png

这样，通过 Frame header: totalLength 就可以将 Frame 之间隔离开来。

在这个基础上，建立协议的第二个抽象：Packet。将 Frame payload 定义为一个 Packet。上图右侧展示的就是 Packet 的结构。

Packet 就是业务层真正需要的消息，每个 Packet 由 Packet 头和 Packet Body 部分组成。
Packet 头就是 commandID，用于标识这个消息的类型；
而 ID 和 payload（packet payload）或 result 字段组成了 Packet 的 Body 部分，对业务层有价值的数据都包含在 Packet Body 部分。

那么到这里，通过 Frame 和 Packet 两个类型结构，完成了程序世界对私有协议规范的抽象。
接下来，要做的就是基于 Frame 和 Packet 这两个概念，实现对私有协议的解包与打包操作。

3）协议的解包与打包
所谓协议的解包（decode），就是指识别 TCP 连接上的字节流，将一组字节“转换”成一个特定类型的协议消息结构，然后这个消息结构会被业务处理逻辑使用。

而打包（encode）刚刚好相反，是指将一个特定类型的消息结构转换为一组字节，然后这组字节数据会被放在连接上发送出去。

具体到这个自定义协议上，解包就是指字节流 -> Frame，打包是指Frame -> 字节流。
可以看一下针对这个协议的服务端解包与打包的流程图：
action/16_action/network/v02/tcp-server-demo1/服务端解包与打包的流程图.png

TCP 流数据先后经过 frame decode 和 packet decode，得到应用层所需的 packet 数据，而业务层回复的响应，则先后经过 packet 的 encode 与 frame 的 encode，写入 TCP 数据流中。

4）Frame 的实现
协议部分最重要的两个抽象是 Frame 和 Packet，于是就在项目中建立 frame 包与 packet 包，分别与两个协议抽象对应。
frame 包的职责是提供识别 TCP 流边界的编解码器，可以很容易为这样的编解码器，定义出一个统一的接口类型 StreamFrameCodec：
action/16_action/network/v02/tcp-server-demo1/frame/frame.go

StreamFrameCodec 接口类型有两个方法 Encode 与 Decode。Encode 方法用于将输入的 Frame payload 编码为一个 Frame，然后写入 io.Writer 所代表的输出（outbound）TCP 流中。
而 Decode 方法正好相反，它从代表输入（inbound）TCP 流的 io.Reader 中读取一个完整 Frame，并将得到的 Frame payload 解析出来并返回。

给出一个针对协议的 StreamFrameCodec 接口的实现：
action/16_action/network/v02/tcp-server-demo1/frame/inner.go

在在这段实现中，有三点事项需要注意：
网络字节序使用大端字节序（BigEndian），因此无论是 Encode 还是 Decode，都是用 binary.BigEndian；

binary.Read 或 Write 会根据参数的宽度，读取或写入对应的字节个数的字节，这里 totalLen 使用 int32，那么 Read 或 Write 只会操作数据流中的 4 个字节；

这里没有设置网络 I/O 操作的 Deadline，io.ReadFull 一般会读满所需的字节数，除非遇到 EOF 或 ErrUnexpectedEOF。

在工程实践中，保证打包与解包正确的最有效方式就是编写单元测试，StreamFrameCodec 接口的 Decode 和 Encode 方法的参数都是接口类型，这可以很容易为 StreamFrameCodec 接口的实现编写测试用例。
下面是为 innerFrameCodec 编写了两个测试用例：
action/16_action/network/v02/tcp-server-demo1/frame/inner_test.go

测试 Encode 方法，其实不需要建立真实的网络连接，只要用一个满足 io.Writer 的 bytes.Buffer 实例“冒充”真实网络连接就可以了，同时 bytes.Buffer 类型也实现了 io.Reader 接口，可以很方便地从中读取出 Encode 后的内容，并进行校验比对。

为了提升测试覆盖率，还需要尽可能让测试覆盖到所有可测的错误执行分支上。
这里，模拟了 Read 或 Write 出错的情况，让执行流进入到 Decode 或 Encode 方法的错误分支中：
action/16_action/network/v02/tcp-server-demo1/frame/codec_test.go

为了实现错误分支的测试，在测试代码源文件中创建了两个类型：ReturnErrorWriter 和 ReturnErrorReader，它们分别实现了 io.Writer 与 io.Reader。

可以控制在第几次调用这两个类型的 Write 或 Read 方法时，返回错误，这样就可以让 Encode 或 Decode 方法按照我们的意图，进入到不同错误分支中去。
有了这两个用例，frame 包的测试覆盖率（通过 go test -cover . 可以查看）就可以达到 90% 以上了。

5）Packet 的实现
接下来，再看看 Packet 这个抽象的实现。
和 Frame 不同，Packet 有多种类型（这里只定义了 Conn、submit、connack、submit ack)。
所以要先抽象一下这些类型需要遵循的共同接口：
// tcp-server-demo1/packet/packet.go
type Packet interface {
    Decode([]byte) error     // []byte -> struct
    Encode() ([]byte, error) //  struct -> []byte
}

其中，Decode 是将一段字节流数据解码为一个 Packet 类型，可能是 conn，可能是 submit 等，具体要根据解码出来的 commandID 判断。而 Encode 则是将一个 Packet 类型编码为一段字节流数据。

这里只完成 submit 和 submitack 类型的 Packet 接口实现，省略了 conn 流程，也省略 conn 以及 connack 类型的实现。
action/16_action/network/v02/tcp-server-demo1/packet/submit.go

这里各种类型的编解码被调用的前提，是明确数据流是什么类型的，因此需要在包级提供一个导出的函数 Decode，这个函数负责从字节流中解析出对应的类型（根据 commandID），并调用对应类型的 Decode 方法。
action/16_action/network/v02/tcp-server-demo1/packet/packet.go:39

同样，也需要包级的 Encode 函数，根据传入的 packet 类型调用对应的 Encode 方法实现对象的编码。
action/16_action/network/v02/tcp-server-demo1/packet/packet.go:67

不过，对 packet 包中各个类型的 Encode 和 Decode 方法的测试，与 frame 包的相似。

好了，万事俱备，只欠东风！下面就来编写服务端的程序结构，将 tcp conn 与 Frame、Packet 连接起来。

5）服务端的组装
按照每个连接一个 Goroutine 的模型，给出了典型 Go 网络服务端程序的结构，这里就以这个结构为基础，将 Frame、Packet 加进来，形成第一版服务端实现：
action/16_action/network/v02/tcp-server-demo1/cmd/server/main.go

这个程序的逻辑非常清晰，服务端程序监听 8888 端口，并在每次调用 Accept 方法后得到一个新连接，服务端程序将这个新连接交到一个新的 Goroutine 中处理。

新 Goroutine 的主函数为 handleConn，有了 Packet 和 Frame 这两个抽象的加持，这个函数同样拥有清晰的代码调用结构：
// handleConn的调用结构
read frame from conn
    ->frame decode
      -> handle packet
        -> packet decode
        -> packet(ack) encode
    ->frame(ack) encode
write ack frame to conn

一个基于 TCP 的自定义应用层协议的经典阻塞式的服务端就完成了。
不过这里的服务端依旧是一个简化的实现，比如这里没有考虑支持优雅退出、没有捕捉某个链接上出现的可能导致整个程序退出的 panic 等。

6）验证测试
要验证服务端的实现是否可以正常工作，需要实现一个自定义应用层协议的客户端。
这里，同样基于 frame、packet 两个包，实现了一个自定义应用层协议的客户端。
下面是客户端的 main 函数：
action/16_action/network/v02/tcp-server-demo1/cmd/client/main.go

关于 startClient 函数，需要简单说明几点。
首先，startClient 函数启动了两个 Goroutine，一个负责向服务端发送 submit 消息请求，另外一个 Goroutine 则负责读取服务端返回的响应；
其次，客户端发送的 submit 请求的负载（payload）是由第三方包 github.com/lucasepe/codename 负责生成的，这个包会生成一些对人类可读的随机字符串，比如：firm-iron、 moving-colleen、game-nova 这样的字符串；
另外，负责读取服务端返回响应的 Goroutine，使用 SetReadDeadline 方法设置了读超时，这主要是考虑该 Goroutine 可以在收到退出通知时，能及时从 Read 阻塞中跳出来。
好了，现在来构建和运行一下这两个程序。
在 tcp-server-demo1 目录下提供了 Makefile，如果你使用的是 Linux 或 macOS 操作系统，可以直接敲入 make 构建两个程序，如果你是在 Windows 下构建，可以直接敲入下面的 go build 命令构建：
$make
go build github.com/bigwhite/tcp-server-demo1/cmd/server
go build github.com/bigwhite/tcp-server-demo1/cmd/client


*/
