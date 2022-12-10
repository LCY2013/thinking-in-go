package main

import (
	"bytes"
	"fmt"

	//_ "fufeng.org/sample05/curl"
	"io"
	"net/http"
	"os"

	"fufeng.org/sample05/structs"
)

// main 程序的入口
func main() {
	// _001()

	// _002()

	// _003()

	// _004()

	_005()
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

// _004 curl 应用程序的入口
func _004() {
	// 从 web 服务器得到响应
	resp, err := http.Get(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// 从 body 到 Stdout
	io.Copy(os.Stdout, resp.Body)
	if err := resp.Body.Close(); err != nil {
		fmt.Println(err)
	}
}

func _005() {
	var b bytes.Buffer

	// 将字符写入buffer
	b.Write([]byte("Hello "))

	// 使用Fprintf将字符拼接到Buffer
	_, _ = fmt.Fprintf(&b, "World!")

	// 将Buffer的内容写入到Stdout
	io.Copy(os.Stdout, &b)
}
