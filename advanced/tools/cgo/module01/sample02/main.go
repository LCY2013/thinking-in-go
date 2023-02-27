package main

//#include <stdio.h>
import "C"

/*
不仅仅通过 import "C" 语句启用 CGO 特性，同时包含 C 语言的 <stdio.h> 头文件。
然后通过 CGO 包的 C.CString 函数将 Go 语言字符串转为 C 语言字符串，
最后调用 CGO 包的 C.puts 函数向标准输出窗口打印转换后的 C 字符串。

相比 sample01中 的 CGO 程序最大的不同是：没有在程序退出前释放 C.CString 创建的 C 语言字符串；
还有我们改用 puts 函数直接向标准输出打印，之前是采用 fputs 向标准输出打印。

没有释放使用 C.CString 创建的 C 语言字符串会导致内存泄漏。
但是对于这个小程序来说，这样是没有问题的，因为程序退出后操作系统会自动回收程序的所有资源。
*/
func main() {
	//println("hello cgo")
	C.puts(C.CString("hello, world!\n"))
}
