/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-09-28
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package users

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

// 定义类型服务中间件为函数类型
type ServiceMiddleware func(service UserService) UserService

// 定义日志中间件结构体
type loggingMiddleware struct {
	UserService
	logger log.Logger
}

// 定义函数含有日志组件的服务中间件
func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next UserService) UserService {
		return &loggingMiddleware{
			next, logger,
		}
	}
}

// 实现接口方法
func (logMiddleware *loggingMiddleware) CheckPassword(ctx context.Context, userName, password string) (bool, error) {
	// 最后记录方法调用时间
	defer func(begin time.Time) {
		logMiddleware.logger.Log(
			"function", "CheckPassword",
			"userName", userName,
			"password", password,
			"took", time.Since(begin),
		)
	}(time.Now())
	// 具体服务调用
	isOk, err := logMiddleware.UserService.CheckPassword(ctx, userName, password)
	return isOk, err
}
