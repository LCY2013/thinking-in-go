/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-09-27
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

	"github.com/go-kit/kit/endpoint"
)

// 定义请求结构体
type RpcRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 定义响应结构体
type RpcResponse struct {
	Ret bool  `json:"ret"`
	Err error `json:"err"`
}

// 定义go kit 的端点
type Endpoints struct {
	UserEndpoint endpoint.Endpoint
}

// 创建User Endpoint
func MakeUserEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// 强制转换登陆请求结构体
		loginRequest := request.(*RpcRequest)
		// 调用业务函数
		isLogin, err := svc.CheckPassword(ctx,
			loginRequest.Username, loginRequest.Password)
		return &RpcResponse{Ret: isLogin, Err: err}, err
	}
}
