package gin

import "github.com/gin-gonic/gin"

type UserController struct {
}

func (c *UserController) GetUser(ctx *gin.Context) {
	ctx.String(200, "hello, world")
}
