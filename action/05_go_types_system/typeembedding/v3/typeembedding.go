package main

import "fmt"

// 嵌入类型，用于修改已有类型以符合新类型

// user 在程序里面定义一个用户类型
type user struct {
	name  string
	email string
}

// notifier 是一个定义了通知类行为的接口
type notifier interface {
	notify()
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

// notify 实现了一个可以通过 admin 类型的指针，调用方法
func (a *admin) notify() {
	fmt.Printf("Sending admin email to %s<%s>\n", a.name, a.email)
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

	// 给admin 用户发送一个通知，其内部接口实现提升到外部类型
	sendNotification(&ad)

	// 通过内部的类型访问
	ad.user.notify()

	// 通过外部的类型访问
	ad.notify()
}

// sendNotification 接受一个 notifier 接口的值并发送通知
func sendNotification(n notifier) {
	n.notify()
}
