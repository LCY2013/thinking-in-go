package main

import (
	"fmt"
	"reflect"
)

/*
defined 类型与 alias 类型的方法集合

Go 语言中，凡通过类型声明语法声明的类型都被称为 defined 类型，下面是一些 defined 类型的声明的例子：
type I interface {
	M1()
	M2()
}
type T int
type NT T // 基于已存在的类型T创建新的defined类型NT
type NI I // 基于已存在的接口类型I创建新defined接口类型NI
*/

// 新定义的 defined 类型与原 defined 类型是不同的类型，那么它们的方法集合上又会有什么关系呢？新类型是否“继承”原 defined 类型的方法集合呢？
// 对于那些基于接口类型创建的 defined 的接口类型，它们的方法集合与原接口类型的方法集合是一致的。
// 但对于基于非接口类型的 defined 类型创建的非接口类型，通过下面例子来看一下：

type T struct{}

func (T) M1()  {}
func (*T) M2() {}

type T1 T

/*
结果：
main.T's method set:
- M1
main.T1's method set is empty!
*main.T's method set:
- M1
- M2
*main.T1's method set is empty!

从输出结果上看，新类型 T1 并没有“继承”原 defined 类型 T 的任何一个方法。
从逻辑上来说，这也符合 T1 与 T 是两个不同类型的语义。

基于自定义非接口类型的 defined 类型的方法集合为空的事实，也决定了即便原类型实现了某些接口，基于其创建的 defined 类型也没有“继承”这一隐式关联。
也就是说，新 defined 类型要想实现那些接口，仍然需要重新实现接口的所有方法。
*/

func main() {
	var t T
	var pt *T
	var t1 T1
	var pt1 *T1
	dumpMethodSet(t)
	dumpMethodSet(t1)
	dumpMethodSet(pt)
	dumpMethodSet(pt1)
}

// dumpMethodSet 列举某个类型的所有方法集合
func dumpMethodSet(i interface{}) {
	dynTyp := reflect.TypeOf(i)
	if dynTyp == nil {
		fmt.Printf("there is no dynamic type\n")
		return
	}
	n := dynTyp.NumMethod()
	if n == 0 {
		fmt.Printf("%s's method set is empty!\n", dynTyp)
		return
	}
	fmt.Printf("%s's method set:\n", dynTyp)
	for j := 0; j < n; j++ {
		fmt.Println("-", dynTyp.Method(j).Name)
	}
	fmt.Printf("\n")
}
