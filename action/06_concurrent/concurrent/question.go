package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

/*
假设有一个超长的切片，切片的元素类型为int，切片中的元素为乱序排列。
限时5秒，使用多个goroutine查找切片中是否存在给定值，在找到目标值或者超时后立刻结束所有goroutine的执行。

比如切片为：[23, 32, 78, 43, 76, 65, 345, 762, …… 915, 86]，查找的目标值为345，
如果切片中存在目标值程序输出:"Found it!"并且立即取消仍在执行查找任务的goroutine。
如果在超时时间未找到目标值程序输出:"Timeout! Not Found"，同时立即取消仍在执行查找任务的goroutine。

首先题目里提到了在找到目标值或者超时后立刻结束所有goroutine的执行，完成这两个功能需要借助计时器、通道和context才行。
我能想到的第一点就是要用context.WithCancel创建一个上下文对象传递给每个执行任务的goroutine，外部在满足条件后（找到目标值或者已超时）通过调用上下文的取消函数来通知所有goroutine停止工作。

func main() {
    timer := time.NewTimer(time.Second * 5)
    ctx, cancel := context.WithCancel(context.Background())
    resultChan := make(chan bool)
  ......
    select {
    case <-timer.C:
        fmt.Fprintln(os.Stderr, "Timeout! Not Found")
        cancel()
    case <- resultChan:
        fmt.Fprintf(os.Stdout, "Found it!\n")
        cancel()
    }
}

执行任务的goroutine们如果找到目标值后需要通知外部等待任务执行的主goroutine，这个工作是典型的应用通道的场景，上面代码也已经看到了，
我们创建了一个接收查找结果的通道，接下来要做的就是把它和上下文对象一起传递给执行任务的goroutine。
func SearchTarget(ctx context.Context, data []int, target int, resultChan chan bool) {
    for _, v := range data {
        select {
        case <- ctx.Done():
            fmt.Fprintf(os.Stdout, "Task cancelded! \n")
            return
        default:
        }
        // 模拟一个耗时查找，这里只是比对值，真实开发中可以是其他操作
        fmt.Fprintf(os.Stdout, "v: %d \n", v)
        time.Sleep(time.Millisecond * 1500)
        if target == v {
            resultChan <- true
            return
        }
    }

}

在执行查找任务的goroutine里接收上下文的取消信号，为了不阻塞查找任务，我们使用了select语句加default的组合：
select {
case <- ctx.Done():
    fmt.Fprintf(os.Stdout, "Task cancelded! \n")
    return
default:
}

在goroutine里面如果找到了目标值，则会通过发送一个true值给resultChan，让外面等待的主goroutine收到一个已经找到目标值的信号。

resultChan <- true

这样通过上下文的Done通道和resultChan通道，goroutine们就能相互通信了。

Go 语言中最常见的、也是经常被人提及的设计模式 — 不要通过共享内存的方式进行通信，而是应该通过通信的方式共享内存

完整的源代码如下：
*/

func problemSolving() {
	timer := time.NewTimer(time.Second * 5)
	data := []int{1, 2, 3, 10, 999, 8, 345, 7, 98, 33, 66, 77, 88, 68, 96}
	dataLen := len(data)
	size := 3
	target := 345
	ctx, cancel := context.WithCancel(context.Background())
	resultChan := make(chan bool)
	for i := 0; i < dataLen; i += size {
		end := i + size
		if end >= dataLen {
			end = dataLen - 1
		}
		go SearchTarget(ctx, data[i:end], target, resultChan)
	}
	select {
	case <-timer.C:
		_, err := fmt.Fprintln(os.Stderr, "Timeout! Not Found")
		if err != nil {
			return
		}
		cancel()
	case <-resultChan:
		_, err := fmt.Fprintf(os.Stdout, "Found it!\n")
		if err != nil {
			return
		}
		cancel()
	}

	time.Sleep(time.Second * 2)
}

func SearchTarget(ctx context.Context, data []int, target int, resultChan chan bool) {
	for _, v := range data {
		select {
		case <-ctx.Done():
			_, err := fmt.Fprintf(os.Stdout, "Task cancelded! \n")
			if err != nil {
				return
			}
			return
		default:
		}
		// 模拟一个耗时查找，这里只是比对值，真实开发中可以是其他操作
		_, err := fmt.Fprintf(os.Stdout, "v: %d \n", v)
		if err != nil {
			return
		}
		time.Sleep(time.Millisecond * 1500)
		if target == v {
			resultChan <- true
			return
		}
	}
}

