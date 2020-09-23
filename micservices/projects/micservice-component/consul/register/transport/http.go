/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 Ramostear.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-09-22
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
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kitHttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"regiter/endpoint"
)

var (
	// 定义错误请求错误
	ErrorBadRequest = errors.New("invalid request parameter")
)

// 构建一个http请求处理方法
func MakeHttpHandler(ctx context.Context, endpoint *endpoint.RegisterEndpoints) http.Handler {
	// 定义路由
	r := mux.NewRouter()

	// 定义系统标准错误输出
	kitLog := log.NewLogfmtLogger(os.Stderr)
	kitLog = log.With(kitLog, "ts", log.DefaultTimestampUTC)
	kitLog = log.With(kitLog, "caller", log.DefaultCaller)

	// 定义kit http 服务参数
	options := []kitHttp.ServerOption{
		kitHttp.ServerErrorHandler(transport.NewLogErrorHandler(kitLog)),
		kitHttp.ServerErrorEncoder(errorEncoder),
	}

	// 定义端点/health
	r.Methods("GET").Path("/health").Handler(kitHttp.NewServer(
		endpoint.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeJsonResponse,
		options...,
	))

	// 定义服务发现端点/discovery/name
	r.Methods("GET").Path("/discovery/name").Handler(kitHttp.NewServer(
		endpoint.DiscoveryEndpoint,
		decodeDiscoveryRequest,
		encodeJsonResponse,
		options...,
	))

	return r
}

// 定义解码服务发现的方法
func decodeDiscoveryRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	// 从url中的请求参数获取服务名称
	serviceName := req.URL.Query().Get("serviceName")
	if serviceName == "" {
		return nil, ErrorBadRequest
	}
	return endpoint.DiscoveryRequest{
		ServiceName: serviceName,
	}, nil
}

// 定义编码成json的响应方法
func encodeJsonResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(writer).Encode(response)
}

// 定义解码健康检测请求方法
func decodeHealthCheckRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	return endpoint.HealthCheckRequest{}, nil
}

// 错误编码方法
func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	// 设置http响应头的Content-Type
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
