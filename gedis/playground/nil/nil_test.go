package main

import (
	"fmt"
	"testing"
)

func TestNil(t *testing.T) {
	// eface
	var a interface{}
	fmt.Println(a == nil)

	var b *int
	fmt.Println(b == nil)

	// eface 有了type就不是nil接口了
	// type eface struct {
	//	_type *_type
	//	data  unsafe.Pointer
	//}
	a = b
	fmt.Println(a == nil)
}
