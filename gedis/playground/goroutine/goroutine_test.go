package goroutine

import (
	"fmt"
	"testing"
)

func do1() {
	// 函数跳转时会调用如下:
	// runtime.morestack()
	do2()
}

func do2() {
	// runtime.morestack()
	do3()
}

func do3() {
	fmt.Println("do3...")
}

// TestGoroutine 协程信息
// 原生类型实现 src/runtime/runtime2.go:407 g
func TestGoroutine(t *testing.T) {

}
