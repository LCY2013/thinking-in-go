package main

import (
	"fmt"
	v2 "fufeng.org/sample05/indentity/counters/v2"
)

// 公开或者未公开的标识符

// main 程序的入口
func main() {
	// 通过公开曝露的函数构建一个未公开访问的 alterCounter
	counter := v2.New(10)

	fmt.Printf("AltreCounter type : %T\n", counter)
	fmt.Printf("AltreCounter type : %d\n", counter)
}
