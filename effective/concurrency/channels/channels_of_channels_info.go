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
)

/*
Channels of channels  信道中的信道

Go 最重要的特性就是信道是一等值，它可以被分配并像其它值到处传递。 这种特性通常被用来实现安全、并行的多路分解。

在上一节的例子中，handle 是个非常理想化的请求处理程序，但我们并未定义它所处理的请求类型。
若该类型包含一个可用于回复的信道，那么每一个客户端都能为其回应提供自己的路径。以下为 Request 类型的大概定义。
type Request struct {
	args       []int
	f          func([]int) int
	resultChan chan int
}

客户端提供了一个函数及其实参，此外在请求对象中还有个接收应答的信道。
func sum(a []int) (s int) {
    for _, v := range a {
		s += v
	}
	return
}
request := &Request{[]int{3, 4, 5}, sum, make(chan int)}
// Send request
clientRequests <- request
// Wait for response.
fmt.Printf("answer: %d\n", <-request.resultChan)

func handle(queue chan *Request) {
    for req := range queue {
        req.resultChan <- req.f(req.args)
    }
}

要使其实际可用还有很多工作要做，这些代码仅能实现一个速率有限、并行、非阻塞 RPC 系统的框架，而且它并不包含互斥锁。
*/

type Request struct {
	args       []int
	f          func([]int) int
	resultChan chan int
}

func sum(a []int) (s int) {
	for _, v := range a {
		s += v
	}
	return
}

func handler(queue chan *Request) {
	for req := range queue {
		req.resultChan <- req.f(req.args)
	}
}

func Server(clientRequests chan *Request, quit chan bool) {
	// 启动处理程序
	for i := 0; i < 100; i++ {
		go handler(clientRequests)
	}
	<-quit // 等待通知退出
}

func Start() {
	request := &Request{[]int{3, 4, 5}, sum, make(chan int)}
	// Send request
	// clientRequests <- request
	chanRequest := make(chan *Request)
	Server(chanRequest, make(chan bool))
	// 等待响应
	fmt.Printf("answer: %d\n", <-request.resultChan)
}
