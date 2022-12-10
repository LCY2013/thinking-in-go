package middleware

import (
	"fmt"

	"github.com/lcy2013/custom-web/coreweb/server/04/framework"
)

func Test1() framework.ControllerHandler {
	// 使用函数回调
	return func(c *framework.Context) error {
		fmt.Println("middleware pre test1")
		if err := c.Next(); err != nil {
			return err
		}
		fmt.Println("middleware post test1")
		return nil
	}
}

func Test2() framework.ControllerHandler {
	// 使用函数回调
	return func(c *framework.Context) error {
		fmt.Println("middleware pre test2")
		if err := c.Next(); err != nil {
			return err
		}
		fmt.Println("middleware post test2")
		return nil
	}
}

func Test3() framework.ControllerHandler {
	// 使用函数回调
	return func(c *framework.Context) error {
		fmt.Println("middleware pre test3")
		if err := c.Next(); err != nil {
			return err
		}
		fmt.Println("middleware post test3")
		return nil
	}
}
