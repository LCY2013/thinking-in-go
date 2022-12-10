package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

/*
如果信号量是一个任意的整数，通常被称为计数信号量（Counting semaphore），或一般信号量（general semaphore）；
如果信号量只有二进制的0或1，称为二进制信号量（binary semaphore）。
在linux系统中，二进制信号量（binary semaphore）又称互斥锁（Mutex）

计数信号量具备两种操作动作，称为V（signal()）与P（wait()）（即部分参考书常称的“PV操作”）。V操作会增加信号量S的数值，P操作会减少它。

运行方式：
1、初始化信号量，给与它一个非负数的整数值。
2、运行P（wait()），信号量S的值将被减少。企图进入临界区的进程，需要先运行P（wait()）。当信号量S减为负值时，进程会被阻塞住，不能继续；当信号量S不为负值时，进程可以获准进入临界区。
3、运行V（signal()），信号量S的值会被增加。结束离开临界区的进程，将会运行V（signal()）。当信号量S不为负值时，先前被阻塞住的其他进程，将可获准进入临界区。

一般用信号量保护一组资源，比如数据库连接池、一组客户端的连接等等。
每次获取资源时都会将信号量中的计数器减去对应的数值，在释放资源时重新加回来。
当遇到信号量资源不够时尝试获取的线程就会进入休眠，等待其他线程释放归还信号量。
如果信号量是只有0和1的二进位信号量，那么，它的 P/V 就和互斥锁的 Lock/Unlock 一样了。

Go 内部使用信号量来控制goroutine的阻塞和唤醒，比如互斥锁sync.Mutex结构体定义的第二个字段就是一个信号量。

type Mutex struct {
    state int32
    sema  uint32
}

信号量的PV操作在Go内部是通过下面这几个底层函数实现的

func runtime_Semacquire(s *uint32)
func runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int)
func runtime_Semrelease(s *uint32, handoff bool, skipframes int)

上面几个函数都是Go语言内部使用的，不能在编程时直接使用。
不过Go 语言的扩展并发原语包中提供了带权重的信号量 semaphore.Weighted

使用信号量前，需先在项目里安装golang.org/x/sync/包

安装方法：go get -u golang.org/x/sync

可以按照不同的权重对资源的访问进行管理，这个结构体对外提供了四个方法：

1、semaphore.NewWeighted 用于创建新的信号量，通过参数(n int64) 指定信号量的初始值。

2、semaphore.Weighted.Acquire 阻塞地获取指定权重的资源，如果当前没有空闲资源，就会陷入休眠等待；相当于 P 操作，你可以一次获取多个资源，如果没有足够多的资源，调用者就会被阻塞。它的第一个参数是 Context，这就意味着，你可以通过 Context 增加超时或者 cancel 的机制。如果是正常获取了资源，就返回 nil；否则，就返回ctx.Err()，信号量不改变。

3、semaphore.Weighted.Release 用于释放指定权重的资源；相当于 V 操作，可以将 n 个资源释放，返还给信号量。

4、semaphore.Weighted.TryAcquire 非阻塞地获取指定权重的资源，如果当前没有空闲资源，就会直接返回 false；

在实际应用Go语言开发程序时，有哪些场景适合使用信号量呢？
在需要控制访问资源的线程数量时就会需要信号量，来举个例子帮助你理解。
假设我们有一组要抓取的页面，资源有限最多允许我们同时执行三个抓取任务，当同时有三个抓取任务在执行时，在执行完一个抓取任务后才能执行下一个排队等待的任务。
当然这个问题用Channel也能解决，不过这次我们使用Go提供的信号量原语来解决这个问题，代码如下：
*/

// doSomething 模拟抓取任务执行
func doSomething(url string) {
	fmt.Println(url)
	time.Sleep(2 * time.Second)
}

const (
	Limit  = 3 // 同时运行的 goroutine 数量
	Weight = 1 // 每个 goroutine 获取信号量的权重
)

func main() {
	urls := []string{
		"https://www.example.org/1",
		"https://www.example.org/2",
		"https://www.example.org/3",
		"https://www.example.org/4",
		"https://www.example.org/5",
	}
	sem := semaphore.NewWeighted(Limit)
	var wg sync.WaitGroup

	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			err := sem.Acquire(context.Background(), Weight)
			if err != nil {
				return
			}
			doSomething(url)
			sem.Release(Weight)
		}(u)
	}

	wg.Wait()
	fmt.Println("All Done")

	func() {
		defer func() {
			if err := recover(); err != nil {
				log.Fatal(err)
			}
		}()

		fmt.Println("func alone")
	}()
}

