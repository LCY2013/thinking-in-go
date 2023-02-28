package main

/*
#include <stdlib.h>

void* makeSlice(size_t memseiz) {
	return malloc(memseiz);
}
*/
import "C"
import "unsafe"

/*
C 语言空间的内存是稳定的，只要不是被人为提前释放，那么在 Go 语言空间可以放心大胆地使用。
在 Go 语言访问 C 语言内存是最简单的情形，我们在之前的例子中已经见过多次。

因为 Go 语言实现的限制，我们无法在 Go 语言中创建大于 2GB 内存的切片（具体请参考 makeslice 实现代码）。
不过借助 cgo 技术，我们可以在 C 语言环境创建大于 2GB 的内存，然后转为 Go 语言的切片使用

我们通过 makeByteSlice 来创建大于 4G 内存大小的切片，从而绕过了 Go 语言实现的限制（需要代码验证）。
而 freeByteSlice 辅助函数则用于释放从 C 语言函数创建的切片。

因为 C 语言内存空间是稳定的，基于 C 语言内存构造的切片也是绝对稳定的，不会因为 Go 语言栈的变化而被移动。
*/

func makeByteSlice(n int) []byte {
	p := C.makeSlice(C.size_t(n))
	//return ((*[1 << 31]byte)(p))[0:n:n] // 4G
	return ((*[1 << 10]byte)(p))[0:n:n]
}

func freeByteSlice(p []byte) {
	C.free(unsafe.Pointer(&p[0]))
}

func main() {
	//s := makeByteSlice(1 << 31)
	s := makeByteSlice(1 << 10)
	s[len(s)-1] = 255
	println(s[len(s)-1])
	freeByteSlice(s)
}
