package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// main 使用 AsciiJSON 生成具有转义的非 ASCII 字符的 ASCII-only JSON
func main() {
	r := gin.Default()

	// 定义一个相对路由
	r.GET("/someJSON", func(context *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "</br>",
		}

		// 输出 : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		context.AsciiJSON(http.StatusOK, data)
	})

	// 默认监听在 0.0.0.0:8080 上面
	r.Run()
}
