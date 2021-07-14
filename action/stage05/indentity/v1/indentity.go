package main

import (
	_ "fufeng.org/sample05/indentity/counters/v1"
)

// 公开或者未公开的标识符

// main 函数程序的入口
func main() {
	// 创建一个未公开的类型的变量
	// counter := counters.alterCounter(10)
	// # command-line-arguments
	// ./indentity.go:13:13: cannot refer to unexported name counters.alterCounter
	// fmt.Printf("Counter : %d\n", counter)
}
