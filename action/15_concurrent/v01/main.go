package main

import (
	"errors"
	"fmt"
	"time"
)

/*
1、Go 并发

Go 的设计者敏锐地把握了 CPU 向多核方向发展的这一趋势，在决定去创建 Go 语言的时候，他们果断将面向多核、原生支持并发作为了 Go 语言的设计目标之一，并将面向并发作为 Go 的设计哲学。
当 Go 语言首次对外发布时，对并发的原生支持成为了 Go 最令开发者着迷的语法特性之一。

怎么去学习 Go 并发呢？我的方法是将“Go 并发”这个词拆开来看，它包含两方面内容：
一个是并发的概念，
另一个是 Go 针对并发设计给出的自身的实现方案，也就是 goroutine、channel、select 这些 Go 并发的语法特性。

2、什么是并发？

并行（parallelism），指的就是在同一时刻，有两个或两个以上的任务（这里指进程）的代码在处理器上执行。

粗略看起来，多进程应用与单进程应用相比并没有什么质的提升。那我们为什么还要将应用设计为多进程呢？

这更多是从应用的结构角度去考虑的，多进程应用由于将功能职责做了划分，并指定专门的模块来负责，所以从结构上来看，要比单进程更为清晰简洁，可读性与可维护性也更好。
这种将程序分成多个可独立执行的部分的结构化程序的设计方法，就是并发设计。采用了并发设计的应用也可以看成是一组独立执行的模块的组合。

并发不是并行，并发关乎结构，并行关乎执行。

3、Go 的并发方案：goroutine

Go 并没有使用操作系统线程作为承载分解后的代码片段（模块）的基本执行单元，而是实现了goroutine这一由 Go 运行时（runtime）负责调度的、轻量的用户级线程，为并发程序设计提供原生支持。

这一方案有啥优势，相比传统操作系统线程来说，goroutine 的优势主要是：

  1、资源占用小，每个 goroutine 的初始栈大小仅为 2k；

  2、由 Go 运行时而不是操作系统调度，goroutine 上下文切换在用户层完成，开销更小；

  3、在语言层面而不是通过标准库提供。goroutine 由go关键字创建，一退出就会被回收或销毁，开发体验更佳；

  4、语言内置 channel 作为 goroutine 间通信原语，为并发设计提供了强大支撑。


4、goroutine 的基本用法

并发是一种能力，它让你的程序可以由若干个代码片段组合而成，并且每个片段都是独立运行的。
goroutine 恰恰就是 Go 原生支持并发的一个具体实现。
无论是 Go 自身运行时代码还是用户层 Go 代码，都无一例外地运行在 goroutine 中。

Go 语言通过go关键字+函数/方法的方式创建一个 goroutine。
创建后，新 goroutine 将拥有独立的代码执行流，并与创建它的 goroutine 一起被 Go 运行时调度。

创建 goroutine 的代码示例：
go fmt.Println("I am a goroutine")
var c = make(chan int)
go func(a, b int) {
    c <- a + b
}(3,4)

// $GOROOT/src/net/http/server.go
c := srv.newConn(rw)
go c.serve(connCtx)

了解了怎么创建，那怎么退出 goroutine 呢？

goroutine 的使用代价很低，Go 官方也推荐你多多使用 goroutine。
而且，多数情况下，不需要考虑对 goroutine 的退出进行控制：goroutine 的执行函数的返回，就意味着 goroutine 退出。

如果 main goroutine 退出了，那么也意味着整个应用程序的退出。
此外，你还要注意的是，goroutine 执行的函数或方法即便有返回值，Go 也会忽略这些返回值。
所以，如果你要获取 goroutine 执行后的返回值，你需要另行考虑其他方法，比如通过 goroutine 间的通信来实现。

4、goroutine 间的通信
传统的编程语言（比如：C++、Java、Python 等）并非面向并发而生的，所以他们面对并发的逻辑多是基于操作系统的线程。
并发的执行单元（线程）之间的通信，利用的也是操作系统提供的线程或进程间通信的原语，比如：共享内存、信号（signal）、管道（pipe）、消息队列、套接字（socket）等。

在这些通信原语中，使用最多、最广泛的（也是最高效的）是结合了线程同步原语（比如：锁以及更为低级的原子操作）的共享内存方式，因此，可以说传统语言的并发模型是基于对内存的共享的。

go 在新并发模型设计中借鉴了著名计算机科学家Tony Hoare提出的 CSP（Communicationing Sequential Processes，通信顺序进程）并发模型。

Tony Hoare 的 CSP 模型旨在简化并发程序的编写，让并发程序的编写与编写顺序程序一样简单。
Tony Hoare 认为输入输出应该是基本的编程原语，数据处理逻辑（也就是 CSP 中的 P）只需调用输入原语获取数据，顺序地处理数据，并将结果数据通过输出原语输出就可以了。

因此，在 Tony Hoare 眼中，一个符合 CSP 模型的并发程序应该是一组通过输入输出原语连接起来的 P 的集合。从这个角度来看，CSP 理论不仅是一个并发参考模型，也是一种并发程序的程序组织方法。它的组合思想与 Go 的设计哲学不谋而合。

Tony Hoare 的 CSP 理论中的 P，也就是“Process（进程）”，是一个抽象概念，它代表任何顺序处理逻辑的封装，它获取输入数据（或从其他 P 的输出获取），并生产出可以被其他 P 消费的输出数据。这里我们可以简单看下 CSP 通信模型的示意图：

注意了，这里的 P 并不一定与操作系统的进程或线程划等号。在 Go 中，与“Process”对应的是 goroutine。
为了实现 CSP 并发模型中的输入和输出原语，Go 还引入了 goroutine（P）之间的通信原语channel。goroutine 可以从 channel 获取输入数据，再将处理后得到的结果数据通过 channel 输出。通过 channel 将 goroutine（P）组合连接在一起，让设计和编写大型并发系统变得更加简单和清晰，我们再也不用为那些传统共享内存并发模型中的问题而伤脑筋了。

比如获取 goroutine 的退出状态，就可以使用 channel 原语实现：
func spawn(f func() error) <-chan error {
    c := make(chan error)
    go func() {
        c <- f()
    }()
    return c
}
func main() {
    c := spawn(func() error {
        time.Sleep(2 * time.Second)
        return errors.New("timeout")
    })
    fmt.Println(<-c)
}

虽然 CSP 模型已经成为 Go 语言支持的主流并发模型，但 Go 也支持传统的、基于共享内存的并发模型，并提供了基本的低级别同步原语（主要是 sync 包中的互斥锁、条件变量、读写锁、原子操作等）。

从程序的整体结构来看，Go 始终推荐以 CSP 并发模型风格构建并发程序，尤其是在复杂的业务层面，这能提升程序的逻辑清晰度，大大降低并发设计的复杂性，并让程序更具可读性和可维护性。

不过，对于局部情况，比如涉及性能敏感的区域或需要保护的结构体数据时，我们可以使用更为高效的低级同步原语（如 mutex），保证 goroutine 对数据的同步访问。



*/

func spawn(f func() error) <-chan error {
	c := make(chan error)
	go func() {
		c <- f()
	}()
	return c
}

/*
这个示例在 main goroutine 与子 goroutine 之间建立了一个元素类型为 error 的 channel，
子 goroutine 退出时，会将它执行的函数的错误返回值写入这个 channel，main goroutine 可以通过读取 channel 的值来获取子 goroutine 的退出状态。
*/
func main() {
	c := spawn(func() error {
		time.Sleep(2 * time.Second)
		return errors.New("timeout")
	})
	fmt.Println(<-c)
}
