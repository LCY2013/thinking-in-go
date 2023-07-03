package main

// 利用 instrument 工具注入跟踪代码
// 有了 instrument 工具后，再来看看如何使用这个工具，在目标 Go 源文件中自动注入跟踪设施。
/*
$cd instrument_trace
$go build github.com/lcy2013/instrument_trace/cmd/instrument
$instrument version
[instrument version]
instrument [-w] xxx.go
  -w  write result to (source) file instead of stdout
*/
func foo() {
	//defer trace.Trace()()
	bar()
}

func bar() {
	//defer trace.Trace()()
}

func main() {
	//defer trace.Trace()()
	foo()
}
