package main

import "fmt"

/*
由于 eface 和 iface 是 runtime 包中的非导出结构体定义，我们不能直接在包外使用，所以也就无法直接访问到两个结构体中的数据。
不过，Go 语言提供了 println 预定义函数，可以用来输出 eface 或 iface 的两个指针字段的值。

在编译阶段，编译器会根据要输出的参数的类型将 println 替换为特定的函数，这些函数都定义在$GOROOT/src/runtime/print.go文件中，而针对 eface 和 iface 类型的打印函数实现如下：
// $GOROOT/src/runtime/print.go
func printeface(e eface) {
    print("(", e._type, ",", e.data, ")")
}
func printiface(i iface) {
    print("(", i.tab, ",", i.data, ")")
}

printeface 和 printiface 会输出各自的两个指针字段的值。下面我们就来使用 println 函数输出各类接口类型变量的内部表示信息，并结合输出结果，解析接口类型变量的等值比较操作。
*/

// 第一种：nil 接口变量
// 未赋初值的接口类型变量的值为 nil，这类变量也就是 nil 接口变量，我们来看这类变量的内部表示输出的例子：
// (0x0,0x0)
// (0x0,0x0)
// i = nil: true
// err = nil: true
// i = err: true
// 无论是空接口类型还是非空接口类型变量，一旦变量值为 nil，那么它们内部表示均为(0x0,0x0)，也就是类型信息、数据值信息均为空。
// 因此下面的变量 i 和 err 等值判断为 true。
func printNilInterface() {
	// nil接口变量
	var i interface{} // 空接口类型
	var err error     // 非空接口类型
	println(i)
	println(err)
	println("i = nil:", i == nil)
	println("err = nil:", err == nil)
	println("i = err:", i == err)
}

// 第二种：空接口类型变量
/*
eif1: (0x1059040,0xc000042768)
eif2: (0x1059040,0xc000042760)
eif1 = eif2: false
eif1: (0x1059040,0xc000042768)
eif2: (0x1059040,0x1075e58)
eif1 = eif2: true
eif1: (0x1059040,0xc000042768)
eif2: (0x1059100,0x1075e58)
eif1 = eif2: false
首先，代码执行到第 11 行时，eif1 与 eif2 已经分别被赋值整型值 17 与 18，这样 eif1 和 eif2 的动态类型的类型信息是相同的（都是 0x10ac580），但 data 指针指向的内存块中存储的值不同，一个是 17，一个是 18，于是 eif1 不等于 eif2。
接着，代码执行到第 16 行的时候，eif2 已经被重新赋值为 17，这样 eif1 和 eif2 不仅存储的动态类型的类型信息是相同的（都是 0x10ac580），data 指针指向的内存块中存储值也相同了，都是 17，于是 eif1 等于 eif2。
然后，代码执行到第 21 行时，eif2 已经被重新赋值了 int64 类型的数值 17。这样，eif1 和 eif2 存储的动态类型的类型信息就变成不同的了，一个是 int，一个是 int64，即便 data 指针指向的内存块中存储值是相同的，最终 eif1 与 eif2 也是不相等的。
从输出结果中我们可以总结一下：对于空接口类型变量，只有 _type 和 data 所指数据内容一致的情况下，两个空接口类型变量之间才能划等号。另外，Go 在创建 eface 时一般会为 data 重新分配新内存空间，将动态类型变量的值复制到这块内存空间，并将 data 指针指向这块内存空间。因此我们多数情况下看到的 data 指针值都是不同的。
*/
func printEmptyInterface() {
	var eif1 interface{} // 空接口类型
	var eif2 interface{} // 空接口类型
	var n, m int = 17, 18

	eif1 = n
	eif2 = m
	println("eif1:", eif1)
	println("eif2:", eif2)
	println("eif1 = eif2:", eif1 == eif2) // false

	eif2 = 17
	println("eif1:", eif1)
	println("eif2:", eif2)
	println("eif1 = eif2:", eif1 == eif2) // true

	eif2 = int64(17)
	println("eif1:", eif1)
	println("eif2:", eif2)
	println("eif1 = eif2:", eif1 == eif2) // false
}

// 第三种：非空接口类型变量
/*
err1: (0x10c0240,0x0)
err1 = nil: false
err1: (0x10c02a0,0x10bfa58)
err2: (0x10c02a0,0x10bfa60)
err1 = err2: false
err1: (0x10c02a0,0x10bfa58)
err2: (0x10c01c0,0xc000010230)
err1 = err2: false

看到上面示例中每一轮通过 println 输出的 err1 和 err2 的 tab 和 data 值，要么 data 值不同，要么 tab 与 data 值都不同。
和空接口类型变量一样，只有 tab 和 data 指的数据内容一致的情况下，两个非空接口类型变量之间才能划等号。这里我们要注意 err1 下面的赋值情况：
err1 = (*T)(nil)
针对这种赋值，println 输出的 err1 是（0x10c0240, 0x0），也就是非空接口类型变量的类型信息并不为空，数据指针为空，因此它与 nil（0x0,0x0）之间不能划等号。

开头的问题中，从 returnsError 返回的 error 接口类型变量 err 的数据指针虽然为空，但它的类型信息（iface.tab）并不为空，而是 *MyError 对应的类型信息，这样 err 与 nil（0x0,0x0）相比自然不相等，这就是我们开头那个问题的答案解析，现在你明白了吗？
*/

type T int

func (t T) Error() string {
	return "bad error"
}
func printNonEmptyInterface() {
	var err1 error // 非空接口类型
	var err2 error // 非空接口类型
	err1 = (*T)(nil)
	println("err1:", err1)
	println("err1 = nil:", err1 == nil)
	err1 = T(5)
	err2 = T(6)
	println("err1:", err1)
	println("err2:", err2)
	println("err1 = err2:", err1 == err2)
	err2 = fmt.Errorf("%d\n", 5)
	println("err1:", err1)
	println("err2:", err2)
	println("err1 = err2:", err1 == err2)
}

// 第四种：空接口类型变量与非空接口类型变量的等值比较
/*
eif: (0x1094220,0x10bfc90)
err: (0x10c04e0,0x10bfc90)
eif = err: true
eif: (0x1094220,0x10bfc90)
err: (0x10c04e0,0x10bfc98)
eif = err: false

空接口类型变量和非空接口类型变量内部表示的结构有所不同（第一个字段：_type vs. tab)，两者似乎一定不能相等。但 Go 在进行等值比较时，类型比较使用的是 eface 的 _type 和 iface 的 tab._type，因此就像我们在这个例子中看到的那样，当 eif 和 err 都被赋值为T(5)时，两者之间是划等号的。
*/
func printEmptyInterfaceAndNonEmptyInterface() {
	var eif interface{} = T(5)
	var err error = T(5)
	println("eif:", eif)
	println("err:", err)
	println("eif = err:", eif == err)
	err = T(6)
	println("eif:", eif)
	println("err:", err)
	println("eif = err:", eif == err)
}

func main() {
	printNilInterface()
	println("=========")
	printEmptyInterface()
	println("=========")
	printNonEmptyInterface()
	println("=========")
	printEmptyInterfaceAndNonEmptyInterface()
}
