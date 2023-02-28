package main

/*
导出 C 函数不能返回 Go 内存

在 Go 语言中，Go 是从一个固定的虚拟地址空间分配内存。而 C 语言分配的内存则不能使用 Go 语言保留的虚拟内存空间。
在 CGO 环境，Go 语言运行时默认会检查导出返回的内存是否是由 Go 语言分配的，如果是则会抛出运行时异常。

下面是 CGO 运行时异常的例子：
/*
extern int* getGoPtr();

static void Main() {
    int* p = getGoPtr();
    *p = 42;
}
*
import "C"

func main() {
	C.Main()
}

//export getGoPtr
func getGoPtr() *C.int {
	return new(C.int)
}

其中 getGoPtr 返回的虽然是 C 语言类型的指针，但是内存本身是从 Go 语言的 new 函数分配，也就是由 Go 语言运行时统一管理的内存。
然后我们在 C 语言的 Main 函数中调用了 getGoPtr 函数，此时默认将发送运行时异常：

$ go run main.go
panic: runtime error: cgo result has Go pointer

goroutine 1 [running]:
main._cgoexpwrap_cfb3840e3af2_getGoPtr.func1(0xc420051dc0)
  command-line-arguments/_obj/_cgo_gotypes.go:60 +0x3a
main._cgoexpwrap_cfb3840e3af2_getGoPtr(0xc420016078)
  command-line-arguments/_obj/_cgo_gotypes.go:62 +0x67
main._Cfunc_Main()
  command-line-arguments/_obj/_cgo_gotypes.go:43 +0x41
main.main()

异常说明 cgo 函数返回的结果中含有 Go 语言分配的指针。指针的检查操作发生在 C 语言版的 getGoPtr 函数中，它是由 cgo 生成的桥接 C 语言和 Go 语言的函数。

下面是 cgo 生成的 C 语言版本 getGoPtr 函数的具体细节（在 cgo 生成的 _cgo_export.c 文件定义）：
int* getGoPtr()
{
    __SIZE_TYPE__ _cgo_ctxt = _cgo_wait_runtime_init_done();
    struct {
        int* r0;
    } __attribute__((__packed__)) a;
    _cgo_tsan_release();
    crosscall2(_cgoexp_95d42b8e6230_getGoPtr, &a, 8, _cgo_ctxt);
    _cgo_tsan_acquire();
    _cgo_release_context(_cgo_ctxt);
    return a.r0;
}

其中 _cgo_tsan_acquire 是从 LLVM 项目移植过来的内存指针扫描函数，它会检查 cgo 函数返回的结果是否包含 Go 指针。

需要说明的是，cgo 默认对返回结果的指针的检查是有代价的，特别是 cgo 函数返回的结果是一个复杂的数据结构时将花费更多的时间。
如果已经确保了 cgo 函数返回的结果是安全的话，可以通过设置环境变量 GODEBUG=cgocheck=0 来关闭指针检查行为。

$ GODEBUG=cgocheck=0 go run main.go

关闭 cgocheck 功能后再运行上面的代码就不会出现上面的异常的。
但是要注意的是，如果 C 语言使用期间对应的内存被 Go 运行时释放了，将会导致更严重的崩溃问题。
cgocheck 默认的值是 1，对应一个简化版本的检测，如果需要完整的检测功能可以将 cgocheck 设置为 2。

关于 cgo 运行时指针检测的功能详细说明可以参考 Go 语言的官方文档。
*/

func main() {

}
