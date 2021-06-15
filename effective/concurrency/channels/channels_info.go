/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-06-15
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package channels

import (
	"fmt"
	"net/http"
	"os"
)

/*
Channels 信道

信道与映射一样，也需要通过 make 来分配内存。
其结果值充当了对底层数据结构的引用。
若提供了一个可选的整数形参，它就会为该信道设置缓冲区大小。
默认值是零，表示不带缓 冲的或同步的信道。
ci := make(chan int)            // unbuffered channel of integers
cj := make(chan int, 0)         // unbuffered channel of integers
cs := make(chan *os.File, 100)  // buffered channel of pointers to Files

无缓冲信道在通信时会同步交换数据，它能确保(两个 Go 程的)计算处于确定状态。

信道有很多惯用法，我们从这里开始了解。在上一节中，我们在后台启动了排序操作。 信道使得启动的 Go 程等待排序完成。
c := make(chan int) // 分配一个信道
// 在 Go 程中启动排序。当它完成后，在信道上发送信号。
go func() {
	list.Sort()
	c <- 1 // 发送信号，什么值无所谓。
}()
doSomethingForAWhile()
<-c // 等待排序结束，丢弃发来的值。

接收者在收到数据前会一直阻塞。
若信道是不带缓冲的，那么在接收者收到值前， 发送者会一直阻塞;
若信道是带缓冲的，则发送者仅在值被复制到缓冲区前阻塞;
若缓冲区已满，发送者会一直等待直到某个接收者取出一个值为止。

带缓冲的信道可被用作信号量，例如限制吞吐量。
在此例中，进入的请求会被传递给 handle，它从信道中接收值，处理请求后将值发回该信道中，以便让该 “信号量” 准备迎接下 一次请求。
信道缓冲区的容量决定了同时调用 process 的数量上限，因此我们在初始化时首 先要填充至它的容量上限。

var sem = make(chan int, MaxOutstanding)
func handle(r *Request) {
	sem <- 1 // 等待活动队列清空。
	process(r) // 可能需要很长时间。
	<-sem // 完成;使下一个请求可以运行。
}
func Serve(queue chan *Request) {
    for {
        req := <-queue
 		go handle(req) // 无需等待 handle 结束。
    }
}

由于数据同步发生在信道的接收端(也就是说发送发生在 > 接受之前，参见 Go 内存模型 https://go-zh.org/ref/mem )，因此信号必须在信道的接收端获取，而非发送端。

然而，它却有个设计问题:尽管只有 MaxOutstanding 个 Go 程能同时运行，但 Serve 还是为每个进入的请求都创建了新的 Go 程。
其结果就是，若请求来得很快， 该程序就会无限地消耗资源。
为了弥补这种不足，我们可以通过修改 Serve 来限制创建 Go 程，这是个明显的解决方案，但要当心我们修复后出现的 Bug。

func Serve(queue chan *Request) {
    for req := range queue {
		sem <- 1
		go func(){
			process(req) // 这里存在bug
			<- sem
	    }()
    }
}

Bug 出现在 Go 的 for 循环中，该循环变量在每次迭代时会被重用，
因此 req 变量会在所有的 Go 程间共享，这不是我们想要的。
我们需要确保 req 对于每个 Go 程来说都是唯一的。
有一 种方法能够做到，就是将 req 的值作为实参传入到该 Go 程的闭包中:
func Serve(queue chan *Request){
	for req := range queue {
		sem <- 1
		go func(req *Request){
			process(req)
			<- sem
		}(req)
	}
}

比较前后两个版本，观察该闭包声明和运行中的差别。 另一种解决方案就是以相同的名字创建新的变量，如例中所示:
func Serve(queue chan *Request){
	for req := range queue {
		req := req // 为该GO 程创建 req 的新实例
		sem <- 1
		go func(){

		}()
	}
}

它的写法看起来有点奇怪
req := req
但在 Go 中这样做是合法且惯用的。
你用相同的名字获得了该变量的一个新的版本， 以此来局部地刻意屏蔽循环变量，使它对每个 Go 程保持唯一。

回到编写服务器的一般问题上来。
另一种管理资源的好方法就是启动固定数量的 handle Go 程，一起从请求信道中读取数据。
Go 程的数量限制了同时调用 process 的数量。
Serve 同样会接收一个通知退出的信道， 在启动所有 Go 程后，它将阻塞并暂停从信道中接收消息。

func handle(queue chan *Request) {
    for r := range queue {
		process(r)
	}
}
func Serve(clientRequests chan *Request, quit chan bool) {
	// 启动处理程序
    for i := 0; i < MaxOutstanding; i++ {
        go handle(clientRequests)
	}
	<-quit // 等待通知退出。
}

*/

func channel() {
	// 整数类型的无缓冲信道
	cj := make(chan int)
	_ = cj

	// 整数类型的无缓冲信道
	cj1 := make(chan int, 0)
	_ = cj1

	// 指向文件指针的带缓冲信道
	cj2 := make(chan *os.File, 10)
	_ = cj2
}

func ChannelCommunication() {
	// 分配一个信道
	c := make(chan int)
	// 在go协程中启动排序，当它完成后，在信道上发送信号。
	go func() {
		// 排序
		sort()
		c <- 1
	}()
	<-c
	doSomethingForAWhile()
}

func sort() {
	fmt.Println("排序...")
}

func doSomethingForAWhile() {
	fmt.Println("排序完成...")
}

var sem = make(chan int, 100)

func handle(r *http.Request) {
	sem <- 1 // 等待活动队列清空
	//process(r) // 可能需要很长时间
	<-sem // 完成，使下一个请求可以运行
}

func Serve(queue chan *http.Request) {
	for {
		req := <-queue
		go handle(req) // 无需等待 handle 结束
	}
}

func Serve1(queue chan *http.Request) {
	for req := range queue {
		sem <- 1
		go func(req *http.Request) {
			// process(req)
			<-sem
		}(req)
	}
}

func Serve2(queue chan *http.Request) {
	for req := range queue {
		req := req // 为该GO 程创建 req 的新实例
		sem <- 1
		go func() {
			// process(req)
			_ = req
			<-sem
		}()
	}
}

// 固定处理数

func handle3(queue chan *http.Request) {
	for r := range queue {
		// process(r)
		_ = r
	}
}

func Serve3(clientRequests chan *http.Request, quit chan bool) {
	// 启动处理程序
	for i := 0; i < 100; i++ {
		go handle3(clientRequests)
	}
	<-quit // 等待通知退出
}
