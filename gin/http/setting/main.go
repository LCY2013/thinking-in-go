package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// main 自定义 HTTP 配置
func mainDefault() {
	router := gin.Default()
	http.ListenAndServe(":8080", router)
}

func main() {
	router := gin.Default()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
