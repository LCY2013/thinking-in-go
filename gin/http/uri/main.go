package main

import (
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Person struct {
	ID   string `uri:"id" form:"id" binding:"required"`
	Name string `uri:"name" form:"name" binding:"required"`
}

// main 绑定 Uri
func main() {
	route := gin.Default()
	//route.GET("/:name/:id", func(c *gin.Context) {
	//	var person Person
	//	if err := c.ShouldBindUri(&person); err != nil {
	//		c.JSON(400, gin.H{"msg": err})
	//		return
	//	}
	//	c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
	//})
	route.GET("/get", func(c *gin.Context) {
		var person Person
		if err := c.BindQuery(&person); err != nil {
			fmt.Printf("bind error: %v", err)
			c.JSON(400, gin.H{"msg": err})
			return
		}
		decodeId, err := base64.StdEncoding.DecodeString(person.ID)
		if err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		c.JSON(200, gin.H{"name": person.Name, "uuid": string(decodeId)})
	})
	route.Run(":8088")
}

/*
curl -v localhost:8088/thinkerou/987fbc97-4bed-5078-9f07-9141ba07c9f3

curl -v localhost:8088/thinkerou/not-uuid
*/
