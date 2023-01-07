package gin

import "github.com/gin-gonic/gin"

// Wrapper 封装具体的函数实现
func Wrapper(exec any) gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}
