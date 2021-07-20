package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	// 创建一个 Buffer 值，并将一个字符串写入到 Buffer
	// 使用实现 io.Writer 的 write 方法
	var b bytes.Buffer
	b.Write([]byte("hello"))

	// 使用 Fprintf 来将一个字符串拼接到 Buffer 中
	// 将 bytes.Buffer 的地址作为 io.Writer 类型值传入
	fmt.Fprintf(&b, " World!")

	// 将 Buffer 中的内容输出到标准输出
	// 将 os.File 值的地址作为 io.Writer 类型值输出
	b.WriteTo(os.Stdout)
}