/*
v: 10
v: 1
v: 33
v: 88
v: 345
Found it!
Task cancelded!
Task cancelded!
Task cancelded!
v: 999
Task cancelded!

因为是并发程序所以每次打印的结果的顺序是不一样的。
而且也并不是先开启的goroutine就一定会先执行，主要还是看调度器先调度哪个。

》Go语言调度器
所有应用程序都是运行在操作系统上，真正用来干活(计算)的是CPU。所以谈到Go语言调度器，我们也绕不开操作系统、进程与线程这些概念。线程是操作系统调度时的最基本单元，而 Linux 在调度器并不区分进程和线程的调度，它们在不同操作系统上也有不同的实现，但是在大多数的实现中线程都属于进程。

多个线程可以属于同一个进程并共享内存空间。因为多线程不需要创建新的虚拟内存空间，所以它们也不需要内存管理单元处理上下文的切换，线程之间的通信也正是基于共享的内存进行的，与重量级的进程相比，线程显得比较轻量。

虽然线程比较轻量，但是在调度时也有比较大的额外开销。每个线程会都占用 1 兆以上的内存空间，在对线程进行切换时不止会消耗较多的内存，恢复寄存器中的内容还需要向操作系统申请或者销毁对应的资源。

大量的线程出现了新的问题
1、高内存占用

2、调度的CPU高消耗

然后工程师们就发现，其实一个线程分为"内核态"线程和"用户态"线程。

一个用户态线程必须要绑定一个内核态线程，但是CPU并不知道有用户态线程的存在，它只知道它运行的是一个内核态线程(Linux的PCB进程控制块)。这样，我们再去细化分类，内核线程依然叫线程(thread)，用户线程叫协程(co-routine)。既然一个协程可以绑定一个线程，那么也可以通过实现协程调度器把多个协程与一个或者多个线程进行绑定。

Go语言的goroutine来自协程的概念，让一组可复用的函数运行在一组线程之上，即使有协程阻塞，该线程的其他协程也可以被runtime调度，转移到其他可运行的线程上。最关键的是，程序员看不到这些底层的细节，这就降低了编程的难度，提供了更容易的并发。

Go中，协程被称为goroutine，它非常轻量，一个goroutine只占几KB，并且这几KB就足够goroutine运行完，这就能在有限的内存空间内支持大量goroutine，支持了更多的并发。虽然一个goroutine的栈只占几KB，但实际是可伸缩的，如果需要更多内存，runtime会自动为goroutine分配。

既然我们知道了goroutine和系统线程的关系，那么最关键的一点就是实现协程调度器了。

Go目前使用的调度器是2012年重新设计的，因为之前的调度器性能存在问题，所以使用4年就被废弃了。重新设计的调度器使用G-M-P模型并一直沿用至今。

然后工程师们就发现，其实一个线程分为"内核态"线程和"用户态"线程。

一个用户态线程必须要绑定一个内核态线程，但是CPU并不知道有用户态线程的存在，它只知道它运行的是一个内核态线程(Linux的PCB进程控制块)。这样，我们再去细化分类，内核线程依然叫线程(thread)，用户线程叫协程(co-routine)。既然一个协程可以绑定一个线程，那么也可以通过实现协程调度器把多个协程与一个或者多个线程进行绑定。

Go语言的goroutine来自协程的概念，让一组可复用的函数运行在一组线程之上，即使有协程阻塞，该线程的其他协程也可以被runtime调度，转移到其他可运行的线程上。最关键的是，程序员看不到这些底层的细节，这就降低了编程的难度，提供了更容易的并发。

Go中，协程被称为goroutine，它非常轻量，一个goroutine只占几KB，并且这几KB就足够goroutine运行完，这就能在有限的内存空间内支持大量goroutine，支持了更多的并发。虽然一个goroutine的栈只占几KB，但实际是可伸缩的，如果需要更多内存，runtime会自动为goroutine分配。

既然我们知道了goroutine和系统线程的关系，那么最关键的一点就是实现协程调度器了。

Go目前使用的调度器是2012年重新设计的，因为之前的调度器性能存在问题，所以使用4年就被废弃了。重新设计的调度器使用G-M-P模型并一直沿用至今。

G — 表示 goroutine，它是一个待执行的任务；

M — 表示操作系统的线程，它由操作系统的调度器调度和管理；

P — 表示处理器，它可以被看做运行在线程上的本地调度器；

》G
gorotuine 就是Go语言调度器中待执行的任务，它在运行时调度器中的地位与线程在操作系统中差不多，但是它占用了更小的内存空间，也降低了上下文切换的开销。

goroutine只存在于Go语言的运行时，它是Go语言在用户态提供的线程，作为一种粒度更细的资源调度单元，如果使用得当能够在高并发的场景下更高效地利用机器的CPU。

》M
Go语言并发模型中的M是操作系统线程。调度器最多可以创建 10000 个线程，但是其中大多数的线程都不会执行用户代码（可能陷入系统调用），最多只会有 GOMAXPROCS 个活跃线程能够正常运行。

在默认情况下，运行时会将 GOMAXPROCS 设置成当前机器的核数，我们也可以使用 runtime.GOMAXPROCS 来改变程序中最大的线程数。一个四核机器上会创建四个活跃的操作系统线程，每一个线程都对应一个运行时中的 runtime.m 结构体。

在大多数情况下，我们都会使用Go的默认设置，也就是活跃线程数等于CPU个数，在这种情况下不会触发操作系统的线程调度和上下文切换，所有的调度都会发生在用户态，由Go语言调度器触发，能够减少非常多的额外开销。

操作系统线程在Go语言中会使用私有结构体 runtime.m 来表示
type m struct {
    g0   *g
    curg *g
    ...
}

其中g0是持有调度栈的goroutine，curg 是在当前线程上运行的用户goroutine，用户goroutine执行完后，线程切换回g0上，g0会从线程M绑定的P上的等待队列中获取goroutine交给线程。

》P
调度器中的处理器P是线程和goroutine 的中间层，它能提供线程需要的上下文环境，也会负责调度线程上的等待队列，通过处理器P的调度，每一个内核线程都能够执行多个 goroutine，它能在goroutine 进行一些 I/O 操作时及时切换，提高线程的利用率。因为调度器在启动时就会创建 GOMAXPROCS 个处理器，所以Go语言程序的处理器数量一定会等于 GOMAXPROCS，这些处理器会绑定到不同的内核线程上并利用线程的计算资源运行goroutine。

此外在调度器里还有一个全局等待队列，当所有P本地的等待队列被占满后，新创建的goroutine会进入全局等待队列。P的本地队列为空后，M也会从全局队列中拿一批待执行的goroutine放到P本地的等待队列中。

1、全局队列：存放等待运行的G。

2、P的本地队列：同全局队列类似，存放的也是等待运行的G，存的数量有限，不超过256个。新建G时，G优先加入到P的本地队列，如果队列已满，则会把本地队列中一半的G移动到全局队列。

3、P列表：所有的P都在程序启动时创建，并保存在数组中，最多有GOMAXPROCS(可配置)个。

4、M：线程想运行任务就得获取P，从P的本地队列获取G，P队列为空时，M也会尝试从全局队列拿一批G放到P的本地队列，或从其他P的本地队列偷一半放到自己P的本地队列。M运行G，G执行之后，M会从P获取下一个G，不断重复下去。

5、goroutine调度器和OS调度器是通过M结合起来的，每个M都代表了1个内核线程，OS调度器负责把内核线程分配到CPU上执行。

》调度器的策略
调度器的一个策略是尽可能的复用现有的活跃线程，通过以下两个机制提高线程的复用：

1、work stealing机制，当本线程无可运行的G时，尝试从其他线程绑定的P偷取G，而不是销毁线程。

2、hand off机制，当本线程因为G进行系统调用阻塞时，线程释放绑定的P，把P转移给其他空闲的线程执行。

Go的运行时并不具备操作系统内核级的硬件中断能力，基于工作窃取的调度器实现，本质上属于先来先服务的协作式调度，为了解决响应时间可能较高的问题，目前运行时实现了协作式调度和抢占式调度两种不同的调度策略，保证在大部分情况下，不同的 G 能够获得均匀的CPU时间片。

协作式调度依靠被调度方主动弃权，系统监控到一个goroutine运行超过10ms会通过 runtime.Gosched 调用主动让出执行机会。抢占式调度则依靠调度器强制将被调度方被动中断。

推荐其他博主的一篇文章Golang调度器GMP原理与调度全分析，里面用几十张图详细展示了全场景的调度策略解析，让我们更容易理解调度器的GMP模型和它的工作原理。

如果想从Go的源码层面了解调度器的实现，可以看看下面链接这个博主的系列文章。

https://changkun.de/golang/zh-cn/part2runtime/ch06sched/

*/

func main() {
	i := -1
	fmt.Println(uint64(i) << 32)
}
