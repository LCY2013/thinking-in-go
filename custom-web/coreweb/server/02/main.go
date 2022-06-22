package main

import (
	"github.com/lcy2013/custom-web/coreweb/server/02/framework"
	"log"
	"net/http"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	serve := &http.Server{
		// 自定义的请求核心处理函数
		Handler: core,
		// 请求监听地址
		Addr: ":8080",
	}
	err := serve.ListenAndServe()
	if err != nil {
		log.Panicln(err)
		return
	}
}
