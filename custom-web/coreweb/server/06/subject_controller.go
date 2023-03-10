package main

import (
	"github.com/lcy2013/custom-web/coreweb/server/06/framework"
	"time"
)

func SubjectAddController(c *framework.Context) error {
	c.JsonResp(200, "ok, SubjectAddController")
	return nil
}

func SubjectListController(c *framework.Context) error {
	c.JsonResp(200, "ok, SubjectListController")
	return nil
}

func SubjectDelController(c *framework.Context) error {
	c.JsonResp(200, "ok, SubjectDelController")
	return nil
}

func SubjectUpdateController(c *framework.Context) error {
	c.JsonResp(200, "ok, SubjectUpdateController")
	return nil
}

func SubjectGetController(c *framework.Context) error {
	c.JsonResp(200, "ok, SubjectGetController")
	return nil
}

func SubjectNameController(c *framework.Context) error {
	c.JsonResp(200, "ok, SubjectNameController")
	return nil
}

func SubjectGracefulShutdownController(c *framework.Context) error {
	foo, _ := c.QueryString("foo", "def")
	// 等待10s才结束执行
	time.Sleep(10 * time.Second)
	// 输出结果
	c.SetOkStatus().Json("ok, graceful shutdown: " + foo)
	return nil
}
