package main

import (
	"fmt"
	"time"
)

/*
Go语言计时器

Go语言的标准库里提供两种类型的计时器Timer和Ticker。Timer经过指定的duration时间后被触发，往自己的时间channel发送当前时间，此后Timer不再计时。Ticker则是每隔duration时间都会把当前时间点发送给自己的时间channel，利用计时器的时间channel可以实现很多与计时相关的功能。

文章主要涉及如下内容：

1、Timer和Ticker计时器的内部结构表示

2、Timer和Ticker的使用方法和注意事项

3、如何正确Reset定时器

》计时器的内部表示

两种计时器都是基于Go语言的运行时计时器runtime.timer实现的，rumtime.timer的结构体表示如下：
type timer struct {
    pp puintptr

    when     int64
    period   int64
    f        func(interface{}, uintptr)
    arg      interface{}
    seq      uintptr
    nextwhen int64
    status   uint32
}

rumtime.timer结构体中的字段含义是

1、when — 当前计时器被唤醒的时间；

2、period — 两次被唤醒的间隔；

3、f — 每当计时器被唤醒时都会调用的函数；

3、arg — 计时器被唤醒时调用 f 传入的参数；

4、nextWhen — 计时器处于 timerModifiedLater/timerModifiedEairlier 状态时，用于设置 when 字段；

5、status — 计时器的状态；

这里的 runtime.timer 只是私有的计时器运行时表示，对外暴露的计时器 time.Timer 和time.Ticker的结构体表示如下：

type Timer struct {
    C <-chan Time
    r runtimeTimer
}

type Ticker struct {
    C <-chan Time
    r runtimeTimer
}

Timer.C和Ticker.C就是计时器中的时间channel，接下来我们看一下怎么使用这两种计时器，以及使用时要注意的地方。

》Timer计时器

time.Timer 计时器必须通过 time.NewTimer、time.AfterFunc 或者 time.After 函数创建。当计时器失效时，失效的时间就会被发送给计时器持有的 channel，订阅 channel 的 goroutine 会收到计时器失效的时间。

通过定时器Timer用户可以定义自己的超时逻辑，尤其是在应对使用select处理多个channel的超时、单channel读写的超时等情形时尤为方便。Timer常见的使用方法如下：

//使用time.AfterFunc：

t := time.AfterFunc(d, f)

//使用time.After：
select {
    case m := <-c:
       handle(m)
    case <-time.After(5 * time.Minute):
       fmt.Println("timed out")
}

// 使用time.NewTimer:
t := time.NewTimer(5 * time.Minute)
select {
    case m := <-c:
       handle(m)
    case <-t.C:
       fmt.Println("timed out")
}

time.AfterFunc这种方式创建的Timer，在到达超时时间后会在单独的goroutine里执行函数f。

func AfterFunc(d Duration, f func()) *Timer {
    t := &Timer{
        r: runtimeTimer{
            when: when(d),
            f:    goFunc,
            arg:  f,
        },
    }
    startTimer(&t.r)
    return t
}

func goFunc(arg interface{}, seq uintptr) {
    go arg.(func())()
}

从上面AfterFunc的源码可以看到外面传入的f参数并非直接赋值给了运行时计时器的f，而是作为包装函数goFunc的参数传入的。
goFunc会启动了一个新的goroutine来执行外部传入的函数f。
这是因为所有计时器的事件函数都是由Go运行时内唯一的goroutine timerproc运行的。
为了不阻塞timerproc的执行，必须启动一个新的goroutine执行到期的事件函数。

对于NewTimer和After这两种创建方法，则是Timer在超时后，执行一个标准库中内置的函数：sendTime。

func NewTimer(d Duration) *Timer {
    c := make(chan Time, 1)
    t := &Timer{
        C: c,
        r: runtimeTimer{
            when: when(d),
            f:    sendTime,
            arg:  c,
        },
    }
    startTimer(&t.r)
    return t
}

func sendTime(c interface{}, seq uintptr) {
    select {
    case c.(chan Time) <- Now():
    default:
    }
}

sendTime将当前时间发送到Timer的时间channel中。那么这个动作不会阻塞timerproc的执行么？
答案是不会，原因是NewTimer创建的是一个带缓冲的channel所以无论Timer.C这个channel有没有接收方sendTime都可以非阻塞的将当前时间发送给Timer.C，而且sendTime中还加了双保险：
通过select判断Timer.C的Buffer是否已满，一旦满了，会直接退出，依然不会阻塞。

Timer的Stop方法可以阻止计时器触发，调用Stop方法成功停止了计时器的触发将会返回true，如果计时器已经过期了或者已经被Stop停止过了，再次调用Stop方法将会返回false。

Go运行时将所有计时器维护在一个最小堆Min Heap中，Stop一个计时器就是从堆中删除该计时器。

》Ticker计时器

Ticker可以周期性地触发时间事件，每次到达指定的时间间隔后都会触发事件。

time.Ticker需要通过time.NewTicker或者time.Tick创建。

// 使用time.Tick:
go func() {
    for t := range time.Tick(time.Minute) {
        fmt.Println("Tick at", t)
    }
}()

// 使用time.Ticker
var ticker *time.Ticker = time.NewTicker(1 * time.Second)

go func() {
    for t := range ticker.C {
        fmt.Println("Tick at", t)
    }
}()

time.Sleep(time.Second * 5)
ticker.Stop()
fmt.Println("Ticker stopped")

不过time.Tick很少会被用到，除非你想在程序的整个生命周期里都使用time.Ticker的时间channel。官文文档里对time.Tick的描述是：
》time.Tick底层的Ticker不能被垃圾收集器恢复；

所以使用time.Tick时一定要小心，为避免意外尽量使用time.NewTicker返回的Ticker替代。

NewTicker创建的计时器与NewTimer创建的计时器持有的时间channel一样都是带一个缓存的channel，每次触发后执行的函数也是sendTime，这样即保证了无论有无接收方Ticker触发时间事件时都不会阻塞：

func NewTicker(d Duration) *Ticker {
    if d <= 0 {
        panic(errors.New("non-positive interval for NewTicker"))
    }
    // Give the channel a 1-element time buffer.
    // If the client falls behind while reading, we drop ticks
    // on the floor until the client catches up.
    c := make(chan Time, 1)
    t := &Ticker{
        C: c,
        r: runtimeTimer{
            when:   when(d),
            period: int64(d),
            f:      sendTime,
            arg:    c,
        },
    }
    startTimer(&t.r)
    return t
}

》Reset计时器时要注意的问题

关于Reset的使用建议，文档里的描述是：

重置计时器时必须注意不要与当前计时器到期发送时间到t.C的操作产生竞争。
如果程序已经从t.C接收到值，则计时器是已知的已过期，并且t.Reset可以直接使用。
如果程序尚未从t.C接收值，计时器必须先被停止，并且-如果使用t.Stop时报告计时器已过期，那么请排空其通道中值。

例如：

if !t.Stop() {
  <-t.C
}
t.Reset(d)

下面的例子里producer goroutine里每一秒向通道中发送一个false值，循环结束后等待一秒再往通道里发送一个true值。
在consumer goroutine里通过循环试图从通道中读取值，用计时器设置了最长等待时间为5秒，
如果计时器超时了，输出当前时间并进行下次循环尝试，
如果从通道中读取出的不是期待的值（预期值是true)，则尝试重新从通道中读取并重置计时器。
*/

