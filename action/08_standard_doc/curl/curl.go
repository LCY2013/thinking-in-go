package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// 这里的 resp 是一个响应， resp.Body 是 io.Reader
	resp, err := http.Get(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	// 创建文件来保存响应内容
	file, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// 使用 MultiWriter，这里同时向标准输出和文件写入内容
	dest := io.MultiWriter(os.Stdout, file)

	// 读出响应的内容，并写入到两个目的地
	io.Copy(dest, resp.Body)
	if err := resp.Body.Close(); err != nil {
		log.Println(err)
	}
}
