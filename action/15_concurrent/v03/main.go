package main

/*
channel中蕴含大智慧

Go 语言实现了基于 CSP（Communicating Sequential Processes）理论的并发方案。

Go 语言的 CSP 模型的实现包含两个主要组成部分：
一个是 Goroutine，它是 Go 应用并发设计的基本构建与执行单元；
另一个就是 channel，它在并发模型中扮演着重要的角色。
channel 既可以用来实现 Goroutine 间的通信，还可以实现 Goroutine 间的同步。
它就好比 Go 并发设计这门“武功”的秘籍口诀，可以说，学会在 Go 并发设计时灵活运用 channel，才能说真正掌握了 Go 并发设计的真谛。

1、作为一等公民的 channel

channel 作为一等公民意味着什么呢？

可以像使用普通变量那样使用 channel，比如，定义 channel 类型变量、给 channel 变量赋值、将 channel 作为参数传递给函数 / 方法、将 channel 作为返回值从函数 / 方法中返回，甚至将 channel 发送到其他 channel 中。
这就大大简化了 channel 原语的使用，提升了我们开发者在做并发设计和实现时的体验。

1）创建 channel
和切片、结构体、map 等一样，channel 也是一种复合数据类型。
也就是说，在声明一个 channel 类型变量时，必须给出其具体的元素类型，比如下面的代码这样：
var ch chan int  // 声明了一个元素为 int 类型的 channel 类型变量 ch。

如果 channel 类型变量在声明时没有被赋予初值，那么它的默认值为 nil。
并且，和其他复合数据类型支持使用复合类型字面值作为变量初始值不同，为 channel 类型变量赋初值的唯一方法就是使用 make 这个 Go 预定义的函数，比如下面代码：
ch1 := make(chan int)
ch2 := make(chan int, 5)

第一行我们通过make(chan T)创建的、元素类型为 T 的 channel 类型，是无缓冲 channel，
而第二行中通过带有 capacity 参数的make(chan T, capacity)创建的元素类型为 T、缓冲区长度为 capacity 的 channel 类型，是带缓冲 channel。

2）发送与接收
Go 提供了<-操作符用于对 channel 类型变量进行发送与接收操作：
ch1 <- 13    // 将整型字面值13发送到无缓冲channel类型变量ch1中
n := <- ch1  // 从无缓冲channel类型变量ch1中接收一个整型值存储到整型变量n中
ch2 <- 17    // 将整型字面值17发送到带缓冲channel类型变量ch2中
m := <- ch2  // 从带缓冲channel类型变量ch2中接收一个整型值存储到整型变量m中

在理解 channel 的发送与接收操作时，你一定要始终牢记：channel 是用于 Goroutine 间通信的，所以绝大多数对 channel 的读写都被分别放在了不同的 Goroutine 中。

由于无缓冲 channel 的运行时层实现不带有缓冲区，所以 Goroutine 对无缓冲 channel 的接收和发送操作是同步的。
也就是说，对同一个无缓冲 channel，只有对它进行接收操作的 Goroutine 和对它进行发送操作的 Goroutine 都存在的情况下，通信才能得以进行，否则单方面的操作会让对应的 Goroutine 陷入挂起状态，比如下面示例代码：
func main() {
    ch1 := make(chan int)
    ch1 <- 13 // fatal error: all goroutines are asleep - deadlock!
    n := <-ch1
    println(n)
}

运行这个示例，我们就会得到 fatal error，提示我们所有 Goroutine 都处于休眠状态，程序处于死锁状态。
要想解除这种错误状态，只需要将接收操作，或者发送操作放到另外一个 Goroutine 中就可以了，比如下面代码：
func main() {
    ch1 := make(chan int)
    go func() {
        ch1 <- 13 // 将发送操作放入一个新goroutine中执行
    }()
    n := <-ch1
    println(n)
}

由此，可以得出结论：对无缓冲 channel 类型的发送与接收操作，一定要放在两个不同的 Goroutine 中进行，否则会导致 deadlock。

和无缓冲 channel 相反，带缓冲 channel 的运行时层实现带有缓冲区，
因此，对带缓冲 channel 的发送操作在缓冲区未满、接收操作在缓冲区非空的情况下是异步的（发送或接收不需要阻塞等待）。

也就是说，对一个带缓冲 channel 来说，在缓冲区未满的情况下，对它进行发送操作的 Goroutine 并不会阻塞挂起；
在缓冲区有数据的情况下，对它进行接收操作的 Goroutine 也不会阻塞挂起。

ch2 := make(chan int, 1)
n := <-ch2 // 由于此时ch2的缓冲区中无数据，因此对其进行接收操作将导致goroutine挂起

ch3 := make(chan int, 1)
ch3 <- 17  // 向ch3发送一个整型数17
ch3 <- 27  // 由于此时ch3中缓冲区已满，再向ch3发送数据也将导致goroutine挂起

使用操作符<-，还可以声明只发送 channel 类型（send-only）和只接收 channel 类型（recv-only），我们接着看下面这个例子：
ch1 := make(chan<- int, 1) // 只发送channel类型
ch2 := make(<-chan int, 1) // 只接收channel类型

<-ch1       // invalid operation: <-ch1 (receive from send-only type chan<- int)
ch2 <- 13   // invalid operation: ch2 <- 13 (send to receive-only type <-chan int)

试图从一个只发送 channel 类型变量中接收数据，或者向一个只接收 channel 类型发送数据，都会导致编译错误。
通常只发送 channel 类型和只接收 channel 类型，会被用作函数的参数类型或返回值，用于限制对 channel 内的操作，或者是明确可对 channel 进行的操作的类型，比如下面这个例子：
func produce(ch chan<- int) {
    for i := 0; i < 10; i++ {
        ch <- i + 1
        time.Sleep(time.Second)
    }
    close(ch)
}
func consume(ch <-chan int) {
    for n := range ch {
        println(n)
    }
}
func main() {
    ch := make(chan int, 5)
    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        produce(ch)
        wg.Done()
    }()
    go func() {
        consume(ch)
        wg.Done()
    }()
    wg.Wait()
}

启动了两个 Goroutine，分别代表生产者（produce）与消费者（consume）。
生产者只能向 channel 中发送数据，我们使用chan<- int作为 produce 函数的参数类型；
消费者只能从 channel 中接收数据，我们使用<-chan int作为 consume 函数的参数类型。

在消费者函数 consume 中，使用了 for range 循环语句来从 channel 中接收数据，for range 会阻塞在对 channel 的接收操作上，直到 channel 中有数据可接收或 channel 被关闭循环，才会继续向下执行。channel 被关闭后，for range 循环也就结束了。

3）关闭 channel
produce 函数在发送完数据后，调用 Go 内置的 close 函数关闭了 channel。channel 关闭后，所有等待从这个 channel 接收数据的操作都将返回。

采用不同接收语法形式的语句，在 channel 被关闭后的返回值的情况：
n := <- ch      // 当ch被关闭后，n将被赋值为ch元素类型的零值
m, ok := <-ch   // 当ch被关闭后，m将被赋值为ch元素类型的零值, ok值为false
for v := range ch { // 当ch被关闭后，for range循环结束
    ... ...
}

通过“comma, ok”惯用法或 for range 语句，可以准确地判定 channel 是否被关闭。
而单纯采用n := <-ch形式的语句，就无法判定从 ch 返回的元素类型零值，究竟是不是因为 channel 被关闭后才返回的。

另外，从前面 produce 的示例程序中，也可以看到，channel 是在 produce 函数中被关闭的，这也是 channel 的一个使用惯例，那就是发送端负责关闭 channel。

这里为什么要在发送端关闭 channel 呢？

这是因为发送端没有像接受端那样的、可以安全判断 channel 是否被关闭了的方法。同时，一旦向一个已经关闭的 channel 执行发送操作，这个操作就会引发 panic，比如下面这个示例：
ch := make(chan int, 5)
close(ch)
ch <- 13 // panic: send on closed channel

4）select
当涉及同时对多个 channel 进行操作时，结合 Go 为 CSP 并发模型提供的另外一个原语 select，一起使用。

通过 select，可以同时在多个 channel 上进行发送 / 接收操作：
select {
case x := <-ch1:     // 从channel ch1接收数据
  ... ...
case y, ok := <-ch2: // 从channel ch2接收数据，并根据ok值判断ch2是否已经关闭
  ... ...
case ch3 <- z:       // 将z值发送到channel ch3中:
  ... ...
default:             // 当上面case中的channel通信均无法实施时，执行该默认分支
}

当 select 语句中没有 default 分支，而且所有 case 中的 channel 操作都阻塞了的时候，整个 select 语句都将被阻塞，直到某一个 case 上的 channel 变成可发送，或者某个 case 上的 channel 变成可接收，select 语句才可以继续进行下去。

2、无缓冲 channel 的惯用法
无缓冲 channel 兼具通信和同步特性，在并发程序中应用颇为广泛。

来看看几个无缓冲 channel 的典型应用：
1）第一种用法：用作信号传递

无缓冲 channel 用作信号传递的时候，有两种情况，分别是 1 对 1 通知信号和 1 对 n 通知信号。

先来分析下 1 对 1 通知信号这种情况。

type signal struct{}
func worker() {
    println("worker is working...")
    time.Sleep(1 * time.Second)
}
func spawn(f func()) <-chan signal {
    c := make(chan signal)
    go func() {
        println("worker start to work...")
        f()
        c <- signal(struct{}{})
    }()
    return c
}
func main() {
    println("start a worker...")
    c := spawn(worker)
    <-c
    fmt.Println("worker work done!")
}

spawn 函数返回的 channel，被用于承载新 Goroutine 退出的“通知信号”，这个信号专门用作通知 main goroutine。
main goroutine 在调用 spawn 函数后一直阻塞在对这个“通知信号”的接收动作上。

有些时候，无缓冲 channel 还被用来实现 1 对 n 的信号通知机制。
这样的信号通知机制，常被用于协调多个 Goroutine 一起工作，比如下面的例子：
func worker(i int) {
    fmt.Printf("worker %d: is working...\n", i)
    time.Sleep(1 * time.Second)
    fmt.Printf("worker %d: works done\n", i)
}
func spawnGroup(f func(i int), num int, groupSignal <-chan signal) <-chan signal {
    c := make(chan signal)
    var wg sync.WaitGroup
    for i := 0; i < num; i++ {
        wg.Add(1)
        go func(i int) {
            <-groupSignal
            fmt.Printf("worker %d: start to work...\n", i)
            f(i)
            wg.Done()
        }(i + 1)
    }
    go func() {
        wg.Wait()
        c <- signal(struct{}{})
    }()
    return c
}
func main() {
    fmt.Println("start a group of workers...")
    groupSignal := make(chan signal)
    c := spawnGroup(worker, 5, groupSignal)
    time.Sleep(5 * time.Second)
    fmt.Println("the group of workers start to work...")
    close(groupSignal)
    <-c
    fmt.Println("the group of workers work done!")
}

main goroutine 创建了一组 5 个 worker goroutine，这些 Goroutine 启动后会阻塞在名为 groupSignal 的无缓冲 channel 上。
main goroutine 通过close(groupSignal)向所有 worker goroutine 广播“开始工作”的信号，收到 groupSignal 后，所有 worker goroutine 会“同时”开始工作，就像起跑线上的运动员听到了裁判员发出的起跑信号枪声。

关闭一个无缓冲 channel 会让所有阻塞在这个 channel 上的接收操作返回，从而实现了一种 1 对 n 的“广播”机制。

2）第二种用法：用于替代锁机制
无缓冲 channel 具有同步特性，这让它在某些场合可以替代锁，让程序更加清晰，可读性也更好，可以对比下两个方案，直观地感受一下。

首先看一个传统的、基于“共享内存”+“互斥锁”的 Goroutine 安全的计数器的实现：
type counter struct {
    sync.Mutex
    i int
}
var cter counter
func Increase() int {
    cter.Lock()
    defer cter.Unlock()
    cter.i++
    return cter.i
}
func main() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            v := Increase()
            fmt.Printf("goroutine-%d: current counter value is %d\n", i, v)
            wg.Done()
        }(i)
    }
    wg.Wait()
}

在这个示例中，使用了一个带有互斥锁保护的全局变量作为计数器，所有要操作计数器的 Goroutine 共享这个全局变量，并在互斥锁的同步下对计数器进行自增操作。

接下来再看更符合 Go 设计惯例的实现，也就是使用无缓冲 channel 替代锁后的实现：
type counter struct {
    c chan int
    i int
}
func NewCounter() *counter {
    cter := &counter{
        c: make(chan int),
    }
    go func() {
        for {
            cter.i++
            cter.c <- cter.i
        }
    }()
    return cter
}
func (cter *counter) Increase() int {
    return <-cter.c
}
func main() {
    cter := NewCounter()
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            v := cter.Increase()
            fmt.Printf("goroutine-%d: current counter value is %d\n", i, v)
            wg.Done()
        }(i)
    }
    wg.Wait()
}

将计数器操作全部交给一个独立的 Goroutine 去处理，并通过无缓冲 channel 的同步阻塞特性，实现了计数器的控制。
这样其他 Goroutine 通过 Increase 函数试图增加计数器值的动作，实质上就转化为了一次无缓冲 channel 的接收动作。

这种并发设计逻辑更符合 Go 语言所倡导的“不要通过共享内存来通信，而是通过通信来共享内存”的原则。

3、带缓冲 channel 的惯用法

带缓冲的 channel 与无缓冲的 channel 的最大不同之处，就在于它的异步性。
也就是说，对一个带缓冲 channel，在缓冲区未满的情况下，对它进行发送操作的 Goroutine 不会阻塞挂起；
在缓冲区有数据的情况下，对它进行接收操作的 Goroutine 也不会阻塞挂起。

这种特性让带缓冲的 channel 有着与无缓冲 channel 不同的应用场合。

1）第一种用法：用作消息队列

channel 经常被 Go 初学者视为在多个 Goroutine 之间通信的消息队列，这是因为，channel 的原生特性与我们认知中的消息队列十分相似，包括 Goroutine 安全、有 FIFO（first-in, first out）保证等。

其实，和无缓冲 channel 更多用于信号 / 事件管道相比，可自行设置容量、异步收发的带缓冲 channel 更适合被用作为消息队列，并且，带缓冲 channel 在数据收发的性能上要明显好于无缓冲 channel。

单接收单发送性能的基准测试
$go test -bench . *.go

多接收多发送性能基准测试
$go test -bench . *.go

结论：
1、无论是 1 收 1 发还是多收多发，带缓冲 channel 的收发性能都要好于无缓冲 channel；
2、对于带缓冲 channel 而言，发送与接收的 Goroutine 数量越多，收发性能会有所下降；
3、对于带缓冲 channel 而言，选择适当容量会在一定程度上提升收发性能。

2）第二种用法：用作计数信号量（counting semaphore）

Go 并发设计的一个惯用法，就是将带缓冲 channel 用作计数信号量（counting semaphore）。
带缓冲 channel 中的当前数据个数代表的是，当前同时处于活动状态（处理业务）的 Goroutine 的数量，而带缓冲 channel 的容量（capacity），就代表了允许同时处于活动状态的 Goroutine 的最大数量。
向带缓冲 channel 的一个发送操作表示获取一个信号量，而从 channel 的一个接收操作则表示释放一个信号量。

这里来看一个将带缓冲 channel 用作计数信号量的例子：
var active = make(chan struct{}, 3)
var jobs = make(chan int, 10)
func main() {
    go func() {
        for i := 0; i < 8; i++ {
            jobs <- (i + 1)
        }
        close(jobs)
    }()
    var wg sync.WaitGroup
    for j := range jobs {
        wg.Add(1)
        go func(j int) {
            active <- struct{}{}
            log.Printf("handle job: %d\n", j)
            time.Sleep(2 * time.Second)
            <-active
            wg.Done()
        }(j)
    }
    wg.Wait()
}

这个示例创建了一组 Goroutine 来处理 job，同一时间允许最多 3 个 Goroutine 处于活动状态。

为了达成这一目标，看到这个示例使用了一个容量（capacity）为 3 的带缓冲 channel: active 作为计数信号量，这意味着允许同时处于活动状态的最大 Goroutine 数量为 3。

4、len(channel) 的应用
len 是 Go 语言的一个内置函数，它支持接收数组、切片、map、字符串和 channel 类型的参数，并返回对应类型的“长度”，也就是一个整型值。

针对 channel ch 的类型不同，len(ch) 有如下两种语义：
1）当 ch 为无缓冲 channel 时，len(ch) 总是返回 0；
2）当 ch 为带缓冲 channel 时，len(ch) 返回当前 channel ch 中尚未被读取的元素个数。

这样一来，针对带缓冲 channel 的 len 调用似乎才是有意义的。
那我们是否可以使用 len 函数来实现带缓冲 channel 的“判满”、“判有”和“判空”逻辑呢？就像下面示例中伪代码这样：
var ch chan T = make(chan T, capacity)

// 判空
if len(ch) == 0 {
    // 此时channel ch空了?
}

// 判有
if len(ch) > 0 {
    // 此时channel ch中有数据?
}

// 判满
if len(ch) == cap(ch) {
    // 此时channel ch满了?
}
上面代码注释的“空了”、“有数据”和“满了”的后面都打上了问号。这是为什么呢？

这是因为，channel 原语用于多个 Goroutine 间的通信，一旦多个 Goroutine 共同对 channel 进行收发操作，len(channel) 就会在多个 Goroutine 间形成“竞态”。
单纯地依靠 len(channel) 来判断 channel 中元素状态，是不能保证在后续对 channel 的收发时 channel 状态是不变的。

Goroutine1 使用 len(channel) 判空后，就会尝试从 channel 中接收数据。
但在它真正从 channel 读数据之前，另外一个 Goroutine2 已经将数据读了出去，所以，Goroutine1 后面的读取就会阻塞在 channel 上，导致后面逻辑的失效。

因此，为了不阻塞在 channel 上，常见的方法是将“判空与读取”放在一个“事务”中，将“判满与写入”放在一个“事务”中，而这类“事务”我们可以通过 select 实现。
来看下面示例：
func producer(c chan<- int) {
    var i int = 1
    for {
        time.Sleep(2 * time.Second)
        ok := trySend(c, i)
        if ok {
            fmt.Printf("[producer]: send [%d] to channel\n", i)
            i++
            continue
        }
        fmt.Printf("[producer]: try send [%d], but channel is full\n", i)
    }
}
func tryRecv(c <-chan int) (int, bool) {
    select {
    case i := <-c:
        return i, true
    default:
        return 0, false
    }
}
func trySend(c chan<- int, i int) bool {
    select {
    case c <- i:
        return true
    default:
        return false
    }
}
func consumer(c <-chan int) {
    for {
        i, ok := tryRecv(c)
        if !ok {
            fmt.Println("[consumer]: try to recv from channel, but the channel is empty")
            time.Sleep(1 * time.Second)
            continue
        }
        fmt.Printf("[consumer]: recv [%d] from channel\n", i)
        if i >= 3 {
            fmt.Println("[consumer]: exit")
            return
        }
    }
}
func main() {
    var wg sync.WaitGroup
    c := make(chan int, 3)
    wg.Add(2)
    go func() {
        producer(c)
        wg.Done()
    }()
    go func() {
        consumer(c)
        wg.Done()
    }()
    wg.Wait()
}

由于用到了 select 原语的 default 分支语义，
当 channel 空的时候，tryRecv 不会阻塞；
当 channel 满的时候，trySend 也不会阻塞。

在特定的场景下，可以用 len(channel) 来实现。比如下面这两种场景：
1）是一个“多发送单接收”的场景，也就是有多个发送者，但有且只有一个接收者。
在这样的场景下，我们可以在接收 goroutine 中使用len(channel)是否大于0来判断是否 channel 中有数据需要接收。

2）是一个“多接收单发送”的场景，也就是有多个接收者，但有且只有一个发送者。
在这样的场景下，我们可以在发送 Goroutine 中使用len(channel)是否小于cap(channel)来判断是否可以执行向 channel 的发送操作。

5、nil channel 的妙用
如果一个 channel 类型变量的值为 nil，称它为 nil channel。
nil channel 有一个特性，那就是对 nil channel 的读写都会发生阻塞。
比如下面示例代码：
func main() {
  var c chan int
  <-c //阻塞
}
或者：
func main() {
  var c chan int
  c<-1  //阻塞
}

nil channel 的这个特性可不是一无是处，有些时候应用 nil channel 的这个特性可以得到事半功倍的效果。

来看一个例子：
func main() {
    ch1, ch2 := make(chan int), make(chan int)
    go func() {
        time.Sleep(time.Second * 5)
        ch1 <- 5
        close(ch1)
    }()
    go func() {
        time.Sleep(time.Second * 7)
        ch2 <- 7
        close(ch2)
    }()
    var ok1, ok2 bool
    for {
        select {
        case x := <-ch1:
            ok1 = true
            fmt.Println(x)
        case x := <-ch2:
            ok2 = true
            fmt.Println(x)
        }
        if ok1 && ok2 {
            break
        }
    }
    fmt.Println("program end")
}
原本期望上面这个在依次输出 5 和 7 两个数字后退出，但实际运行的输出结果却是在输出 5 之后，程序输出了许多的 0 值，之后才输出 7 并退出。

简单分析一下这段代码的运行过程：
1）前 5s，select 一直处于阻塞状态；
2）第 5s，ch1 返回一个 5 后被 close，select 语句的case x := <-ch1这个分支被选出执行，程序输出 5，并回到 for 循环并重新 select；
3）由于 ch1 被关闭，从一个已关闭的 channel 接收数据将永远不会被阻塞，于是新一轮 select 又把case x := <-ch1这个分支选出并执行。由于 ch1 处于关闭状态，从这个 channel 获取数据，我们会得到这个 channel 对应类型的零值，这里就是 0。于是程序再次输出 0；程序按这个逻辑循环执行，一直输出 0 值；
4）2s 后，ch2 被写入了一个数值 7。这样在某一轮 select 的过程中，分支case x := <-ch2被选中得以执行，程序输出 7 之后满足退出条件，于是程序终止。

5、与 select 结合使用的一些惯用法
第一种用法：利用 default 分支避免阻塞
select 语句的 default 分支的语义，就是在其他非 default 分支因通信未就绪，而无法被选择的时候执行的，这就给 default 分支赋予了一种“避免阻塞”的特性。

其实在前面已经用到了“利用 default 分支”实现的trySend和tryRecv两个函数：

而且，无论是无缓冲 channel 还是带缓冲 channel，这两个函数都能适用，并且不会阻塞在空 channel 或元素个数已经达到容量的 channel 上。

在 Go 标准库中，这个惯用法也有应用，比如：
// $GOROOT/src/time/sleep.go
func sendTime(c interface{}, seq uintptr) {
    // 无阻塞的向c发送当前时间
    select {
    case c.(chan Time) <- Now():
    default:
    }
}

第二种用法：实现超时机制
带超时机制的 select，是 Go 中常见的一种 select 和 channel 的组合用法。
通过超时事件，既可以避免长期陷入某种操作的等待中，也可以做一些异常处理工作。

比如，下面示例代码实现了一次具有 30s 超时的 select：
func worker() {
  select {
  case <-c:
       // ... do some stuff
  case <-time.After(30 *time.Second):
      return
  }
}
不过，在应用带有超时机制的 select 时，我们要特别注意 timer 使用后的释放，尤其在大量创建 timer 的时候。

Go 语言标准库提供的 timer 实际上是由 Go 运行时自行维护的，而不是操作系统级的定时器资源，它的使用代价要比操作系统级的低许多。
但即便如此，作为 time.Timer 的使用者，我们也要尽量减少在使用 Timer 时给 Go 运行时和 Go 垃圾回收带来的压力，要及时调用 timer 的 Stop 方法回收 Timer 资源。

第三种用法：实现心跳机制
结合 time 包的 Ticker，可以实现带有心跳机制的 select。
这种机制让可以在监听 channel 的同时，执行一些周期性的任务，比如下面这段代码：
func worker() {
  heartbeat := time.NewTicker(30 * time.Second)
  defer heartbeat.Stop()
  for {
    select {
    case <-c:
      // ... do some stuff
    case <- heartbeat.C:
      //... do heartbeat stuff
    }
  }
}

这里使用 time.NewTicker，创建了一个 Ticker 类型实例 heartbeat。
这个实例包含一个 channel 类型的字段 C，这个字段会按一定时间间隔持续产生事件，就像“心跳”一样。
这样 for 循环在 channel c 无数据接收时，会每隔特定时间完成一次迭代，然后回到 for 循环进行下一次迭代。

和 timer 一样，在使用完 ticker 之后，也不要忘记调用它的 Stop 方法，避免心跳事件在 ticker 的 channel（上面示例中的 heartbeat.C）中持续产生。


*/
