package curl

import (
	"fmt"
	"os"
)

// 实现一个简单版本的curl

// init 在main函数之前执行
func init() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./example <url>")
		os.Exit(-1)
	}
}
