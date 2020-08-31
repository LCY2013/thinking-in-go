/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 Ramostear.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-08-31
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kitHttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	goLog "log"
	"net/http"
	"os"
	"user/endpoint"
)

var (
	ErrBadRequest = errors.New("invalid request parameter")
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints *endpoint.UserEndpoints) http.Handler {
	// 创建一个路由
	r := mux.NewRouter()

	// 创建日志输出信息
	kitLogger := log.NewLogfmtLogger(os.Stderr)

	// 绑定日志参数信息
	kitLogger = log.With(kitLogger, "ts", log.DefaultTimestampUTC)
	kitLogger = log.With(kitLogger, "caller", log.DefaultCaller)

	// 创建服务参数嘻嘻
	options := []kitHttp.ServerOption{
		kitHttp.ServerErrorHandler(transport.NewLogErrorHandler(kitLogger)),
		kitHttp.ServerErrorEncoder(encodeError),
	}

	// 开始绑定服务路由信息
	r.Methods("POST").
		Path("/register").
		Handler(kitHttp.NewServer(
			endpoints.RegisterEndpoint,
			decodeRegisterRequest,
			encodeJsonResponse,
			options...,
		))

	r.Methods("POST").
		Path("/login").
		Handler(kitHttp.NewServer(
			endpoints.LoginEndpoint,
			decodeLoginRequest,
			encodeJsonResponse,
			options...,
		))

	return r
}

// 编码响应的json信息
func encodeJsonResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return json.NewEncoder(writer).Encode(response)
}

// 解码登陆请求信息
func decodeLoginRequest(_ context.Context, req *http.Request) (interface{}, error) {
	// 获取参数信息
	email := req.FormValue("email")
	password := req.FormValue("password")
	// 参数校验
	if email == "" || password == "" {
		return req, ErrBadRequest
	}

	// 可以利用这里做json -> DTO
	// 读取body信息
	if req.Body != nil {
		//包装bufio
		bodyReader := bufio.NewReader(req.Body)
		content, err := ioutil.ReadAll(bodyReader)
		if err == nil {
			contentStr := fmt.Sprintf("%s", content)
			goLog.Println("json : " + contentStr)
		}
		_ = req.Body.Close()
	}

	return &endpoint.LoginRequest{
		Email:    email,
		Password: password,
	}, nil
}

// 解码注册请求信息
func decodeRegisterRequest(_ context.Context, req *http.Request) (interface{}, error) {
	// 获取参数信息
	email := req.FormValue("email")
	username := req.FormValue("username")
	password := req.FormValue("password")
	// 参数校验
	if email == "" || username == "" || password == "" {
		return req, ErrBadRequest
	}
	return &endpoint.RegisterRequest{
		Email:    email,
		Username: username,
		Password: password,
	}, nil
}

// 创建一个error encode处理器
func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

/*
	测试流程如下:
	1、curl -X POST http://127.0.0.1:9527/login\?email\=magic@fufeng.com\&password\=123456
	2、curl -X POST http://127.0.0.1:9527/register\?email\=magic@fufeng.com\&password\=123456\&username\=fufeng
*/
