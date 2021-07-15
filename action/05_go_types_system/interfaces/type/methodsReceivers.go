package main

import "fmt"

// duration 定义一个封装的int类型
type duration int

// 使用更可读的方式格式化duration的值
func (d *duration) pretty() string {
	return fmt.Sprintf("Duration: %d", *d)
}

// main 程序的主入口
func main() {
	// # command-line-arguments
	//./methodsReceivers.go:15:14: cannot call pointer method on duration(21)
	//./methodsReceivers.go:15:14: cannot take the address of duration(21)
	// duration(21).pretty()

	// 正确使用
	d := duration(21)
	fmt.Printf("duration(21) type is : %T\n", duration(21))
	fmt.Printf("duration type is : %T\n", d)
	d.pretty()
}
