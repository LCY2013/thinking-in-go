package main

import (
	"fmt"
	"reflect"
)

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

//==================
//结构体类型中嵌入接口类型
//==================

type I interface {
	M1()
	M2()
}
type T struct {
	I
}

func (T) M3() {}

func (*T) M4() {}

// main 结构体类型的方法集合，包含嵌入的接口类型的方法集合。
func main() {
	var t T
	var p *T
	dumpMethodSet(t)
	dumpMethodSet(p)
}
