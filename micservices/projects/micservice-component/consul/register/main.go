/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 Ramostear.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-08-30
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	discoverys "user/discovery"
	"user/endpoint"
	"user/service"
	"user/transport"
)

// 程序主函数信息
func main() {
	var (
		// waitTime = flag.Int("wait.time", 10, "wait time")
		// consul address
		consulAddr = flag.String("consul.addr", "localhost", "consul address")
		// consul port
		consulPort = flag.Int("consul.port", 8500, "consul port")
		// service name
		serviceName = flag.String("service.name", "register", "service name")
		// service addr
		serviceAddr = flag.String("service.addr", "localhost", "service addr")
		// 服务地址、服务名称
		servicePort = flag.Int("service.port", 9527, "service port")
	)

	// 通过命令行参数或者环境变量获取服务配置信息
	flag.Parse()

	// 创建DiscoveryClient
	client := discoverys.NewDiscoveryClient(*consulAddr, *consulPort)

	ctx := context.Background()
	// 错误通知的channel
	errChan := make(chan error)

	// 获取注册服务实现
	srv := service.NewRegisterServiceImpl(client)

	// 定义注册需要的端点
	endpoints := endpoint.RegisterEndpoints{
		DiscoveryEndpoint:   endpoint.MakeDiscoveryEndpoint(srv),
		HealthCheckEndpoint: endpoint.MakeHealCheckEndpoint(srv),
	}

	// http 处理器构建
	handler := transport.MakeHttpHandler(ctx, &endpoints)

	// 协程异步监听http服务
	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		// 监控系统信号，等待Ctrl+C 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// 创建一个服务实例ID
	instanceId := *serviceName + "-" + uuid.New().String()

	// 注册服务
	err := client.Register(ctx, *serviceName, instanceId,
		"/health", *serviceAddr, *servicePort,
		nil, nil)

	if err != nil {
		log.Printf("register service err : %s\n", err)
		os.Exit(-1)
	}

	errorMsg := <-errChan
	log.Printf("listen error : %s\n", errorMsg)

	// 取消注册
	err = client.Deregister(ctx, instanceId)

}

func init() {
	// 定义注册服务的日志文件
	file := "./" + "register.log"
	openFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	// 将文件设置为log输出文件
	log.SetOutput(openFile)
	log.SetFlags(log.Ldate | log.Lshortfile)
}
