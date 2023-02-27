package main

/*
#include <stdio.h>

static void SayHello(const char* s) {
    puts(s);
}
*/
import "C"

/*
使用自己的 C 函数

先自定义一个叫 SayHello 的 C 函数来实现打印，然后从 Go 语言环境中调用这个 SayHello 函数。

可以将 SayHello 函数放到当前目录下的一个 C 语言源文件中（后缀名必须是 .c）。
因为是编写在独立的 C 文件中，为了允许外部引用，所以需要去掉函数的 static 修饰符。
*/
func main() {
	C.SayHello(C.CString("hello, world!\n"))
}
