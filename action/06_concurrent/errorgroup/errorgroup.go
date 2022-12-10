package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

/*
在并发编程里，sync.WaitGroup并发原语的使用频率非常高，经常用于协同等待场景：一个goroutine在检查点(Check Point)等待一组执行任务的 worker goroutine 全部完成，如果在执行任务的这些worker goroutine 还没全部完成，等待的 goroutine 就会阻塞在检查点，直到所有woker  goroutine 都完成后才能继续执行。

如果在woker goroutine的执行过程中遇到错误想要通知在检查点等待的协程处理该怎么办呢？WaitGroup并没有提供传播错误的功能。Go语言在扩展库提供的ErrorGroup并发原语正好适合在这种场景下使用，它在WaitGroup的功能基础上还提供了，错误传播以及上下文取消的功能。

Go扩展库通过errorgroup.Group提供ErrorGroup原语的功能，它有三个方法可调用：
func WithContext(ctx context.Context) (*Group, context.Context)
func (g *Group) Go(f func() error)
func (g *Group) Wait() error

调用errorgroup包的WithContext方法会返回一个Group 实例，同时还会返回一个使用 context.WithCancel 生成的新Context。一旦有一个子任务返回错误，或者是Wait 调用返回，这个新 Context 就会被 cancel。
Go方法，接收类型为func() error 的函数作为子任务函数，如果任务执行成功，就返回nil，否则就返回 error，并且会cancel 那个新的Context。
Wait方法，类似WaitGroup的 Wait 方法，调用后会阻塞地等待所有的子任务都完成，它才会返回。如果有多个子任务返回错误，它只会返回第一个出现的错误，如果所有的子任务都执行成功，就返回nil。

接下来我们让主goroutine使用ErrorGroup代替WaitGroup等待所有子任务的完成，ErrorGroup有一个特点是会返回所有执行任务的goroutine遇到的第一个错误。我们试着执行一下下面的程序，注意观察程序的输出。
*/

/*func main() {
	var eg errgroup.Group
	for i := 0; i < 100; i++ {
		i := i
		eg.Go(func() error {
			time.Sleep(time.Second * 2)
			if i > 90 {
				fmt.Println("Error: ", i)
				return fmt.Errorf("error occurred: %d", i)
			}
			fmt.Println("End: ", i)
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}*/
/*
上面程序，遇到i大于90的都会产生错误结束执行，但是只有第一个产生的错误被ErrorGroup返回，程序的输出大概如下：
End:  73
End:  39
Error:  92
End:  2
End:  21
End:  89
End:  0
End:  82
End:  32
End:  8
End:  31
error occurred: 98

最早执行遇到错误的goroutine输出了Error: 98但是所有未执行完的其他任务并没有停止执行。
那么想让程序遇到错误就终止其他子任务该怎么办呢？
可以用errgroup.Group提供的WithContext方法创建一个带可取消上下文功能的ErrorGroup。

*/

func main() {
	eg, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < 100; i++ {
		i := i
		eg.Go(func() error {
			time.Sleep(time.Second * 2)
			select {
			case <-ctx.Done():
				fmt.Println("Canceled: ", i)
				return nil
			default:
				if i > 90 {
					fmt.Println("Error: ", i)
					return fmt.Errorf("error occurred: %d", i)
				}
				fmt.Println("End: ", i)
				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}

/*
Go方法单独开启的goroutine在执行参数传递进来的函数时，如果函数返回了错误，会对ErrorGroup持有的err字段进行赋值并及时调用cancel函数，通过上下文通知其他子任务取消执行任务。

ErrorGroup原语的结构体类型errorgroup.Group定义如下：
type Group struct {
 cancel func()

 wg sync.WaitGroup

 errOnce sync.Once
 err     error
}

cancel — 创建 context.Context 时返回的取消函数，用于在多个 goroutine 之间同步取消信号；

wg — 用于等待一组 goroutine 完成子任务的同步原语；

errOnce — 用于保证只接收一个子任务返回的错误的同步原语；

通过 errgroup.WithContext构造器创建errgroup.Group 结构体：
func WithContext(ctx context.Context) (*Group, context.Context) {
 ctx, cancel := context.WithCancel(ctx)
 return &Group{cancel: cancel}, ctx
}

运行新的并行子任务需要使用errgroup.Group.Go方法，这个方法的执行过程如下：

调用 sync.WaitGroup.Add 增加待处理的任务数；
创建一个新的 goroutine 并在 goroutine 内部运行子任务；
返回错误时及时调用 cancel 并对 err 赋值，只有最早返回的错误才会被上游感知到，后续的错误都会被舍弃：

func (g *Group) Go(f func() error) {
 g.wg.Add(1)

 go func() {
  defer g.wg.Done()

  if err := f(); err != nil {
   g.errOnce.Do(func() {
    g.err = err
    if g.cancel != nil {
     g.cancel()
    }
   })
  }
 }()
}

用于等待的errgroup.Group.Wait方法只是调用了 sync.WaitGroup.Wait方法，阻塞地等待所有子任务完成。
在子任务全部完成时会通过调用在errorgroup.WithContext创建Group和Context对象时存放在Group.cancel字段里的函数，取消Context对象并返回可能出现的错误。

func (g *Group) Wait() error {
 g.wg.Wait()
 if g.cancel != nil {
  g.cancel()
 }
 return g.err
}

Go语言通过errorgroup.Group结构提供的ErrorGroup原语，通过封装WaitGroup、Once基本原语结合上下文对象，提供了除同步等待外更加复杂的错误传播和执行任务取消的功能。

在使用时，我们也需要注意它的两个特点：

errgroup.Group在出现错误或者等待结束后都会调用 Context对象 的 cancel 方法同步取消信号。
只有第一个出现的错误才会被返回，剩余的错误都会被直接抛弃。
*/
