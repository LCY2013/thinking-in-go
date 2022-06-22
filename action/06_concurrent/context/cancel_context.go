package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

/*
》context取消goroutine执行的方法

》为什么需要取消功能
简单来说，我们需要取消功能来防止系统做一些不必要的工作。

考虑以下常见的场景：一个HTTP服务器查询数据库并将查询到的数据作为响应返回给客户端：
HTTP服务器查询数据库.png

如果一切正常，时序图将如下所示：

时序图.png

但是，如果客户端在中途取消了请求会发生什么？
这种情况可以发生在，比如用户在请求中途关闭了浏览器。
如果不支持取消功能，HTTP服务器和数据库会继续工作，由于客户端已经关闭所以他们工作的成果也就被浪费了。
这种情况的时序图如下所示：

取消时序图.png

理想情况下，如果我们知道某个处理过程（在此示例中为HTTP请求）已停止，则希望该过程的所有下游组件都停止运行：

理想时序图.png

》使用context实现取消功能

现在我们知道了应用程序为什么需要取消功能，接下来我们开始探究在Go中如何实现它。因为“取消事件”与正在执行的操作高度相关，因此很自然地会将它与上下文捆绑在一起。

取消功能需要从两方面实现才能完成：

1、监听取消事件

2、发出取消事件

》》监听取消事件

Go语言context标准库的Context类型提供了一个Done()方法，该方法返回一个类型为<-chan struct{}的channel。
每次context收到取消事件后这个channel都会接收到一个struct{}类型的值。
所以在Go语言里监听取消事件就是等待接收<-ctx.Done()。

举例来说，假设一个HTTP服务器需要花费两秒钟来处理一个请求。
如果在处理完成之前请求被取消，我们想让程序能立即中断不再继续执行下去：
*/
//func main() {
func main_() {
	// 创建一个监听8080端口的服务器
	err := http.ListenAndServe(":8080", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		// 输出到STDOUT展示处理已经开始
		_, err := fmt.Fprint(os.Stdout, "processing request\n")
		if err != nil {
			return
		}
		// 通过select监听多个channel
		select {
		case <-time.After(2 * time.Second):
			// 如果两秒后接受到了一个消息后，意味请求已经处理完成
			// 我们写入"request processed"作为响应
			_, err = writer.Write([]byte("request processed"))
			if err != nil {
				return
			}
		case res := <-ctx.Done():
			// 如果处理完成前取消了，在STDERR中记录请求被取消的消息
			_, err = fmt.Fprintf(os.Stderr, "request cancelled: %s\n", res)
			if err != nil {
				return
			}
		}
	}))
	if err != nil {
		log.Fatal(err)
		return
	}
}

/*
你可以通过运行服务器并在浏览器中打开localhost:8000进行测试。
如果你在2秒钟前关闭浏览器，则应该在终端窗口上看到“request cancelled”字样。

》》发出取消事件
如果你有一个可以取消的操作，则必须通过context发出取消事件。
可以通过context包的WithCancel函数返回的取消函数来完成此操作（withCancel还会返回一个支持取消功能的上下文对象）。
该函数不接受参数也不返回任何内容，当需要取消上下文时会调用该函数，发出取消事件。

考虑有两个相互依赖的操作的情况。在这里，“依赖”是指如果其中一个失败，那么另一个就没有意义，而不是第二个操作依赖第一个操作的结果（那种情况下，两个操作不能并行）。
在这种情况下，如果我们很早就知道其中一个操作失败，那么我们就会希望能取消所有相关的操作。

如果你有一个可以取消的操作，则必须通过context发出取消事件。
可以通过context包的WithCancel函数返回的取消函数来完成此操作（withCancel还会返回一个支持取消功能的上下文对象）。
该函数不接受参数也不返回任何内容，当需要取消上下文时会调用该函数，发出取消事件。

考虑有两个相互依赖的操作的情况。
在这里，“依赖”是指如果其中一个失败，那么另一个就没有意义，而不是第二个操作依赖第一个操作的结果（那种情况下，两个操作不能并行）。
在这种情况下，如果我们很早就知道其中一个操作失败，那么我们就会希望能取消所有相关的操作。
*/

func operation1(ctx context.Context) error {
	// 让我们假设这个操作会因为某种原因失败
	// 我们使用time.Sleep来模拟一个资源密集型操作
	time.Sleep(100 * time.Millisecond)
	return errors.New("failed")
}

func operation2(ctx context.Context) {
	// 我们使用在前面HTTP服务器例子里使用过的类似模式
	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("halted operation2")
	}
}

func contextCancel() {
	// 新建一个上下文
	ctx := context.Background()
	// 在初始上下文的基础上创建一个有取消功能的上下文
	ctx, cancel := context.WithCancel(ctx)
	// 在不同的goroutine中运行operation2
	go func() {
		operation2(ctx)
	}()

	err := operation1(ctx)
	// 如果这个操作返回错误，取消所有使用相同上下文的操作
	if err != nil {
		cancel()
	}
}

/*
》基于时间的取消

任何需要在请求的最大持续时间内维持SLA（服务水平协议）的应用程序，都应使用基于时间的取消。
该API与前面的示例几乎相同，但有一些补充：

// 这个上下文将会在3秒后被取消
// 如果需要在到期前就取消可以像前面的例子那样使用cancel函数
ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

// 上下文将在2009-11-10 23:00:00被取消
ctx, cancel := context.WithDeadline(ctx, time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))

例如，程序在对外部服务进行HTTP API调用时设置超时时间。如果被调用服务花费的时间太长，到时间后就会取消请求：
*/

func timeoutHTTP() {
	// 创建一个超时时间为100毫秒的上下文
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)

	// 创建一个访问Google主页的请求
	req, _ := http.NewRequest(http.MethodGet, "https://www.baidu.com", nil)
	// 将超时上下文关联到创建的请求上
	req = req.WithContext(ctx)

	// 创建一个HTTP客户端并执行请求
	client := &http.Client{}
	res, err := client.Do(req)
	// 如果请求失败了，记录到STDOUT
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	// 请求成功后打印状态码
	fmt.Println("Response received, status code:", res.StatusCode)
}

/*
》context使用上的一些陷阱

尽管Go中的上下文取消功能是一种多功能工具，但是在继续操作之前，你需要牢记一些注意事项。
其中最重要的是，上下文只能被取消一次。如果您想在同一操作中传播多个错误，那么使用上下文取消可能不是最佳选择。
使用取消上下文的场景是你实际上确实要取消某项操作，而不仅仅是通知下游进程发生了错误。
还需要记住的另一件事是，应该将相同的上下文实例传递给你可能要取消的所有函数和goroutine。

用WithTimeout或WithCancel包装一个已经支持取消功能的上下文将会造成多种可能会导致你的上下文被取消的情况，应该避免这种二次包装。
*/
