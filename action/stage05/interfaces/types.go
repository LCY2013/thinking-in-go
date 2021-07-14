package main

import "fmt"

/*
// 规范里描述的方法集
Values			  Methods Receivers
  T						(t T)
 *T					(t T) & (t *T)

// 从接收者类型的角度看方法集
Methods Receivers 		Values
	(t T)				T & *T
   (t *T)				  *T

当使用指针类型作为接口的接收者时，只能只用该接收者类型的指针类型去使用，因为编译器不是每次都能自动获得一个值的地址
*/

// notifier 是一个定义了通知类行为的接口
type notifier interface {
	notify()
}

// user 定义一个用户类型
type user struct {
	name  string
	email string
}

// notify 是使用指针接收者
func (u *user) notify() {
	fmt.Printf("Send message to user email %s<%s>\n",
		u.name,
		u.email)
}

func main() {
	u := user{"Bill", "bill@email.com"}

	// # command-line-arguments
	//./types.go:26:18: cannot use u (type user) as type notifier in argument to sendNotification:
	//        user does not implement notifier (notify method has pointer receiver)
	// sendNotification(u) // 不能将u（user）作为sendNotification的参数类型notifier，notifier是一个指针的接受者声明

	sendNotification(&u)
}

// sendNotification 接受一个实现了notifier接口的值，并发送通知
func sendNotification(u notifier) {
	u.notify()
}
