package trace_test

import (
	trace "github.com/lcy2013/instrument_trace"
)

func a() {
	defer trace.Trace()()
	b()
}

func b() {
	defer trace.Trace()()
	c()
}

func c() {
	defer trace.Trace()()
	d()
}

func d() {
	defer trace.Trace()()
}

// 在 example_test.go 文件中，我们用 ExampleXXX 形式的函数表示一个示例，go test 命令会扫描 example_test.go 中的以 Example 为前缀的函数并执行这些函数。
// 每个 ExampleXXX 函数需要包含预期的输出，就像上面 ExampleTrace 函数尾部那样，
// 在一大段注释中提供这个函数执行后的预期输出，预期输出的内容从// Output:的下一行开始。
// go test 会将 ExampleTrace 的输出与预期输出对比，如果不一致，会报测试错误。
// 从这一点，可以看出 example_test.go 也是 trace 包单元测试的一部分。
func ExampleTrace() {
	a()
	// Output:
	// g[00001]: 	->github.com/lcy2013/instrument_trace_test.a
	// g[00001]: 		->github.com/lcy2013/instrument_trace_test.b
	// g[00001]: 			->github.com/lcy2013/instrument_trace_test.c
	// g[00001]: 				->github.com/lcy2013/instrument_trace_test.d
	// g[00001]: 				<-github.com/lcy2013/instrument_trace_test.d
	// g[00001]: 			<-github.com/lcy2013/instrument_trace_test.c
	// g[00001]: 		<-github.com/lcy2013/instrument_trace_test.b
	// g[00001]: 	<-github.com/lcy2013/instrument_trace_test.a
}
