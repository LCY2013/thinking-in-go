package main

import "fufeng.org/sample05/structs"

func main() {
	// _001()

	// _002()

	_003()

}

func _001() {
	structs.UserFunc()
}

func _002() {
	// 创建一个 User 类型的结构体
	user := structs.CreateUser("fufeng", "fufeng@gmail.com")

	// 使用该结构体的函数,输出其内部信息
	(*user).Notify()

	// 改变 User 类型的结构体内部数据
	user.ChangeEmail("fufeng@qq.com")

	// 使用该结构体的函数，输出其内部信息
	user.Notify()
}

func _003() {
	// 创建一个 User 类型的结构体
	user := structs.NewUser("fufeng", "fufeng@gmail.com")

	// 使用该结构体的函数,输出其内部信息
	user.Notify()

	// 改变 User 类型的结构体内部数据
	(&user).ChangeEmail("fufeng@qq.com")

	user.ChangeEmail("fufeng@qq.com")

	// 使用该结构体的函数，输出其内部信息
	user.Notify()
}
