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
package discoveries

import (
	"context"
	"os"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

// 定义服务发现结构体
type DiscoveryClient struct {
	client       consul.Client
	register     sd.Registrar
	config       *api.Config
	registration *api.AgentServiceRegistration
}

// 定义生成consul注册信息的方法
func NewAgentServiceRegistration(
	serviceName, instanceId, healthCheckUrl, serviceAddr string,
	servicePort int,
	meta map[string]string) *api.AgentServiceRegistration {
	return &api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    serviceName,
		Address: serviceAddr,
		Port:    servicePort,
		Meta:    meta,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + serviceAddr + ":" + strconv.Itoa(servicePort) + healthCheckUrl,
			Interval:                       "15s",
		},
	}
}

// 创建服务发现客户端
func CreateDiscoveryClient(host string, port int, registration *api.AgentServiceRegistration) (*DiscoveryClient, error) {
	// 获取默认配置信息
	config := api.DefaultConfig()

	// 设置配置的地址
	config.Address = host + ":" + strconv.Itoa(port)

	// 创建客户端
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	// 创建consul client
	consulClient := consul.NewClient(client)

	// 返回自定义的服务发现结构体
	return &DiscoveryClient{
		client:       consulClient,
		register:     consul.NewRegistrar(consulClient, registration, log.NewLogfmtLogger(os.Stderr)),
		config:       config,
		registration: registration,
	}, nil
}

// 定义服务注册方法
func (discoveryClient *DiscoveryClient) Register(ctx context.Context) {
	discoveryClient.register.Register()
}

// 定义服务取消注册方法
func (discoveryClient *DiscoveryClient) Deregister(ctx context.Context) {
	discoveryClient.register.Deregister()
}

// 定义服务发现方法
func (discoveryClient *DiscoveryClient) DiscoverServices(ctx context.Context, serviceName string) ([]*api.AgentService, error) {
	// 调用具体的服务发现客户端去发现服务
	serviceEntries, _, err :=
		discoveryClient.client.Service(serviceName, "", false, nil)
	if err != nil {
		return nil, err
	}
	// 创建一个服务返回切片
	services := make([]*api.AgentService, 0, len(serviceEntries))
	for _, serviceEntry := range serviceEntries {
		services = append(services, serviceEntry.Service)
	}
	// 返回获取到的服务
	return services, nil
}
