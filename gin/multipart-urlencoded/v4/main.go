package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// main 查询字符串参数
func main() {
	router := gin.Default()

	// 使用现有的基础请求对象解析查询字符串参数。
	// 示例 URL： /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // c.Request.URL.Query().Get("lastname") 的一种快捷方式

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	router.Run(":8080")
}

// curl http://127.0.0.1:8080/welcome\?firstname\=Jane\&lastname\=Doe
