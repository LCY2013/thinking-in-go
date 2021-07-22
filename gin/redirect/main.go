package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// main 重定向 HTTP 重定向很容易。 内部、外部重定向均支持。
func main() {
	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.Request.URL.Path = "/test2"
		router.HandleContext(c)
	})

	router.GET("/test2", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	router.GET("/baidu", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
	})

	router.POST("/notfound", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/foo")
	})

	router.Run(":8080")
}
