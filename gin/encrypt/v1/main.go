package main

import (
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"log"
)

/*
支持 Let's Encrypt
一行代码支持 LetsEncrypt HTTPS servers 示例。
*/

func main() {
	r := gin.Default()

	// Ping handler
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	log.Fatal(autotls.Run(r, "example1.com", "example2.com"))
}
