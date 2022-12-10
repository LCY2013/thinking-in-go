/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-06-10
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package interfaces_methods

import (
	"fmt"
	"net/http"
	"os"
)

/*
Interfaces and methods
接口和方法

由于几乎任何类型都能添加方法，因此几乎任何类型都能满足一个接口。一个很直观的例子 就是 http 包中定义的 Handler 接口。
任何实现了 Handler 的对象都能够处理 HTTP 请求。
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

ResponseWriter 接口提供了对方法的访问，这些方法需要响应客户端的请求。
由于这些方法 包含了标准的 Write 方法，因此 http.ResponseWriter 可用于任何 io.Writer 适用的场景。
Request 结构体包含已解析的客户端请求。

为简单起见，我们假设所有的 HTTP 请求都是 GET 方法，而忽略 POST 方法， 这种简化不会影响处理程序的建立方式。
这里有个短小却完整的处理程序实现， 它用于记录某个页面被访问的次数。
*/

// Counter 简单的计数器服务
type Counter struct {
	n int
}

// ServeHTTP 计数器服务
func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctr.n++
	_, err := fmt.Fprintf(w, "counter = %d\n", ctr.n)
	if err != nil {
		return
	}
}

/*
紧跟我们的主题，注意 Fprintf 如何能输出到 http.ResponseWriter。) 作为参考，这里演示了如何将这样一个服务器添加到 URL 树的一个节点上。
*/
func handler() {
	ctr := new(Counter)
	http.Handle("/counter", ctr)
}

//CounterInt 计数器
/*
但为什么 Counter 要是结构体呢?一个整数就够了。
An integer is all that's needed. (接收者必须为指针，增量操作对于调用者才可见。)
*/
type CounterInt int

func (ctr *CounterInt) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	*ctr++
	_, err := fmt.Fprintf(w, "counter = %d\n", *ctr)
	if err != nil {
		return
	}
}

func handlerCounterInt() {
	ctr := new(CounterInt)
	http.Handle("/counter/int", ctr)
}

//Chan
/*
当页面被访问时，怎样通知你的程序去更新一些内部状态呢?为 Web 页面绑定个信道吧。
*/
// 每次浏览时信道会发一个提醒
// (可能需要带缓存的信道)
type Chan chan *http.Request

func (ch Chan) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ch <- req
	_, err := fmt.Fprint(w, "notification sent")
	if err != nil {
		return
	}
}

// ArgServer 最后，假设我们需要输出调用服务器二进制程序时使用的实参 /args。 很简单，写个打印实参 的函数就行了。
//func ArgServer() {
//	fmt.Println(os.Args)
//}

/*
我们如何将它转换为 HTTP 服务器呢?
我们可以将 ArgServer 实现为某种可忽略值的方法，不过还有种更简单的方法。
既然我们可以为除指针和接口以外的任何类型定义方法，同样也能为一个函数写一个方法。
http包 包含以下代码如下所示、
*/

// HandlerFunc 类型是一个适配器，它允许将普通函数用做 HTTP 处理程序。
// 若 f 是个具有适当签名的函数，HandlerFunc(f) 就是个调用 f 的处理程序对象。
//type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(c, req).
//func (f HandlerFunc) ServeHTTP(w ResponseWriter, req *Request) {
//	f(w, req)
//}

// HandlerFunc 定义函数签名
type HandlerFunc func(int, *interface{})

func (f HandlerFunc) HandlerInteger(num int, obj *interface{}) {
	f(num, obj)
}

/*
HandlerFunc 是个具有 ServeHTTP 方法的类型， 因此该类型的值就能处理 HTTP 请求。
我们来看看该方法的实现:接收者是一个函数 f，而该方法调用 f。
这看起来很奇怪，但不必大惊小怪， 区别在于接收者变成了一个信道，而方法通过该信道发送消息。
*/

// ArgServer 为了将 ArgServer 实现成 HTTP 服务器，首先我们得让它拥有合适的签名。
// 实参服务器
func ArgServer(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintln(w, os.Args)
	if err != nil {
		return
	}
}

// StartHTTP 启动http服务
/*
ArgServer 和 HandlerFunc 现在拥有了相同的签名，
因此我们可将其转换为这种类型以访问它的方法，就像我们将 Sequence 转换为 IntSlice 以访问 IntSlice.Sort 那样。
建立代码非常简单:
*/
func StartHTTP() {
	http.Handle("/args", http.HandlerFunc(ArgServer))
}

/*
当有人访问 /args 页面时，该页面的处理程序就有了值 ArgServer 和类型 HandlerFunc。
HTTP 服务器会以 ArgServer 为接收者，调用该类型的 ServeHTTP 方法，
它会反过来调用 ArgServer(通过 f(c, req))，接着实参就会被显示出来。

通过一个结构体，一个整数，一个信道和一个函数，建立了一个 HTTP 服务器，这一切都是因为接口只是方法的集和，而几乎任何类型都能定义方法。
*/
