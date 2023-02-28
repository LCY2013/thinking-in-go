package main

/*
#include <stdlib.h>
#include <stdio.h>

void printString(const char* s) {
	printf("%s", s);
}
*/
import "C"
import "unsafe"

/*
C 临时访问传入的 Go 内存

cgo 之所以存在的一大因素是为了方便在 Go 语言中接纳吸收过去几十年来使用 C/C++ 语言软件构建的大量的软件资源。
C/C++ 很多库都是需要通过指针直接处理传入的内存数据的，因此 cgo 中也有很多需要将 Go 内存传入 C 语言函数的应用场景。

假设一个极端场景：我们将一块位于某 goroutine 的栈上的 Go 语言内存传入了 C 语言函数后，在此 C 语言函数执行期间，
此 goroutinue 的栈因为空间不足的原因发生了扩展，也就是导致了原来的 Go 语言内存被移动到了新的位置。
但是此时此刻 C 语言函数并不知道该 Go 语言内存已经移动了位置，仍然用之前的地址来操作该内存——这将将导致内存越界。
以上是一个推论（真实情况有些差异），也就是说 C 访问传入的 Go 内存可能是不安全的！

当然有 RPC 远程过程调用的经验的用户可能会考虑通过完全传值的方式处理：
借助 C 语言内存稳定的特性，在 C 语言空间先开辟同样大小的内存，然后将 Go 的内存填充到 C 的内存空间；
返回的内存也是如此处理。下面的例子是这种思路的具体实现：
*/

func printString(s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	C.printString(cs)
}

func main() {
	s := "hello, world!"
	printString(s)
}
