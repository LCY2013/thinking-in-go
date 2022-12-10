package main

import (
	"unsafe"

	"github.com/lcy2013/dumpinterface/dump"
)

// 输出接口类型变量内部表示的详细信息
/*
println 输出的接口类型变量的内部表示信息，在一般情况下都是足够的，但有些时候又显得过于简略，
比如在v06中，如果仅凭eif: (0x10b3b00,0x10eb4d0)和err: (0x10ed380,0x10eb4d8)的输出，是无法想到两个变量是相等的。

那这时如果我们能输出接口类型变量内部表示的详细信息（比如：tab._type），那势必可以取得事半功倍的效果。

eface 和 iface 以及组成它们的 itab 和 _type 都是 runtime 包下的非导出结构体，我们无法在外部直接引用它们。
但我们发现，组成 eface、iface 的类型都是基本数据类型，完全可以通过“复制代码”的方式将它们拿到 runtime 包外面来。

不过，这里要注意，由于 runtime 中的 eface、iface，或者它们的组成可能会随着 Go 版本的变化发生变化，因此这个方法不具备跨版本兼容性。
也就是说，基于 Go 1.17 版本复制的代码，可能仅适用于使用 Go 1.17 版本编译。这里我们就以 Go 1.17 版本为例看看：
dumpinterface.go
*/

func main() {
	var eif interface{} = dump.T(5)
	var err error = dump.T(5)
	println("eif:", eif)
	println("err:", err)
	println("eif = err:", eif == err)

	dump.DumpEface(eif)
	dump.DumpItabOfIface(unsafe.Pointer(&err))
	dump.DumpDataOfIface(err)
}

/*
eif: (0x1095a80,0x10c1e90)
err: (0x10c2740,0x10c1e90)
eif = err: true
eface: {_type:0x1095a80 data:0x10c1e90}
         _type: {size:8 ptrdata:0 hash:664027975 tflag:15 align:8 fieldAlign:8 kind:2 equal:0x1002c80 gcdata:0x10c1dc3 str:3532 ptrToThis:40256}
         data: bad error

iface: {tab:0x10c2740 data:0x10c1e90}
         itab: {inter:0x10979e0 _type:0x1095a80 hash:664027975 _:[0 0 0 0] fun:[17350560]}
                 inter: {typ:{size:16 ptrdata:16 hash:235953867 tflag:7 align:8 fieldAlign:8 kind:20 equal:0x1002de0 gcdata:0x10ac050 str:2568 ptrToThis:17536} pkgpath:{bytes:<nil>} mhdr:[{name:1458 ityp:29568}]}
                 _type: {size:8 ptrdata:0 hash:664027975 tflag:15 align:8 fieldAlign:8 kind:2 equal:0x1002c80 gcdata:0x10c1dc3 str:3532 ptrToThis:40256}
                 fun: [0x108bfa0(17350560),]
         data: bad error

从输出结果中，我们看到 eif 的 _type（0x10b38c0）与 err 的 tab._type（0x10b38c0）是一致的，data 指针所指内容（“bad error”）也是一致的，因此eif == err表达式的结果为 true。
再次强调一遍，上面这个实现可能仅在 Go 1.17 版本上测试通过，并且在输出 iface 或 eface 的 data 部分内容时只列出了 int、float64 和 T 类型的数据读取实现，没有列出全部类型的实现，你可以根据自己的需要实现其余数据类型。
*/
