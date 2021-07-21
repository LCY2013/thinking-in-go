package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

/*
POST /post?id=1234&page=1 HTTP/1.1
Content-Type: application/x-www-form-urlencoded

name=manu&message=this_is_great
*/
func main() {
	router := gin.Default()

	router.POST("/post", func(c *gin.Context) {

		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")

		fmt.Printf("id: %s; page: %s; name: %s; message: %s\n", id, page, name, message)
	})
	router.Run(":8080")
}

// curl -v --form name=fufeng --form message=this_is_great http://127.0.0.1:8080/post?id=1&page=1
