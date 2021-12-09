package main

import "fmt"

const Pi float64 = 3.14159265358979323846 // 单行常量声明
// 以const代码块形式声明常量
const (
	size    int64 = 4096
	i, j, s       = 13, 14, "bar" // 单行声明多个常量
)

type myInt int

const n myInt = 13

//const m int = n + 5 // 编译器报错：cannot use n + 5 (type myInt) as type int in const initializer
const m int = int(n) + 5 // 编译器报错：cannot use n + 5 (type myInt) as type int in const initializer

// 无类型常量
const unTypedConst = 13

// 对于int8溢出
const overflowInt8 = 133333333

// 通过常量实现枚举
const (
	Apple, Banana     = 11, 22
	Strawberry, Grape // 这里会使用上一行的表达式计算
	Pear, Watermelon
)

// == 等价于下面
/*
const (
	Apple, Banana = 11, 22
	Strawberry,Grape = 11, 22
	Pear, Watermelon = 11, 22
)
*/

// iota 从0开始，一行一增加
const (
	a = iota + 1 // 1 iota = 0
	b            // 2 iota = 1
	c            // 3 iota = 2
)

func main() {
	constAdd()
	unTypeConstAdd()
	overflow()
}

func overflow() {
	//var k int8 = 1
	//j := k + m //Invalid operation: k + m (mismatched types int8 and int)
}

func unTypeConstAdd() {
	var a myInt = 5
	fmt.Println(a + unTypedConst)
}

func constAdd() {
	var a int = 5
	//fmt.Println(a + n) // 编译器报错：invalid operation: a + n (mismatched types int and myInt)
	fmt.Println(a + int(n)) // 编译器报错：invalid operation: a + n (mismatched types int and myInt)
}
