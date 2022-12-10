package main

import (
	"fmt"
	"time"
)

type field struct {
	name string
}

// func (p *field) print() {
func (p field) print() {
	fmt.Println(p.name)
}

func main() {
	data1 := []*field{{"one"}, {"two"}, {"three"}}
	for _, v := range data1 {
		go v.print()
	}
	data2 := []field{{"four"}, {"five"}, {"six"}}
	for _, v := range data2 {
		go v.print()
	}
	time.Sleep(3 * time.Second)
	/*
		运行结果：
			three
			one
			two
			six
			six
			six
		这位读者的问题显然是：为什么对 data2 迭代输出的结果是三个“six”，而不是 four、five、six？

		Go 方法的本质，也就是一个以方法的 receiver 参数作为第一个参数的普通函数，对这个程序做个等价变换。这里利用 Method Expression 方式，等价变换后的源码如下：

		data1 := []*field{{"one"}, {"two"}, {"three"}}
		for _, v := range data1 {
			go (*field).print(v)
		}
		data2 := []field{{"four"}, {"five"}, {"six"}}
		for _, v := range data2 {
			go (*field).print(&v)
		}
		time.Sleep(3 * time.Second)
	*/

}
