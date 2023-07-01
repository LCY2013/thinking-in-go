package main

import (
	"sync"
)

/*
》数据竞争
要解释什么是数据竞争我们先来看一段程序：

下面程序getNumber函数中开启了一个单独的goroutine设置变量i的值，同时在不知道开启的goroutine是否已经执行完成的情况下返回了i。所以现在正在发生两个操作：

变量i的值正在被设置成5。

函数getNumber返回了变量i的值。

现在，根据这两个操作中哪一个先完成，最后程序打印出来的值将是0或5。

这就是为什么它被称为数据竞争：getNumber返回的值根据操作1或操作2中的哪一个最先完成而不同。

》检测数据竞争
我们上面代码是一个高度简化的数据竞争示例。在较大的应用程序中，仅靠自己检查代码很难检测到数据竞争。
幸运的是，Go(从V1.1开始)有一个内置的数据竞争检测器，我们可以使用它来确定应用程序里潜在的数据竞争条件。

使用它非常简单，只需在使用Go命令行工具时添加-race标志。例如，让我们尝试使用-race标志来运行我们刚刚编写的程序：
go run -race concurrent.go

fufeng@magic ~/s/g/p/t/a/0/concurrent>go run -race concurrent.go                                                                                                                                                               master!
0
==================
WARNING: DATA RACE
Write at 0x00c0000bc018 by goroutine 7:
  main.getNumber.func1()
      /Users/magicLuoMacBook/software/go/projects/thinking-in-go/action/06_concurrent/concurrent/concurrent.go:32 +0x30

Previous read at 0x00c0000bc018 by main goroutine:
  main.getNumber()
      /Users/magicLuoMacBook/software/go/projects/thinking-in-go/action/06_concurrent/concurrent/concurrent.go:34 +0xb8
  main.main()
      /Users/magicLuoMacBook/software/go/projects/thinking-in-go/action/06_concurrent/concurrent/concurrent.go:38 +0x24

Goroutine 7 (running) created at:
  main.getNumber()
      /Users/magicLuoMacBook/software/go/projects/thinking-in-go/action/06_concurrent/concurrent/concurrent.go:31 +0xae
  main.main()
      /Users/magicLuoMacBook/software/go/projects/thinking-in-go/action/06_concurrent/concurrent/concurrent.go:38 +0x24
==================
Found 1 data race(s)
exit status 66

第一个0是打印结果(因此我们现在知道是操作2首先完成)。接下来的几行给出了在代码中检测到的数据竞争的信息。我们可以看到关于数据竞争的信息分为三个部分：

第一部分告诉我们，在getNumber函数里创建的goroutine中尝试写入（这是我们将值5赋给i的位置）

第二部分告诉我们，在主goroutine里有一个在同时进行的读操作。

第三部分描述了导致数据竞争的goroutine是在哪里被创建的。

除了go run命令外，go build和go test命令也支持使用-race标志。
这个会使编译器创建的应用程序能够记录所有运行期间对共享变量访问，并且会记录下每一个读或者写共享变量的goroutine的身份信息。

竞争检查器会报告所有的已经发生的数据竞争。然而，它只能检测到运行时的竞争条件，并不能证明之后不会发生数据竞争。
由于需要额外的记录，因此构建时加了竞争检测的程序跑起来会慢一些，且需要更大的内存，即使是这样，这些代价对于很多生产环境的工作来说还是可以接受的。
对于一些偶发的竞争条件来说，使用附带竞争检查器的应用程序可以节省很多花在Debug上的时间。

》解决数据竞争的方案
Go提供了很多解决它的选择。所有这些解决方案的思路都是确保在我们写入变量时阻止对该变量的访问。
一般常用的解决数据竞争的方案有：使用WaitGroup锁，使用通道阻塞以及使用Mutex锁，下面我们一个个来看他们的用法并比较一下这几种方案的不同点。

1、使用WaitGroup
解决数据竞争的最直接方法是（如果需求允许的情况下）阻止读取访问，直到写入操作完成。

2、使用通道阻塞
这个方法原则上与上一种方法类似，只是我们使用了通道而不是WaitGroup：

3、使用Mutex
到目前为止，使用的解决方案只有在确定写入操作完成后再去读取i的值时才适用。现在让我们考虑一个更通常的情况，程序读取和写入的顺序并不是固定的，我们只要求它们不能同时发生就行。这种情况下我们应该考虑使用Mutex互斥锁。

》Mutex vs Channel
上面我们使用互斥锁和通道两种方法解决了并发程序的数据竞争问题。那么我们该在什么情况下使用互斥锁，什么情况下又该使用通道呢？答案就在你试图解决的问题中。如果你试图解决的问题更适合互斥锁，那么就继续使用互斥锁。。如果问题似乎更适合渠道，则使用它。

大多数Go新手都试图使用通道来解决所有并发问题，因为这是Go语言的一个很酷的特性。这是不对的。语言为我们提供了使用Mutex或Channel的选项，选择两者都没有错。

通常，当goroutine需要相互通信时使用通道，当确保同一时间只有一个goroutine能访问代码的关键部分时使用互斥锁。在我们上面解决的问题中，我更倾向于使用互斥锁，因为这个问题不需要goroutine之间的任何通信。只需要确保同一时间只有一个goroutine拥有共享变量的使用权，互斥锁本来就是为解决这种问题而生的，所以使用互斥锁是更自然的一种选择。
*/

// SafeNumber 首先，创建一个结构体包含我们想用互斥锁保护的值和一个mutex实例
type SafeNumber struct {
	val int
	m   sync.Mutex
}

func (i *SafeNumber) Get() int {
	i.m.Lock()
	defer i.m.Unlock()
	return i.val
}

func (i *SafeNumber) Set(val int) {
	i.m.Lock()
	defer i.m.Unlock()
	i.val = val
}

func getNumberMutex() int {
	// 创建一个sageNumber实例
	i := &SafeNumber{}
	// 使用Set和Get代替常规赋值和读取操作。
	// 我们现在可以确保只有在写入完成时才能读取，反之亦然
	go func() {
		i.Set(5)
	}()
	return i.Get()
}

func getNumberWg() int {
	var i int
	// 初始化一个WaitGroup
	var wg sync.WaitGroup
	// Add(1) 通知程序有一个需要等待完成的任务
	wg.Add(1)
	go func() {
		i = 5
		// 调用wg.Done 表示正在等待的程序已经执行完成了
		wg.Done()
	}()
	// wg.Wait会阻塞当前程序直到等待的程序都执行完成为止
	wg.Wait()
	return i
}

func getNumber() int {
	var i int
	go func() {
		i = 5
	}()
	return i
}

func getNumberChannel() int {
	var i int
	// 创建一个通道，在等待的任务完成时会向通道发送一个空结构体
	done := make(chan struct{})
	go func() {
		i = 5
		// 执行完成后向通道发送一个空结构体
		done <- struct{}{}
	}()
	// 从通道接收值将会阻塞程序，直到有值发送给done通道为止
	<-done
	return i
}

//func main() {
//	//fmt.Println(getNumber())
//	problemSolving()
//}
