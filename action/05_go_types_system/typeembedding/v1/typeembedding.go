package main

import "fmt"

// 嵌入类型，用于修改已有类型以符合新类型

// user 在程序里面定义一个用户类型
type user struct {
	name  string
	email string
}

// notify 实现了一个可以通过 user 类型值的指针，调用方法
func (u *user) notify() {
	fmt.Printf("Sending user email to %s<%s>\n",
		u.name,
		u.email)
}

// admin 代表一个拥有权限的管理员用户
type admin struct {
	user  // 嵌入类型
	level string
}

// main 程序的主入口
func main() {
	// 创建一个 admin 用户
	ad := admin{
		user: user{
			name:  "john smith",
			email: "john@gmail.com",
		},
		level: "super",
	}

	// 我们可以直接访问内部类型的方法
	ad.user.notify()

	// 内部类型的方法也可以被提升到外部类型
	ad.notify()
}
