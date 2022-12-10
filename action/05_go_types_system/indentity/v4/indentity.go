package main

import (
	"fmt"

	v2 "fufeng.org/sample05/indentity/entities/v2"
)

// 公开或者未公开的标识符

// main 应用程序入口
func main() {
	// 创建 entities 包中的 Admin 类型的值
	a := v2.Admin{
		Rights: 10,
	}

	// 设置未公开的内部类型的公开类型的值
	a.Name = "Bill"
	a.Email = "bill@email.com"

	fmt.Printf("User: %v\n", a)
}
