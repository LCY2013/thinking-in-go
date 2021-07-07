package main

import (
	_ "fufeng.org/sample/matchers"
	"fufeng.org/sample/search"
	"log"
	"os"
)

//init 在main函数之前执行
func init() {
	// 将日志输出到标准输出
	log.SetOutput(os.Stdout)
}

//main 程序主入口
func main() {
	// 使用特定项作为搜索
	search.Run("president")
}
