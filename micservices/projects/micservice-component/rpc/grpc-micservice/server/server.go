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
	"context"
	"flag"
	"grpc-mic/pb"
	"grpc-mic/users"
	logGo "log"
	"net"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func main() {
	// 获取命令行参数信息
	flag.Parse()

	// 定义日志相关参数信息
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, log.DefaultTimestampUTC)
		logger = log.With(logger, log.DefaultCaller)
	}

	// 定义应用上下文信息
	ctx := context.Background()

	// 创建Service
	var svc users.UserService
	// 创建UserService实现
	svc = &users.UserServiceImpl{}

	// 构建日志中间件
	service := users.LoggingMiddleware(logger)(svc)

	// 创建Endpoint
	endpoint := users.MakeUserEndpoint(service)

	// 构建限流中间件
	limiter := rate.NewLimiter(rate.Every(time.Second*1), 1)
	endpoint = users.NewTokenBucketLimiterWithBuildIn(limiter)(endpoint)

	endpoints := users.Endpoints{
		UserEndpoint: endpoint,
	}

	// 构建UserService
	userServiceServer := users.NewUserServer(ctx, endpoints)

	// grpc 启动，监听端口、注册grpc服务信息
	listen, err := net.Listen("tcp", "localhost:9527")
	if err != nil {
		logGo.Printf("gpc listen err : %s\n", err)
	}
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, userServiceServer)
	_ = server.Serve(listen)
}
