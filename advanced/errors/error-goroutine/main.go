package main

import "log"

/*
在程序启动的时候，如果有强依赖的服务出现故障时 panic 退出
在程序启动的时候，如果发现有配置明显不符合要求， 可以 panic 退出（防御编程）
其他情况下只要不是不可恢复的程序错误，都不应该直接 panic 应该返回 error
在程序入口处，例如 gin 中间件需要使用 recover 预防 panic 程序退出

在程序中我们应该避免使用野生的 goroutine
1、如果是在请求中需要执行异步任务，应该使用异步 worker ，消息通知的方式进行处理，避免请求量大时大量 goroutine 创建
2、如果需要使用 goroutine 时，应该使用同一的 Go 函数进行创建，这个函数中会进行 recover ，避免因为野生 goroutine panic 导致主进程退出
*/
func main() {

}

// GO 全局声明的异步 - 执行异步
func GO(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v\n", err)
			}
		}()

		f()
	}()
}

/*
error：
1、我们在应用程序中使用 github.com/pkg/errors 处理应用错误，注意在公共库当中，我们一般不使用这个
2、error 应该是函数的最后一个返回值，当 error 不为 nil 时，函数的其他返回值是不可用的状态，不应该对其他返回值做任何期待
	func f() (io.Reader, *S1, error) 在这里，我们不知道 io.Reader 中是否有数据，可能有，也有可能有一部分
3、错误处理的时候应该先判断错误， if err != nil 出现错误及时返回，使代码是一条流畅的直线，避免过多的嵌套.
*/
