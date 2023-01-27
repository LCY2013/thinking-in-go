package size

import (
	"reflect"
	"testing"
	"unsafe"
)

type Message struct {
}

type Input struct {
	msg Message
	bit int
}

func TestSizeOfNilStruct(t *testing.T) {
	a := Message{}
	b := int(0)
	c := Message{}
	// src/runtime/malloc.go:778
	// 所有空结构体都指向了zerobase, 但是里面有成员变量不是空结构体的就不是, 如果里面都是结构体那么就还是可以指向zerobase
	// var zerobase uintptr
	t.Log(unsafe.Sizeof(a))
	t.Logf("%p\n", &a)
	t.Logf("%p\n", &b)
	t.Logf("%p\n", &c)

	d := Input{}
	t.Logf("%p\n", &d.msg)
	t.Logf("%p\n", &d)
}

// TestSizeOfString
// 原生string src/runtime/string.go:238   stringStruct
// reflect src/reflect/value.go:2670	StringHeader
func TestSizeOfString(t *testing.T) {
	t.Log(unsafe.Sizeof("123"))
	t.Log(unsafe.Sizeof("123456"))
}

// TestTypeOfString
// 查询string类型在go里面的真实存储形式
func TestTypeOfString(t *testing.T) {
	str := "扶风111"

	rsh := (*reflect.StringHeader)(unsafe.Pointer(&str))

	t.Log(rsh.Len)
}

// TestRangeOfString src/runtime/utf8.go:60 utf8解码
func TestRangeOfString(t *testing.T) {
	str := "扶风111"

	for _, c := range str {
		t.Logf("%c", c)
	}

	for i := 0; i < len(str); i++ {
		t.Logf("%c", str[i])
	}
}
