package main

import (
	"log"
)

// 使用最基本的 log 包
// 日志前缀，日志时间戳，日志由那个源文件记录，日志记录源文件所在行，日志消息

// init 设置日志相关配置
func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

// main 程序的入口
func main() {
	// Println 写入到标准日志记录器
	log.Println("message")

	// Fatalln 在调用 Println() 之后会接着调用 os.Exit()
	log.Fatalln("fatal message")

	// Panicln 在调用 Println() 之后接着调用 Panic
	log.Panicln("panic message")
}
