package main

import "fmt"

/*
接口类型的装箱（boxing）原理
装箱（boxing）是编程语言领域的一个基础概念，一般是指把一个值类型转换成引用类型，比如在支持装箱概念的 Java 语言中，将一个 int 变量转换成 Integer 对象就是一个装箱操作。

在 Go 语言中，将任意类型赋值给一个接口类型变量也是装箱操作。面对接口类型变量内部表示，知道接口类型的装箱实际就是创建一个 eface 或 iface 的过程。

接下来就来简要描述一下这个过程，也就是接口类型的装箱原理。

*/
type T struct {
	n int
	s string
}

func (T) M1() {}
func (T) M2() {}

type NonEmptyInterface interface {
	M1()
	M2()
}

/*
这个例子中，对 ei 和 i 两个接口类型变量的赋值都会触发装箱操作，要想知道 Go 在背后做了些什么，需要“下沉”一层，也就是要输出上面 Go 代码对应的汇编代码：

$ go tool compile -S main.go > main.s

对应ei = t一行的汇编如下：
	0x0026 00038 (main.go:37)	MOVQ	$17, ""..autotmp_15+104(SP)
	0x002f 00047 (main.go:37)	LEAQ	go.string."hello, interface"(SB), CX
	0x0036 00054 (main.go:37)	MOVQ	CX, ""..autotmp_15+112(SP)
	0x003b 00059 (main.go:37)	MOVQ	$16, ""..autotmp_15+120(SP)
	0x0044 00068 (main.go:37)	LEAQ	type."".T(SB), AX
	0x004b 00075 (main.go:37)	LEAQ	""..autotmp_15+104(SP), BX
	0x0050 00080 (main.go:37)	PCDATA	$1, $0
	0x0050 00080 (main.go:37)	CALL	runtime.convT2E(SB)

对应 i = t 一行的汇编如下：
	0x005f 00095 (main.go:40)	MOVQ	$17, ""..autotmp_15+104(SP)
	0x0068 00104 (main.go:40)	LEAQ	go.string."hello, interface"(SB), CX
	0x006f 00111 (main.go:40)	MOVQ	CX, ""..autotmp_15+112(SP)
	0x0074 00116 (main.go:40)	MOVQ	$16, ""..autotmp_15+120(SP)
	0x007d 00125 (main.go:40)	LEAQ	go.itab."".T,"".NonEmptyInterface(SB), AX
	0x0084 00132 (main.go:40)	LEAQ	""..autotmp_15+104(SP), BX
	0x0089 00137 (main.go:40)	PCDATA	$1, $1
	0x0089 00137 (main.go:40)	CALL	runtime.convT2I(SB)

在将动态类型变量赋值给接口类型变量语句对应的汇编代码中，看到了convT2E和convT2I两个 runtime 包的函数。
这两个函数的实现位于$GOROOT/src/runtime/iface.go中：

func convT2E(t *_type, elem unsafe.Pointer) (e eface) {
	if raceenabled {
		raceReadObjectPC(t, elem, getcallerpc(), funcPC(convT2E))
	}
	if msanenabled {
		msanread(elem, t.size)
	}
	x := mallocgc(t.size, t, true)
	// TODO: We allocate a zeroed object only to overwrite it with actual data.
	// Figure out how to avoid zeroing. Also below in convT2Eslice, convT2I, convT2Islice.
	typedmemmove(t, x, elem)
	e._type = t
	e.data = x
	return
}

func convT2I(tab *itab, elem unsafe.Pointer) (i iface) {
	t := tab._type
	if raceenabled {
		raceReadObjectPC(t, elem, getcallerpc(), funcPC(convT2I))
	}
	if msanenabled {
		msanread(elem, t.size)
	}
	x := mallocgc(t.size, t, true)
	typedmemmove(t, x, elem)
	i.tab = tab
	i.data = x
	return
}

convT2E 用于将任意类型转换为一个 eface，convT2I 用于将任意类型转换为一个 iface。
两个函数的实现逻辑相似，主要思路就是根据传入的类型信息（convT2E 的 _type 和 convT2I 的 tab._type）分配一块内存空间，并将 elem 指向的数据拷贝到这块内存空间中，最后传入的类型信息作为返回值结构中的类型信息，返回值结构中的数据指针（data）指向新分配的那块内存空间。

由此也可以看出，经过装箱后，箱内的数据，也就是存放在新分配的内存空间中的数据与原变量便无瓜葛了，比如下面这个例子：
func main() {
  var n int = 61
  var ei interface{} = n
  n = 62  // n的值已经改变
  fmt.Println("data in box:", ei) // 输出仍是61
}

那么 convT2E 和 convT2I 函数的类型信息是从何而来的呢？

其实这些都依赖 Go 编译器的工作。
编译器知道每个要转换为接口类型变量（toType）和动态类型变量的类型（fromType），它会根据这一对类型选择适当的 convT2X 函数，并在生成代码时使用选出的 convT2X 函数参与装箱操作。

不过，装箱是一个有性能损耗的操作，因此 Go 也在不断对装箱操作进行优化，包括对常见类型如整型、字符串、切片等提供系列快速转换函数：
// $GOROOT/src/runtime/iface.go
func convT16(val any) unsafe.Pointer     // val must be uint16-like
func convT32(val any) unsafe.Pointer     // val must be uint32-like
func convT64(val any) unsafe.Pointer     // val must be uint64-like
func convTstring(val any) unsafe.Pointer // val must be a string
func convTslice(val any) unsafe.Pointer  // val must be a slice

这些函数去除了 typedmemmove 操作，增加了零值快速返回等特性。

同时 Go 建立了 staticuint64s 区域，对 255 以内的小整数值进行装箱操作时不再分配新内存，而是利用 staticuint64s 区域的内存空间，下面是 staticuint64s 的定义：
/ $GOROOT/src/runtime/iface.go
// staticuint64s is used to avoid allocating in convTx for small integer values.
var staticuint64s = [...]uint64{
    0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
    0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
  ... ...
}
*/
func main() {
	var t = T{
		n: 17,
		s: "hello, interface",
	}
	var ei interface{}
	ei = t

	var i NonEmptyInterface
	i = t
	fmt.Println(ei)
	fmt.Println(i)
}

/*
在 Go 语言中有着很高的地位。它这个地位的取得离不开它拥有的“动静兼备”的语法特性。
Go 接口的动态特性让 Go 拥有与动态语言相近的灵活性，而静态特性又在编译阶段保证了这种灵活性的安全。

要更好地理解 Go 接口的这两种特性，我们需要深入到 Go 接口在运行时的表示层面上去。
接口类型变量在运行时表示为 eface 和 iface，eface 用于表示空接口类型变量，iface 用于表示非空接口类型变量。
只有两个接口类型变量的类型信息（eface._type/iface.tab._type）相同，且数据指针（eface.data/iface.data）所指数据相同时，两个接口类型变量才是相等的。

我们可以通过 println 输出接口类型变量的两部分指针变量的值。
而且，通过拷贝 runtime 包 eface 和 iface 相关类型源码，我们还可以自定义输出 eface/iface 详尽信息的函数，不过要注意的是，由于 runtime 层代码的演进，这个函数可能不具备在 Go 版本间的移植性。

最后，接口类型变量的赋值本质上是一种装箱操作，装箱操作是由 Go 编译器和运行时共同完成的，有一定的性能开销，对于性能敏感的系统来说，我们应该尽量避免或减少这类装箱操作。
*/
