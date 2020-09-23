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
package service

import (
	"context"
	"errors"
	"log"
	discoverys "regiter/discovery"
)

// 定义一个service的接口
type Service interface {
	HealthCheck() string // 健康检查方法定义

	DiscoveryService(ctx context.Context, serviceName string) ([]*discoverys.InstanceInfo, error) // 定义服务发现方法
}

// 定义服务发现的错误
var ErrorNotServiceInstances = errors.New("service instance not found")

// 定义注册服务结构体
type RegisterServiceImpl struct {
	discoveryClient *discoverys.DiscoveryClient
}

// 生成一个Service 实现RegisterServiceImpl
func NewRegisterServiceImpl(discoveryClient *discoverys.DiscoveryClient) Service {
	return &RegisterServiceImpl{
		discoveryClient: discoveryClient,
	}
}

// 服务发现
func (service *RegisterServiceImpl) DiscoveryService(ctx context.Context, serviceName string) ([]*discoverys.InstanceInfo, error) {
	services, err := service.discoveryClient.DiscoveryServices(ctx, serviceName)

	if err != nil {
		log.Printf("get service info err : %s\n", err)
	}

	if service == nil || len(services) == 0 {
		return nil, ErrorNotServiceInstances
	}

	return services, nil
}

// 健康检查接口实现
func (service *RegisterServiceImpl) HealthCheck() string {
	return "OK"
}
