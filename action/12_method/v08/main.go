package main

/*
类型嵌入（Type Embedding）

Go 语言支持两种类型嵌入，分别是接口类型的类型嵌入和结构体类型的类型嵌入。

*/

// E 接口类型的类型嵌入
// 接口类型声明了由一个方法集合代表的接口，比如下面接口类型 E：
type E interface {
	M1()
	M2()
}

type I interface {
	M1()
	M2()
	M3()
}

// I1 在一个接口类型（I1）定义中，嵌入另外一个接口类型（E）的方式，就是我们说的接口类型的类型嵌入。
// 这个带有类型嵌入的接口类型 I 的定义与上面那个包含 M1、M2 和 M3 的接口类型 I 的定义，是等价的。
// 因此，我们可以得到一个结论，这种接口类型嵌入的语义就是新接口类型（如接口类型 I）将嵌入的接口类型（如接口类型 E）的方法集合，并入到自己的方法集合中。
type I1 interface {
	E
	M3()
}

// 通过嵌入其他接口类型来创建新接口类型的方式，在 Go 1.14 版本之前是有约束的：
// 如果新接口类型嵌入了多个接口类型，这些嵌入的接口类型的方法集合不能有交集，
// 同时嵌入的接口类型的方法集合中的方法名字，也不能与新接口中的其他方法同名。

type Interface1 interface {
	M1()
}
type Interface2 interface {
	M1()
	M2()
}
type Interface3 interface {
	Interface1
	Interface2 // Error: duplicate method M1
}
type Interface4 interface {
	Interface2
	M2() // Error: duplicate method M2
}

func main() {
}
