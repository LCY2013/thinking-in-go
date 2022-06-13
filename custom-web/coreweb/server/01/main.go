package main

import (
	"github.com/lcy2013/custom-web/coreweb/framework"
	"log"
	"net/http"
)

func main() {
	serve := &http.Server{
		// 自定义的请求核心处理函数
		Handler: framework.NewCore(),
		// 请求监听地址
		Addr: "8080",
	}
	err := serve.ListenAndServe()
	if err != nil {
		log.Panicln(err)
		return
	}
}
