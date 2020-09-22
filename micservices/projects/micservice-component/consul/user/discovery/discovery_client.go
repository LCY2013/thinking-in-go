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
package discoverys

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// 服务实例结构体定义
type InstanceInfo struct {
	ID                string                    `json:"ID"`                // 服务实例ID
	Service           string                    `json:"Service,omitempty"` // 服务发现时返回服务名称
	Name              string                    `json:"Name"`              // 服务名称
	Tags              []string                  `json:"Tags,omitempty"`    // 标签用于服务过滤
	Address           string                    `json:"Address"`           // 服务实例Host
	Port              int                       `json:"Port"`              // 服务实例端口号
	Meta              map[string]string         `json:"Meta,omitempty"`    // 元信息
	EnableTagOverride bool                      `json:"EnableTagOverride"` // 是否运行标签覆盖
	Check             `json:"Check,omitempty"`  // 健康检测
	Weight            `json:"Weight,omitempty"` // 权重
}

// 定义健康检测结构体
type Check struct {
	DeregisterCriticalServiceAfter string   `json:"DeregisterCriticalServiceAfter"` // 多久之后注销服务
	Args                           []string `json:"Args,omitempty"`                 // 请求参数
	HTTP                           string   `json:"HTTP"`                           // 健康检查地址
	Interval                       string   `json:"Interval,omitempty"`             // consul 主动检查间隔
	TTL                            string   `json:"TTL,omitempty"`                  // 服务实例主动维持心跳间隔，与Interval只存其一
}

// 定义权重结构体
type Weight struct {
	Passing int `json:"Passing"`
	Warning int `json:"Warning"`
}

// 定义服务发现客户端
type DiscoveryClient struct {
	host string // consul host
	port int    // consul port
}

// 创建一个DiscoveryClient
func NewDiscoveryClient(host string, port int) *DiscoveryClient {
	return &DiscoveryClient{
		host: host,
		port: port,
	}
	// return &DiscoveryClient{host,port}
}

// 注册实例服务实例
func (consulClient *DiscoveryClient) Register(ctx context.Context,
	serviceName, instanceId, healCheckUrl, instanceHost string,
	instancePort int,
	meta map[string]string,
	weight *Weight) error {

	// 构建一个服务实例相关信息
	instanceInfo := &InstanceInfo{
		ID:                instanceId,
		Name:              serviceName,
		Address:           instanceHost,
		Port:              instancePort,
		Meta:              meta,
		EnableTagOverride: false,
		Check: Check{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healCheckUrl,
			Interval:                       "15s",
		},
	}

	// 判断权重是否为nil
	if weight == nil {
		instanceInfo.Weight = Weight{
			Passing: 10,
			Warning: 1,
		}
	} else {
		instanceInfo.Weight = *weight
	}

	// 使用json序列化服务实例
	byteDataInstance, err := json.Marshal(instanceInfo)
	if err != nil {
		log.Printf("json format err : %s\n", err)
		return err
	}

	// 向consul注册服务实例
	req, err := http.NewRequest("PUT",
		"http://"+consulClient.host+":"+strconv.Itoa(consulClient.port)+"/v1/agent/service/register",
		bytes.NewReader(byteDataInstance))
	if err != nil {
		return err
	}
	// 设置请求头信息
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	client.Timeout = time.Second * 2
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("register instance service : %s\n", err)
		return err
	}
	// 最后关闭响应返回
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("register service http request errCode : %v\n", resp.StatusCode)
		return fmt.Errorf("register service http request errCode : %v\n", resp.StatusCode)
	}

	log.Printf("register service %s success\n", serviceName)
	return nil
}

func (consulClient *DiscoveryClient) DiscoveryServices(ctx context.Context,
	serviceName string) ([]*InstanceInfo, error) {

	// 构建查询请求
	req, err := http.NewRequest("GET",
		"http://"+consulClient.host+":"+strconv.Itoa(consulClient.port)+"/v1/health/service/"+serviceName,
		nil)
	if err != nil {
		log.Printf("http request format error : %s\n", err)
		return nil, err
	}

	client := http.Client{}
	client.Timeout = time.Second * 2

	// 发送http请求
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("discovery service err :%s\n", err)
		return nil, err
	}

	// 关闭http响应体
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("discover service http request errCode : %v\n", resp.StatusCode)
		return nil, fmt.Errorf("discover service http request errCode : %v\n", resp.StatusCode)
	}

	// 定义服务列表结构体
	var serviceList []struct {
		Service InstanceInfo `json:"Service"`
	}

	// 序列化解析响应体
	err = json.NewDecoder(resp.Body).Decode(&serviceList)
	if err != nil {
		log.Printf("format service info err : %s\n", err)
		return nil, fmt.Errorf("format service info err : %s\n", err)
	}

	// 构建一个切片容器用来存储json解析后的数据
	instances := make([]*InstanceInfo, len(serviceList))
	for inx, service := range serviceList {
		instances[inx] = &service.Service
	}

	return instances, nil
}
