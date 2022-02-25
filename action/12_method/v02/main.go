package main

import "fmt"

type T struct {
	a int
}

func (t T) Get() int {
	return t.a
}

func (t *T) Set(a int) int {
	t.a = a
	return t.a
}

// Get 类型T的方法Get的等价函数
func Get(t T) int {
	return t.a
}

// Set 类型*T的方法Set的等价函数
func Set(t *T, a int) int {
	t.a = a
	return t.a
}

func main() {
	var t T
	t.Get()
	t.Set(1)

	// == 等价于

	/*
		这种直接以类型名 T 调用方法的表达方式，被称为 Method Expression。
		通过 Method Expression 这种形式，类型 T 只能调用 T 的方法集合（Method Set）中的方法，同理类型 *T 也只能调用 *T 的方法集合中的方法。
	*/
	var t1 T
	T.Get(t1)
	(*T).Set(&t1, 1)

	/*
	 Method Expression 对方法进行调用的方式，与之前所做的方法到函数的等价转换是如出一辙的。
	 所以，Go 语言中的方法的本质就是，一个以方法的 receiver 参数作为第一个参数的普通函数。

	 Method Expression 就是 Go 方法本质的最好体现，因为方法自身的类型就是一个普通函数的类型，甚至可以将它作为右值，赋值给一个函数类型的变量。

	 结论：方法本质上也是函数
	*/
	var t2 T
	f1 := (*T).Set                           // f1的类型，也是T类型Set方法的类型：func (t *T, int)int
	f2 := T.Get                              // f2的类型，也是T类型Get方法的类型：func(t T)int
	fmt.Printf("the type of f1 is %T\n", f1) // the type of f1 is func(*main.T, int) int
	fmt.Printf("the type of f2 is %T\n", f2) // the type of f2 is func(main.T) int
	f1(&t2, 3)
	fmt.Println(f2(t2)) // 3
}
