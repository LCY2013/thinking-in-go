package main

import (
	"fmt"
	"io"
	"strings"
)

type T struct {
}
type P struct {
}

// S 结构体类型的类型嵌入
//type S struct {
//	A int
//	b string
//	c T
//	p *P
//	_ [10]int8
//	F func()
//}

// 结构体类型 S 中的每个字段（field）都有唯一的名字与对应的类型，即便是使用空标识符占位的字段，它的类型也是明确的，但这还不是 Go 结构体类型的“完全体”。
// Go 结构体类型定义还有另外一种形式，那就是带有嵌入字段（Embedded Field）的结构体定义。看下面这个例子：
type T1 int

type t2 struct {
	n int
	m int
}

type I interface {
	M1()
}

// S1 这种以某个类型名、类型的指针类型名或接口类型名，直接作为结构体字段的方式就叫做结构体的类型嵌入，这些字段也被叫做嵌入字段（Embedded Field）。
type S1 struct {
	T1  // 标识符 T1 表示字段名为 T1，它的类型为自定义类型 T1；
	*t2 // 标识符 t2 表示字段名为 t2，它的类型为自定义结构体类型 t2 的指针类型；
	I   // 标识符 I 表示字段名为 I，它的类型为接口类型 I。
	a   int
	b   string
}

type MyInt int

func (n *MyInt) Add(m int) {
	*n = *n + MyInt(m)
}

type t struct {
	a int
	b int
}
type S struct {
	*MyInt
	t
	io.Reader
	s string
	n int
}

func main() {
	m := MyInt(17)
	r := strings.NewReader("hello, go")
	s := S{
		MyInt: &m,
		t: t{
			a: 1,
			b: 2,
		},
		Reader: r,
		s:      "demo",
	}
	var sl = make([]byte, len("hello, go"))
	_, _ = s.Reader.Read(sl)
	fmt.Println(string(sl)) // hello, go
	s.MyInt.Add(5)
	fmt.Println(*(s.MyInt)) // 22
}
