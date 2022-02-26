package main

import "fmt"

// Go 规定：如果一个类型 T 的方法集合是某接口类型 I 的方法集合的等价集合或超集，我们就说类型 T 实现了接口类型 I，那么类型 T 的变量就可以作为合法的右值赋值给接口类型 I 的变量。

// 果一个变量的类型是空接口类型，由于空接口类型的方法集合为空，这就意味着任何类型都实现了空接口的方法集合，所以我们可以将任何类型的值作为右值，赋值给空接口类型的变量，比如下面例子：
func emptyInterface() {
	var i interface{} = 15 // ok
	i = "hello, golang"    // ok
	type T struct{}
	var t T
	i = t  // ok
	i = &t // ok
	fmt.Println(i)
}

// Go 语言还支持接口类型变量赋值的“逆操作”，也就是通过接口类型变量“还原”它的右值的类型与值信息，这个过程被称为“类型断言（Type Assertion）”。
// 类型断言通常使用下面的语法形式：
// v, ok := i.(T)
// 其中 i 是某一个接口类型变量，如果 T 是一个非接口类型且 T 是想要还原的类型，那么这句代码的含义就是断言存储在接口类型变量 i 中的值的类型为 T。
//
// 类型断言也支持下面这种语法形式：
// v := i.(T)
// 但在这种形式下，一旦接口变量 i 之前被赋予的值不是 T 类型的值，那么这个语句将抛出 panic。如果变量 i 被赋予的值是 T 类型的值，那么变量 v 的类型为 T，它的值就会是之前变量 i 的右值。由于可能出现 panic，所以我们并不推荐使用这种类型断言的语法形式。
func typeAssertion() {
	var a int64 = 13
	var i interface{} = a
	v1, ok := i.(int64)
	fmt.Printf("v1=%d, the type of v1 is %T, ok=%t\n", v1, v1, ok) // v1=13, the type of v1 is int64, ok=true
	v2, ok := i.(string)
	fmt.Printf("v2=%s, the type of v2 is %T, ok=%t\n", v2, v2, ok) // v2=, the type of v2 is string, ok=false
	v3 := i.(int64)
	fmt.Printf("v3=%d, the type of v3 is %T\n", v3, v3) // v3=13, the type of v3 is int64
	v4 := i.([]int)                                     // panic: interface conversion: interface {} is int64, not []int
	fmt.Printf("the type of v4 is %T\n", v4)
}

func main() {
	typeAssertion()
}
