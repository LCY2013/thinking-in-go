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
package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	discoverys "user/discovery"
	"user/service"
)

// 注册端点结构体
type RegisterEndpoints struct {
	DiscoveryEndpoint   endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

// 服务发现请求结构体
type DiscoveryRequest struct {
	ServiceName string
}

// 服务发现响应结构体
type DiscoveryResponse struct {
	Instances []*discoverys.InstanceInfo `json:"instances"`
	Error     string                     `json:"error"`
}

// 健康检查请求HealthRequest结构体
type HealthCheckRequest struct {
}

// 健康检查响应HealthResponse结构体
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// endpoint 设置

// 创建服务发现的endpoint
func MakeDiscoveryEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		// 强制转换任意类型的请求到服务发现请求结构体
		req := request.(DiscoveryRequest)
		// 通过服务发现去consul获取服务实例
		discoveryServices, err := svc.DiscoveryService(ctx, req.ServiceName)

		// 定义错误消息内容
		var errMsg = ""
		if err != nil {
			errMsg = err.Error()
		}
		// 返回服务发现响应结构体
		return DiscoveryResponse{
			Instances: discoveryServices,
			Error:     errMsg,
		}, nil
	}
}

// 创建健康检查的endpoint
func MakeHealCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		// 利用服务去查询健康状态
		checkStatus := svc.HealthCheck()
		// 返回健康检查响应结构体
		return HealthCheckResponse{
			Status: checkStatus,
		}, nil
	}
}
