package main

import (
	"fmt"
	"reflect"
)

// T 基于类型别名（type alias）定义的新类型有没有“继承”原类型的方法集合呢？
type T struct{}

func (T) M1()  {}
func (*T) M2() {}

type T1 = T

/*
结果：
main.T's method set:
- M1

main.T's method set:
- M1

*main.T's method set:
- M1
- M2

*main.T's method set:
- M1
- M2

这个输出结果，可以看到dumpMethodSet 函数甚至都无法识别出“类型别名”，无论类型别名还是原类型，输出的都是原类型的方法集合。

结论：无论原类型是接口类型还是非接口类型，类型别名都与原类型拥有完全相同的方法集合。
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
