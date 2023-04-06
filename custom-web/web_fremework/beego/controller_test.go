package beego

import (
	"github.com/beego/beego/v2/server/web"
	"testing"
)

func TestUserController(t *testing.T) {
	// go func() {
	// 	s := web.NewHttpSever()
	// 	s.Run(":8082")
	// }()
	web.BConfig.CopyRequestBody = true
	c := &UserController{}
	web.Router("/user", c, "get:GetUser")
	// 监听8080端口
	web.Run(":8080")
}
