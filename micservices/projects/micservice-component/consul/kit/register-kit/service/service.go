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
package service

import (
	"context"
	"errors"
	"log"
	discoverys "register-kit/discovery"

	"github.com/hashicorp/consul/api"
)

// 定义一个服务接口
type Service interface {
	// 健康检测方法定义
	HealthCheck() string
	// 服务发现方法定义
	DiscoverService(ctx context.Context, serviceName string) ([]*api.AgentService, error)
}

// 定义服务实例没找到异常
var ErrNotServiceInstances = errors.New("services instance not found")

// 定义服务注册实现
type RegisterServiceImpl struct {
	discoveryClient *discoverys.DiscoveryClient
}

// 创建一个结构体实现
func CreateRegisterServiceImpl(discoveryClient *discoverys.DiscoveryClient) Service {
	return &RegisterServiceImpl{
		discoveryClient: discoveryClient,
	}
}

// 结构体RegisterServiceImpl实现接口Service方法

func (*RegisterServiceImpl) HealthCheck() string {
	return "OK"
}

func (service *RegisterServiceImpl) DiscoverService(ctx context.Context, serviceName string) ([]*api.AgentService, error) {
	// 获取服务名称为ServiceName的所有服务实例
	services, err := service.discoveryClient.DiscoverServices(ctx, serviceName)
	if err != nil {
		log.Printf("get [%s] service instance error : %s\n", serviceName, err)
	}
	if services == nil || len(services) == 0 {
		return nil, ErrNotServiceInstances
	}
	return services, nil
}
