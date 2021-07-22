package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// main 使用 HTTP 方法
func main() {
	// 禁用控制台颜色
	gin.DisableConsoleColor()

	// 使用默认中间件（logger 和 recovery 中间件）创建 gin 路由
	router := gin.Default()

	// curl -X GET http://127.0.0.1:8080/some
	router.GET("/some", rest)
	// curl -X POST http://127.0.0.1:8080/some
	router.POST("/some", rest)
	// curl -X PUT http://127.0.0.1:8080/some
	router.PUT("/some", rest)
	// curl -X DELETE http://127.0.0.1:8080/some
	router.DELETE("/some", rest)
	// curl -X PATCH http://127.0.0.1:8080/some
	router.PATCH("/some", rest)
	// curl -I http://127.0.0.1:8080/some
	// curl --head http://127.0.0.1:8080/some
	router.HEAD("/some", rest)
	// curl -X OPTIONS http://127.0.0.1:8080/some
	router.OPTIONS("/some", rest)

	// 默认在 8080 端口启动服务，除非定义了一个 PORT 的环境变量。
	router.Run()
	// router.Run(":3000") hardcode 端口号
}

func rest(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"httpMethod": context.Request.Method,
	})
}
