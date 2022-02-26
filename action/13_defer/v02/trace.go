package main

import "runtime"

// 1、自动获取所跟踪函数的函数名
// 要解决“调用 Trace 时需要手动显式传入要跟踪的函数名”的问题，也就是要让我们的 Trace 函数能够自动获取到它跟踪函数的函数名信息。
// 以跟踪 foo 为例，看看这样做能带来什么好处。

// 在手动显式传入的情况下，需要用下面这个代码对 foo 进行跟踪：
// defer Trace("foo")()

// 实现了自动获取函数名，所有支持函数调用链跟踪的函数都只需使用下面调用形式的 Trace 函数就可以了
// defer Trace()()

// Trace 这种一致的 Trace 函数调用方式也为后续的自动向代码中注入 Trace 函数奠定了基础。那么如何实现 Trace 函数对它跟踪函数名的自动获取呢？需要借助 Go 标准库 runtime 包的帮助。
// 新版 Trace 函数的实现以及它的使用方法如下：
func Trace() func() {
	// 通过 runtime.Caller 函数获得当前 Goroutine 的函数调用栈上的信息
	// runtime.Caller 的参数标识的是要获取的是哪一个栈帧的信息。
	// 当参数为 0 时，返回的是 Caller 函数的调用者的函数信息，在这里就是 Trace 函数。
	// 但我们需要的是 Trace 函数的调用者的信息，于是我们传入 1。
	//
	// Caller 函数有四个返回值：
	// 第一个返回值代表的是程序计数（pc）；
	// 第二个和第三个参数代表对应函数所在的源文件名以及所在行数，这里我们暂时不需要；
	// 最后一个参数代表是否能成功获取这些信息，如果获取失败，抛出 panic。
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}
	// 通过 runtime.FuncForPC 函数和程序计数器（PC）得到被跟踪函数的函数名称。
	// runtime.FuncForPC 返回的名称中不仅仅包含函数名，还包含了被跟踪函数所在的包名。
	fn := runtime.FuncForPC(pc)
	name := fn.Name()
	println("enter:", name)
	return func() {
		println("exit:", name)
	}
}

func foo() {
	defer Trace()()
	bar()
}

func bar() {
	defer Trace()()
}

func main() {
	defer Trace()()
	foo()
}
