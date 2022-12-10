package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
POST /post?ids[a]=1234&ids[b]=hello HTTP/1.1
Content-Type: application/x-www-form-urlencoded

names[first]=thinkerou&names[second]=tianou
*/
// main 映射查询字符串或表单参数
func main() {
	router := gin.Default()

	router.POST("/post", func(c *gin.Context) {

		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")

		// ids: map[b:hello a:1234], names: map[second:tianou first:thinkerou]
		fmt.Printf("ids: %v; names: %v\n", ids, names)
	})
	router.Run(":8080")
}

// curl -X POST http://127.0.0.1:8080/post\?ids\[a\]\=123\&ids\[b\]\=hello --form names\[first\]=thinkerou --form names\[second\]=tianou
// 输出： ids: map[a:123 b:hello]; names: map[first:thinkerou second:tianou]
