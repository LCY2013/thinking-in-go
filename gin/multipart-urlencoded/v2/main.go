package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/form_post", func(context *gin.Context) {
		message := context.PostForm("message")
		name := context.DefaultPostForm("name", "fufeng")

		context.JSON(http.StatusOK, gin.H{
			"statu":   "posted",
			"name":    name,
			"message": message,
		})
	})

	router.Run(":8080")
}

// curl -v --form message=hello --form name=lcy http://127.0.0.1:8080/form_post
// curl -v --form message=hello http://127.0.0.1:8080/form_post