func main() {
	c := make(chan bool)

	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second * 6)
			c <- false
		}

		time.Sleep(time.Second * 6)
		c <- true
	}()

	go func() {
		// try to read from channel, block at most 5s.
		// if timeout, print time event and go on loop.
		// if read a message which is not the type we want(we want true, not false),
		// retry to read.
		timer := time.NewTimer(time.Second * 5)
		for {
			// timer is active , not fired, stop always returns true, no problems occurs.
			if !timer.Stop() {
				//<-timer.C
				select {
				case <-timer.C: //try to drain from the channel
				default:
				}
			}
			timer.Reset(time.Second * 5)

			select {
			case b := <-c:
				if b == false {
					fmt.Println(time.Now(), ":recv false. continue")
					continue
				}
				//we want true, not false
				fmt.Println(time.Now(), ":recv true. return")
				return
			case <-timer.C:
				fmt.Println(time.Now(), ":timer expired")
				continue
			}
		}
	}()

	//to avoid that all goroutine blocks.
	var s string
	_, err := fmt.Scanln(&s)
	if err != nil {
		return
	}
}

/*
2022-05-29 17:58:05.862874 +0800 CST m=+1.001459437 :recv false. continue
2022-05-29 17:58:06.863207 +0800 CST m=+2.001832659 :recv false. continue
2022-05-29 17:58:07.864081 +0800 CST m=+3.002650312 :recv false. continue
2022-05-29 17:58:08.864289 +0800 CST m=+4.002850270 :recv false. continue
2022-05-29 17:58:09.864458 +0800 CST m=+5.003011634 :recv false. continue
2022-05-29 17:58:10.864601 +0800 CST m=+6.003145987 :recv true. return

目前来看没什么问题，使用Reset重置计时器也起作用了，接下来我们对producer goroutin做一些更改，
我们把producer goroutine里每秒发送值的逻辑改成每6秒发送值，而consumer gouroutine里和计时器还是5秒就到期。

再次运行会发现程序发生了deadlock在第一次报告计时器过期后直接阻塞住了：
2022-05-29 17:58:51.629279 +0800 CST m=+5.001445413 :timer expired

那程序是在哪阻塞住的呢？对就是在抽干timer.C通道时阻塞住了（英文叫做drain channel比喻成流干管道里的水，在程序里就是让timer.C管道中不再存在未接收的值)。

if !timer.Stop() {
    <-timer.C
}
timer.Reset(time.Second * 5)

producer goroutine的发送行为发生了变化，comsumer goroutine在收到第一个数据前有了一次计时器过期的事件，for循环进行一下次循环。
这时timer.Stop函数返回的不再是true，而是false，因为计时器已经过期了，上面提到的维护着所有活跃计时器的最小堆中已经不包含该计时器了。
而此时timer.C中并没有数据，接下来用于drain channel的代码会将consumer goroutine阻塞住。

这种情况，我们应该直接Reset计时器，而不用显式drain channel。如何将这两种情形合二为一呢？
我们可以利用一个select来包裹drain channel的操作，这样无论channel中是否有数据，drain都不会阻塞住。

//consumer
go func() {
	// try to read from channel, block at most 5s.
	// if timeout, print time event and go on loop.
	// if read a message which is not the type we want(we want true, not false),
	// retry to read.
	timer := time.NewTimer(time.Second * 5)
	for {
		// timer may be not active, and fired
		if !timer.Stop() {
			select {
			case <-timer.C: //try to drain from the channel
			default:
			}
		}
		timer.Reset(time.Second * 5)
		select {
		case b := <-c:
			if b == false {
				fmt.Println(time.Now(), ":recv false. continue")
				continue
			}
			//we want true, not false
			fmt.Println(time.Now(), ":recv true. return")
			return
		case <-timer.C:
			fmt.Println(time.Now(), ":timer expired")
			continue
		}
	}
}()

运行修改后的程序，发现程序不会被阻塞住，能正常进行通道读取，读取到true值后会自行退出。输出结果如下：
2022-05-29 18:04:01.517911 +0800 CST m=+5.000452379 :timer expired
2022-05-29 18:04:02.517841 +0800 CST m=+6.000373914 :recv false. continue
2022-05-29 18:04:07.518138 +0800 CST m=+11.000630379 :timer expired
2022-05-29 18:04:08.518351 +0800 CST m=+12.000835110 :recv false. continue
2022-05-29 18:04:13.518558 +0800 CST m=+17.001000990 :timer expired
2022-05-29 18:04:14.519554 +0800 CST m=+18.001989831 :recv false. continue
2022-05-29 18:04:19.520552 +0800 CST m=+23.002946404 :timer expired
2022-05-29 18:04:20.520706 +0800 CST m=+24.003092144 :recv false. continue
2022-05-29 18:04:25.521995 +0800 CST m=+29.004341058 :timer expired
2022-05-29 18:04:26.521351 +0800 CST m=+30.003688886 :recv false. continue
2022-05-29 18:04:31.522568 +0800 CST m=+35.004864964 :timer expired
2022-05-29 18:04:32.521519 +0800 CST m=+36.003808551 :recv true. return

》总结

Timer和Ticker都是在运行时计时器runtime.timer的基础上实现的。

运行时里的所有计时器都由运行时内唯一的timerproc触发。

time.Tick创建的Ticker在运行时不会被gc回收，能不用就不用。

Timer和Ticker的时间channel都是带有一个缓冲的通道。

time.After，time.NewTimer，time.NewTicker创建的计时器触发时都会执行sendTime。

sendTime和计时器带缓冲的时间通道保证了计时器不会阻塞程序。

Reset计时器时要注意drain channel和计时器过期存在竞争条件。

*/
