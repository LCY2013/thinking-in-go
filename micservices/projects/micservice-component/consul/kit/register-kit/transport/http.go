/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-09-24
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"register-kit/endpoint"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// 定义http请求参数不合法异常
	ErrBadRequest = errors.New("invalid request parameter")
)

// 创建http端点处理器
func CreateHttpHandler(ctx context.Context, endpoint *endpoint.RegisterEndpoints) http.Handler {

	// 创建路由
	router := mux.NewRouter()

	// 定义日志信息
	kitLogger := log.NewLogfmtLogger(os.Stderr)
	kitLogger = log.With(kitLogger, "ts", log.DefaultTimestampUTC)
	kitLogger = log.With(kitLogger, "caller", log.DefaultCaller)

	// 定义http options
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(kitLogger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	router.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoint.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeJsonResponse,
		options...,
	))

	router.Methods("GET").Path("/discovery/name").Handler(kithttp.NewServer(
		endpoint.DiscoveryEndpoint,
		decodeDiscoveryServiceRequest,
		encodeJsonResponse,
		options...,
	))

	return router
}

// 解码服务发现请求
func decodeDiscoveryServiceRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	serviceName := req.URL.Query().Get("serviceName")

	if serviceName == "" {
		return nil, ErrBadRequest
	}

	return endpoint.DiscoveryRequest{
		ServiceName: serviceName,
	}, nil
}

// 解码健康检测请求
func decodeHealthCheckRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	return endpoint.HealthCheckRequest{}, nil
}

// 编码服务响应请求
func encodeJsonResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(writer).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
