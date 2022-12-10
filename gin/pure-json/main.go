package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 通常，JSON 使用 unicode 替换特殊 HTML 字符，例如 < 变为 \ u003c。
//如果要按字面对这些字符进行编码，则可以使用 PureJSON。Go 1.6 及更低版本无法使用此功能。

func main() {
	router := gin.Default()

	// 提供unicode示例
	router.GET("/json", func(context *gin.Context) {
		// {"html":"\u003cb\u003eHello, world!\u003c/b\u003e"}
		context.JSON(http.StatusOK, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})

	// 提供字面常量
	router.GET("/pure_json", func(context *gin.Context) {
		// {"html":"<b>Hello, world!</b>"}
		context.PureJSON(http.StatusOK, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})

	// 监听 0.0.0.0:8080 并启动服务
	router.Run(":8080")
}

// curl -v http://127.0.0.1:8080/json
// curl -v http://127.0.0.1:8080/pure_json
