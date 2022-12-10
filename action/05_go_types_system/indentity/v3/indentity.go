package main

import (
	"fmt"
)

// 公开或者未公开的标识符

// main 程序的主入口
func main() {
	// 创建 entities 中的用户类型
	u := v1.User{
		Name: "Bill",
		//email: "bill@email.com",
	}

	//# fufeng.org/sample05/indentity/v3
	//./indentity.go:13:3: cannot refer to unexported field 'email' in struct literal of type entities.User

	fmt.Printf("User: %v\n", u)
}
