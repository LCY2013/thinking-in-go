package main

import "github.com/gin-gonic/gin"

type myForm struct {
	Colors []string `form:"colors[]"`
}

func formHandler(c *gin.Context) {
	var fakeForm myForm
	c.ShouldBind(&fakeForm)
	c.JSON(200, gin.H{"color": fakeForm.Colors})
}

// main 绑定 HTML 复选框
func main() {
	router := gin.Default()
	router.GET("/check", formHandler)

	router.Run(":8080")
}
