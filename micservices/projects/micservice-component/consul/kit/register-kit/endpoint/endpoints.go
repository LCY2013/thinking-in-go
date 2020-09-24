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
package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/hashicorp/consul/api"
	"register-kit/service"
)

// 注册端点结构体定义
type RegisterEndpoints struct {
	// 服务发现http端点
	DiscoveryEndpoint endpoint.Endpoint
	// 健康检测http端点
	HealthCheckEndpoint endpoint.Endpoint
}

// 定义服务发现请求结构体
type DiscoveryRequest struct {
	ServiceName string
}

// 定义服务发现响应结构体
type DiscoveryResponse struct {
	// 服务注册响应所有的实例信息
	Instance []*api.AgentService `json:"instances"`
	// 异常信息
	Error string `json:"error"`
}

// 创建服务发现端点
func CreateServiceDiscoverEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(DiscoveryRequest)
		instances, err := svc.DiscoverService(ctx, req.ServiceName)
		var errMsg = ""

		if err != nil {
			errMsg = err.Error()
		}

		return &DiscoveryResponse{
			Instance: instances,
			Error:    errMsg,
		}, nil
	}
}

// 定义健康检测请求结构体
type HealthCheckRequest struct {
}

// 定义健康检测响应结构体
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// 创建健康检测http端点
func CreateHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return HealthCheckResponse{
			Status: svc.HealthCheck(),
		}, nil
	}
}
