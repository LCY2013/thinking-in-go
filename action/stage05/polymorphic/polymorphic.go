package main

import "fmt"

// 多态
// notifier 是一个定义了通知类行为的接口
type notifier interface {
	notify()
}

// user 在程序里定义了一个用户类型
type user struct {
	name  string
	email string
}

// notify 使用指针接收者实现了 notifier 接口
func (u *user) notify() {
	fmt.Printf("Sending user email to %s<%s>\n",
		u.name,
		u.email)
}

// admin 定义程序里面的管理员信息
type admin struct {
	name  string
	email string
}

// notify 使用指针接收者实现 notifier 接口
func (a *admin) notify() {
	fmt.Printf("Sending admin email to %s<%s>\n",
		a.name,
		a.email)
}

// sendNotification 接受一个实现了 notifier 接口的值，并发送通知
func sendNotification(n notifier) {
	n.notify()
}

// main 程序的主入口
func main() {
	// 创建一个 user 值，并传递给 sendNotification
	bill := user{"Bill", "bill@email.com"}
	sendNotification(&bill)

	// 创建一个 admin 类型的值，并传递给 sendNotification
	lisa := admin{"Lisa", "lisa@email.com"}
	sendNotification(&lisa)
}
