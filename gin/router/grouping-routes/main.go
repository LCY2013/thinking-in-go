package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// main 路由组
func main() {
	router := gin.Default()

	// 简单的路由组: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}

	// 简单的路由组: v2
	v2 := router.Group("/v2")
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
		v2.POST("/read", readEndpoint)
	}

	router.Run(":8080")
}

func readEndpoint(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"path": fmt.Sprintf("read %s", context.Request.RequestURI),
	})
}

func submitEndpoint(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"path": fmt.Sprintf("submit %s", context.Request.RequestURI),
	})
}

func loginEndpoint(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"path": fmt.Sprintf("login %s", context.Request.RequestURI),
	})
}
