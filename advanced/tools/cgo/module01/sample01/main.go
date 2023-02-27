//go:build cgo

package main

import "C"

/*
import "C" 语句启用 CGO 特性，主函数只是通过 Go 内置的 println 函数输出字符串，其中并没有任何和 CGO 相关的代码。

虽然没有调用 CGO 的相关函数，但是 go build 命令会在编译和链接阶段启动 gcc 编译器，这已经是一个完整的 CGO 程序了。
*/
func main() {
	println("hello cgo")
}
