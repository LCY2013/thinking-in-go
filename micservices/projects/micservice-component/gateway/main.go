/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-09-28
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"
)

func main() {

	// 获取命令行环境参数
	var (
		consulHost = flag.String("consul.host", "127.0.0.1", "consul server ip address")
		consulPort = flag.String("consul.port", "127.0.0.1", "consul server port")
	)
	// 处理环境参数
	flag.Parse()

	// 创建日志组件，设置日志组件相关内容
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, log.DefaultTimestampUTC)
		logger = log.With(logger, log.DefaultCaller)
	}

	// 创建consul api 的客户端信息
	consulConfig := api.DefaultConfig()
	// 设置consul的访问地址
	consulConfig.Address = "http://" + *consulHost + ":" + *consulPort
	// 利用consulConfig创建一个consul客户端
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		_ = logger.Log("log", err)
		// 定制退出码，用于容器退出捕获
		os.Exit(7)
	}

	// 创建反向代理，通过传入consul的client端以及日志组件
	proxy := NewReverseProxy(consulClient, logger)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// 启动协程监听端口网络服务
	go func() {
		_ = logger.Log("transport", "HTTP", "addr", "9527")
		errs <- http.ListenAndServe(":9527", proxy)
	}()

	// 等待退出指令并且打印退出请求
	logger.Log("exit", <-errs)
}

// 利用go提供的httputils.ReverseProxy 实现一个简单的反向代理，通过查询consul的服务实例
func NewReverseProxy(consulClient *api.Client, logger log.Logger) *httputil.ReverseProxy {

	// 创建ReverseProxy需要的Director
	director := func(req *http.Request) {
		// 获取请求原始路径
		reqPath := req.URL.Path
		if reqPath == "" {
			return
		}

		// 按照"/"对路径进行分割，获取到实例名称ServiceName
		reqs := strings.Split(reqPath, "/")
		serviceName := reqs[1]

		// 根据服务名称去Consul查询所有的服务实例
		services, _, err := consulClient.Catalog().Service(serviceName, "", nil)
		if err != nil {
			_ = logger.Log("reverse proxy fail", "no such service instance", err.Error())
			return
		}

		// 判断服务实例的数量
		if len(services) == 0 {
			_ = logger.Log("reverse proxy fail", "no such service instance", serviceName)
			return
		}

		// 重新组织请求路径，去掉原有的ServiceName
		realRequestPath := strings.Join(reqs[2:], "/")

		// 随机选择一个服务实例
		targetInstance := services[rand.Int()%len(services)]
		// 添加日志打印，输出获取到的目标实例ID
		logger.Log("service id", targetInstance.ServiceID)

		// 设置代理相关的配置信息
		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("%s:%d", targetInstance.ServiceAddress, targetInstance.ServicePort)
		req.URL.Path = realRequestPath
	}

	// 返回一个反向代理实现
	return &httputil.ReverseProxy{
		Director: director,
	}
}
