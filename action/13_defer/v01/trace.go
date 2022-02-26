package main

func Trace(name string) func() {
	println("enter:", name)
	return func() {
		println("exit:", name)
	}
}

// foo 以 foo 函数中的defer Trace("foo")()这行代码为例，Go 会对 defer 后面的表达式Trace("foo")()进行求值。
// 由于这个表达式包含一个函数调用Trace("foo")，所以这个函数会被执行。
func foo() {
	defer Trace("foo")()
	bar()
}

func bar() {
	defer Trace("bar")()
}

// main Go 会在 defer 设置 deferred 函数时对 defer 后面的表达式进行求值。
// 这里还是有一些“瑕疵”，也就是离我们期望的“跟踪函数调用链”的实现还有一些不足之处。这里列举了几点：
// 1、调用 Trace 时需手动显式传入要跟踪的函数名；
// 2、如果是并发应用，不同 Goroutine 中函数链跟踪混在一起无法分辨；
// 3、输出的跟踪结果缺少层次感，调用关系不易识别；
// 4、对要跟踪的函数，需手动调用 Trace 函数。
// 结论：最终实现一个自动注入跟踪代码，并输出有层次感的函数调用链跟踪命令行工具。
func main() {
	defer Trace("main")()
	foo()
}
