package main

import (
	"fmt"
	"fufeng.org/advanced/errors/sentinel_error/v1"
)

/**
调用者要使用类型断言和类型 switch，就要让自定义的 error 变为 public。这种模型会导致和调用者产生强耦合，从而导致 API 变得脆弱。
结论是尽量避免使用 error types，虽然错误类型比 sentinel errors 更好，因为它们可以捕获关于出错的更多上下文，但是 error types 共享 error values 许多相同的问题。
因此，我的建议是避免错误类型，或者至少避免将它们作为公共 API 的一部分。

*/

func main() {
	err := v1.Test()
	// 因为 MyError 是一个type ，调用者可以使用断言转换成这个类型，来获取更多的上下文信息
	switch err := err.(type) {
	case nil:
		// call succeeded , nothing to do
	case *v1.MyError:
		fmt.Println("error occurred on line:", err.Line)
	default:
		// unknown error
	}
}
