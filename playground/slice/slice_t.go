package main

import "fmt"

func main() {
	//var a = make([]int, 0, 10)
	//appendInt(a)
	//appendIntPointer(&a)
	//fmt.Println(a)
	/*
		uintptr
		unsafe.Pointer()
	*/
	a := []int{1, 2, 3}
	var b = a
	var c = b
	c[0] = 0
	fmt.Println(c)
	fmt.Printf("%p\n", &c)
	fmt.Println(b)
	fmt.Printf("%p\n", &b)
	fmt.Println(a)
	fmt.Printf("%p\n", &a)
}

func appendInt(a []int) {
	for i := 0; i < 100000; i++ {
		a = append(a, i)
	}
}

func appendIntPointer(a *[]int) {
	for i := 0; i < 100000; i++ {
		*a = append(*a, i)
	}
}
