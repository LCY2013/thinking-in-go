package main

import "github.com/lcy2013/custom-web/coreweb/server/04/framework"

func UserLoginController(c *framework.Context) error {
	c.Json(200, "ok, UserLoginController")
	return nil
}
