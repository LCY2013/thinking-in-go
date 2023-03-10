package main

import "github.com/lcy2013/custom-web/coreweb/server/05/framework"

func UserLoginController(c *framework.Context) error {
	c.JsonResp(200, "ok, UserLoginController")
	return nil
}
