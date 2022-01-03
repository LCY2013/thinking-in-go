package main

/*
map 类型对 value 的类型没有限制，但是对 key 的类型却有严格要求，因为 map 类型要保证 key 的唯一性。
Go 语言中要求，key 的类型必须支持“==”和“!=”两种比较操作符。


*/
func main() {
	typeCompare()
	declaredMap()
	mapCouture()
}

/*
和切片相比，map 类型的内部实现要更加复杂。Go 运行时使用一张哈希表来实现抽象的 map 类型。
运行时实现了 map 类型操作的所有功能，包括查找、插入、删除等。在编译阶段，Go 编译器会将 Go 语法层面的 map 操作，重写成运行时对应的函数调用
*/
func mapCouture() {
	// 创建map类型变量实例
	//m := make(map[keyType]valType, capacityhint) → m := runtime.makemap(maptype, capacityhint, m)
	// 插入新键值对或给键重新赋值
	//m["key"] = "value" → v := runtime.mapassign(maptype, m, "key") v是用于后续存储value的空间的地址
	// 获取某键的值
	//v := m["key"]      → v := runtime.mapaccess1(maptype, m, "key")
	//v, ok := m["key"]  → v, ok := runtime.mapaccess2(maptype, m, "key")
	// 删除某键
	//delete(m, "key")   → runtime.mapdelete(maptype, m, “key”)
}

/*
map 变量的声明和初始化

和切片类型变量一样，如果没有显式地赋予 map 变量初值，map 类型变量的默认值为 nil。

不过切片变量和 map 变量在这里也有些不同。初值为零值 nil 的切片类型变量，可以借助内置的 append 的函数进行操作，这种在 Go 语言中被称为“零值可用”。
定义“零值可用”的类型，可以提升我们开发者的使用体验，我们不用再担心变量的初始状态是否有效。
但 map 类型，因为它内部实现的复杂性，无法“零值可用”。
所以，如果我们对处于零值状态的 map 变量直接进行操作，就会导致运行时异常（panic），从而导致程序进程异常退出：
var m map[string]int // m = nil
m["key"] = 1         // 发生运行时异常：panic: assignment to entry in nil map

对 map 类型变量进行显式初始化后才能使用。那我们怎样对 map 类型变量进行初始化呢？
和切片一样，为 map 类型变量显式赋值有两种方式：一种是使用复合字面值；另外一种是使用 make 这个预声明的内置函数。
m := map[int]string{}
显式初始化了 map 类型变量 m。不过，你要注意，虽然此时 map 类型变量 m 中没有任何键值对，但变量 m 也不等同于初值为 nil 的 map 变量。这个时候，我们对 m 进行键值对的插入操作，不会引发运行时异常。

通过稍微复杂一些的复合字面值，对 map 类型变量进行初始化：
*/
func declaredMap() {
	// 字面量声明
	//var m map[string]int // 一个map[string]int类型的变量
	//m["key"] = 1
	m1 := map[int][]string{
		1: []string{"val1_1", "val1_2"},
		3: []string{"val3_1", "val3_2", "val3_3"},
		7: []string{"val7_1"},
	}
	println(m1)
	type Position struct {
		x float64
		y float64
	}
	m2 := map[Position]string{
		Position{29.935523, 52.568915}:  "school",
		Position{25.352594, 113.304361}: "shopping-mall",
		Position{73.224455, 111.804306}: "hospital",
	}
	println(m2)

	// 上面虽然完成了对两个 map 类型变量 m1 和 m2 的显式初始化，但不知道你有没有发现一个问题，作为初值的字面值似乎有些“臃肿”。
	// 你看，作为初值的字面值采用了复合类型的元素类型，而且在编写字面值时还带上了各自的元素类型，比如作为 map[int] []string 值类型的[]string，以及作为 map[Position]string 的 key 类型的 Position。

	// 针对这种情况，Go 提供了“语法糖”。这种情况下，Go 允许省略字面值中的元素类型。因为 map 类型表示中包含了 key 和 value 的元素类型，Go 编译器已经有足够的信息，来推导出字面值中各个值的类型了。
	//我们以 m2 为例，这里的显式初始化代码和上面变量 m2 的初始化代码是等价的：
	m3 := map[Position]string{
		{29.935523, 52.568915}:  "school",
		{25.352594, 113.304361}: "shopping-mall",
		{73.224455, 111.804306}: "hospital",
	}
	println(m3)

	// 使用 make 为 map 类型变量进行显式初始化。
	//和切片通过 make 进行初始化一样，通过 make 的初始化方式，我们可以为 map 类型变量指定键值对的初始容量，但无法进行具体的键值对赋值，就像下面代码这样：
	m4 := make(map[int]string)    // 未指定初始容量
	m5 := make(map[int]string, 8) // 指定初始容量为8
	println(m4)
	println(m5)
	// 不过，map 类型的容量不会受限于它的初始容量值，当其中的键值对数量超过初始容量后，Go 运行时会自动增加 map 类型的容量，保证后续键值对的正常插入。

}

/*
在 Go 语言中，函数类型、map 类型自身，以及切片只支持与 nil 的比较，而不支持同类型两个变量的比较。
如果像下面代码这样，进行这些类型的比较，Go 编译器将会报错：

函数类型、map 类型自身，以及切片类型是不能作为 map 的 key 类型的。
*/
func typeCompare() {
	//s1 := make([]int, 1)
	//s2 := make([]int, 2)
	//f1 := func() {}
	//f2 := func() {}
	//m1 := make(map[int]string)
	//m2 := make(map[int]string)
	//println(s1 == s2) // 错误：invalid operation: s1 == s2 (slice can only be compared to nil)
	//println(f1 == f2) // 错误：invalid operation: f1 == f2 (func can only be compared to nil)
	//println(m1 == m2) // 错误：invalid operation: m1 == m2 (map can only be compared to nil)
}
