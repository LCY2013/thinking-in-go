package beego

import "github.com/beego/beego/v2/server/web"

type UserController struct {
	web.Controller
}

func (c *UserController) GetUser() {
	c.Ctx.WriteString("hello, i'm fufeng!")
}

func (c *UserController) CreateUser() {
	u := &User{}
	err := c.Ctx.BindJSON(u)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	err = c.Ctx.JSONResp(u)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
}

type User struct {
	Name string
}
