package main

import (
	"fmt"
	"reflect"
)

// 结构体类型中嵌入结构体类型
/*
在结构体类型中嵌入结构体类型，为 Gopher 们提供了一种“实现继承”的手段，外部的结构体类型 T 可以“继承”嵌入的结构体类型的所有方法的实现。
并且，无论是 T 类型的变量实例还是 *T 类型变量实例，都可以调用所有“继承”的方法。
但这种情况下，带有嵌入类型的新类型究竟“继承”了哪些方法，我们还要通过下面这个具体的示例来看一下。
*/

type T1 struct{}

func (T1) T1M1()   { println("T1's M1") }
func (*T1) PT1M2() { println("PT1's M2") }

type T2 struct{}

func (T2) T2M1()   { println("T2's M1") }
func (*T2) PT2M2() { println("PT2's M2") }

/*
T1 的方法集合包含：T1M1；
*T1 的方法集合包含：T1M1、PT1M2；
T2 的方法集合包含：T2M1；
*T2 的方法集合包含：T2M1、PT2M2。
*/

// T 结构体包含
type T struct {
	T1
	*T2
}

/*
通过输出结果，我们看到了 T 和 *T 类型的方法集合果然有差别的：
类型 T 的方法集合 = T1 的方法集合 + *T2 的方法集合
类型 *T 的方法集合 = *T1 的方法集合 + *T2 的方法集合

注意 *T 类型的方法集合，它包含的可不是 T1 类型的方法集合，而是 *T1 类型的方法集合。
这和结构体指针类型的方法集合包含结构体类型方法集合，是一个道理。
*/
// main 测试
func main() {
	t := T{
		T1: T1{},
		T2: &T2{},
	}
	dumpMethodSet(t)
	dumpMethodSet(&t)
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
