package main

import (
	"fmt"
	"unsafe"
)

func arrayToSlice() {
	a := [3]int{1, 2, 3}
	// 通过切片化将数组转化成切片
	b := a[:]
	b[1] += 10
	// [1 12 3]
	fmt.Printf("%v\n", a)
}

// beforeGo1_17SliceToArray go 1.17 前的slice 转 array
func beforeGo1_17SliceToArray() {
	a := []int{1, 2, 3}
	var p = (*[3]int)(unsafe.Pointer(&a[0]))
	p[1] += 10
	// [1 12 3]
	fmt.Printf("%v\n", a)
}

// go1_17SliceToArray go 1.17 的slice 转 array
func go1_17SliceToArray() {
	a := []int{1, 2, 3}
	p := (*[3]int)(a)
	p[1] += 10
	// [1 12 3]
	fmt.Printf("%v\n", a)
}

// sliceToArrayTips 转换后的数组长度不能大于原切片的长度
func sliceToArrayTips() {
	var b = []int{11, 12, 13}
	//var p1 = (*[4]int)(b)     // cannot convert slice with length 3 to pointer to array with length 4
	var p2 = (*[0]int)(b) // ok，*p = []
	var p3 = (*[1]int)(b) // ok，*p = [11]
	var p4 = (*[2]int)(b) // ok，*p = [11, 12]
	var p5 = (*[3]int)(b) // ok，*p = [11, 12, 13]
	//var p6 = (*[3]int)(b[:1]) // cannot convert slice with length 1 to pointer to array with length 3
	//fmt.Printf("%v\n", p1)
	fmt.Printf("%v\n", p2)
	fmt.Printf("%v\n", p3)
	fmt.Printf("%v\n", p4)
	fmt.Printf("%v\n", p5)
	//fmt.Printf("%v\n", p6)
}

func main() {
	arrayToSlice()
	beforeGo1_17SliceToArray()
	go1_17SliceToArray()
	sliceToArrayTips()
}
