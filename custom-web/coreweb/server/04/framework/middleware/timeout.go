package middleware

import (
	"context"
	"fmt"
	"github.com/lcy2013/custom-web/coreweb/server/04/framework"
	"log"
	"time"
)

func Timeout(d time.Duration) framework.ControllerHandler {
	// 使用函数回调
	return func(c *framework.Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)
		// 执行业务逻辑前预操作：初始化超时context
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			// 使用next执行具体的业务逻辑
			if err := c.Next(); err != nil {
				return
			}

			finish <- struct{}{}
		}()
		// 执行业务逻辑后操作
		select {
		case p := <-panicChan:
			c.Json(500, "time out")
			log.Println(p)
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.Json(500, "time out")
		}
		return nil
	}
}