/*
Go语言扩展库中的信号量是使用互斥锁和List 实现的。
互斥锁实现其它字段的保护，而 List 实现了一个等待队列，等待者的通知是通过 Channel 的通知机制实现的。

《信号量semaphore.Weighted的数据结构》
type Weighted struct {
    size    int64         // 最大资源数
    cur     int64         // 当前已被使用的资源
    mu      sync.Mutex    // 互斥锁，对字段的保护
    waiters list.List     // 等待队列
}

size: 字段用来记录信号量拥有的最大资源数。
cur: 标识当前已被使用的资源数。
mu: 是一个互斥锁用来提供对其他字段的临界区保护。
waiters: 表示申请资源时由于可使用资源不够而陷入阻塞等待的调用者列表。

《Acquire请求信号量资源》
Acquire方法会监控资源是否可用，而且还要检测传递进来的context.Context对象是否发送了超时过期或者取消的信号，我们来看一下它的代码实现：
func (s *Weighted) Acquire(ctx context.Context, n int64) error {
    s.mu.Lock()
    // 如果恰好有足够的资源，也没有排队等待获取资源的goroutine，
    // 将cur加上n后直接返回
    if s.size-s.cur >= n && s.waiters.Len() == 0 {
      s.cur += n
      s.mu.Unlock()
      return nil
    }

    // 请求的资源数大于能提供的最大的资源数
    // 这个任务处理不了，走错误处理逻辑
    if n > s.size {
      s.mu.Unlock()
      // 依赖ctx的状态返回，否则一直等待
      <-ctx.Done()
      return ctx.Err()
    }
    // 现存资源不够, 需要把调用者加入到等待队列中
    // 创建了一个ready chan,以便被通知唤醒
    ready := make(chan struct{})
    w := waiter{n: n, ready: ready}
    elem := s.waiters.PushBack(w)
    s.mu.Unlock()


    // 等待
    select {
    case <-ctx.Done(): // context的Done被关闭
      err := ctx.Err()
      s.mu.Lock()
      select {
      case <-ready: // 如果被唤醒了，忽略ctx的状态
        err = nil
      default: // 通知waiter
        isFront := s.waiters.Front() == elem
        s.waiters.Remove(elem)
        // 通知其它的waiters,检查是否有足够的资源
        if isFront && s.size > s.cur {
          s.notifyWaiters()
        }
      }
      s.mu.Unlock()
      return err
    case <-ready: // 等待者被唤醒了
      return nil
    }
  }
}

如果调用者请求不到信号量的资源就会被加入等待者列表里，这里等待者列表的结构体定义是：
type waiter struct {
 n     int64
 ready chan<- struct{} // 当调用者可以获取到信号量资源时, close调这个chan
}
包含了两个字段，调用者请求的资源数，以及一个ready 通道。ready通道会在调用者可以被重新唤醒的时候被close调，从而起到通知正在阻塞读取ready通道的等待者的作用。

《NotifyWaiters 通知等待者》
notifyWaiters方法会逐个检查队列里等待的调用者，如果现存资源够等待者请求的数量n，或者是没有等待者了，就返回：

func (s *Weighted) notifyWaiters() {
    for {
      next := s.waiters.Front()
      if next == nil {
        break // 没有等待者了，直接返回
      }


      w := next.Value.(waiter)
      if s.size-s.cur < w.n {
        // 如果现有资源不够队列头调用者请求的资源数，就退出所有等待者会继续等待
        // 这里还是按照先入先出的方式处理是为了避免饥饿
        break
      }

      s.cur += w.n
      s.waiters.Remove(next)
      close(w.ready)
    }
  }
}

notifyWaiters方法是按照先入先出的方式唤醒调用者。
当释放 100 个资源的时候，如果第一个等待者需要 101 个资源，那么，队列中的所有等待者都会继续等待，即使队列后面有的等待者只需要 1 个资源。
这样做的目的是避免饥饿，否则的话，资源可能总是被那些请求资源数小的调用者获取，这样一来，请求资源数巨大的调用者，就没有机会获得资源了。

《Release归还信号量资源》
Release方法就很简单了，它将当前计数值减去释放的资源数 n，并调用notifyWaiters方法，尝试唤醒等待队列中的调用者，看是否有足够的资源被获取。

func (s *Weighted) Release(n int64) {
    s.mu.Lock()
    s.cur -= n
    if s.cur < 0 {
      s.mu.Unlock()
      panic("semaphore: released more than held")
    }
    s.notifyWaiters()
    s.mu.Unlock()
}

在Go语言中信号量有时候也会被Channel类型所取代，因为一个 buffered chan 也可以代表 n 个资源。不过既然Go语言通过golang.orgx/sync扩展库对外提供了semaphore.Weight这一种信号量实现，遇到使用信号量的场景时还是尽量使用官方提供的实现。在使用的过程中我们需要注意以下的几个问题：

Acquire和 TryAcquire方法都可以用于获取资源，前者会阻塞地获取信号量。后者会非阻塞地获取信号量，如果获取不到就返回false。
Release归还信号量后，会以先进先出的顺序唤醒等待队列中的调用者。如果现有资源不够处于等待队列前面的调用者请求的资源数，所有等待者会继续等待。
如果一个goroutine申请较多的资源，由于上面说的归还后唤醒等待者的策略，它可能会等待比较长的时间。
*/
