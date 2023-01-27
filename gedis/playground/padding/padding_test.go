package main

import (
	"testing"
	"unsafe"
)

type S1 struct {
	num1 int32
	num2 int32
}

type S2 struct {
	num1 int16
	num2 int32
}

// TestPadding 内存对齐
func TestPadding(t *testing.T) {
	t.Log(unsafe.Sizeof(S1{}))
	t.Log(unsafe.Sizeof(S2{}))
}

// TestPaddingAlign 内存对齐系数
func TestPaddingAlign(t *testing.T) {
	t.Logf("bool size: %d, align: %d", unsafe.Sizeof(bool(true)), unsafe.Alignof(bool(true)))
	t.Logf("int size: %d, align: %d", unsafe.Sizeof(int(1)), unsafe.Alignof(int(1)))
	t.Logf("int8 size: %d, align: %d", unsafe.Sizeof(int8(1)), unsafe.Alignof(int8(1)))
	t.Logf("int16 size: %d, align: %d", unsafe.Sizeof(int16(1)), unsafe.Alignof(int16(1)))
	t.Logf("int32 size: %d, align: %d", unsafe.Sizeof(int32(1)), unsafe.Alignof(int32(1)))
	t.Logf("int64 size: %d, align: %d", unsafe.Sizeof(int64(1)), unsafe.Alignof(int64(1)))
	t.Logf("string size: %d, align: %d", unsafe.Sizeof(""), unsafe.Alignof(""))
}
