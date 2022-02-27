package main

type E1 interface {
	M1()
	M2()
	M3()
}

type E2 interface {
	M1()
	M2()
	M4()
}

type T struct {
	E1
	E2
}

// main 嵌入了其他类型的结构体类型本身是一个代理，在调用其实例所代理的方法时，Go 会首先查看结构体自身是否实现了该方法。
// 如果实现了，Go 就会优先使用结构体自己实现的方法。
// 如果没有实现，那么 Go 就会查找结构体中的嵌入字段的方法集合中，是否包含了这个方法。
// 如果多个嵌入字段的方法集合中都包含这个方法，那么我们就说方法集合存在交集。
// 这个时候，Go 编译器就会因无法确定究竟使用哪个方法而报错。
func main() {
	t := T{}
	//t.M1()
	//t.M2()
	t.M3()
}