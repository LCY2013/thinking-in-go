package main

import (
	"fmt"
)

func main() {
	var a = make([]int, 0, 10)
	//appendInt(a)
	appendIntPointer(&a)
	fmt.Println(a)
	/*
		uintptr
		unsafe.Pointer()
	*/
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
