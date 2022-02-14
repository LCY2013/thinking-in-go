package main

import (
	"bytes"
	"fufeng.org/struct/basic"
	"unsafe"
)

func main() {
	// 输出空结构体信息
	printEmptyStruct()
	// empty struct to chan
	emptyStructToChan()

	// 输出内嵌struct
	printEmbeddedStruct()

	// init struct
	printInitStruct()

	// fill struct
	printFillStruct()
}

func printFillStruct() {
	var t basic.TFill
	println(unsafe.Sizeof(t))
	var s basic.SFill
	println(unsafe.Sizeof(s))
}

func printInitStruct() {
	var b bytes.Buffer
	b.Write([]byte("Hello, go"))
	println(b.String())
}

func printEmbeddedStruct() {
	var book basic.Book
	println(book.Author.Phone)

	var inform basic.Inform
	println(inform.Phone)
	println(inform.Person.Phone)
}

func emptyStructToChan() {
	// 内存占用最小的goroutine通信
	var ch = make(chan basic.Empty, 1)
	ch <- basic.Empty{}
}

// printEmptyStruct 输出空结构体信息
func printEmptyStruct() {
	// 空结构体内存占用长度为 0
	var empty basic.Empty
	println(unsafe.Sizeof(empty))
}
