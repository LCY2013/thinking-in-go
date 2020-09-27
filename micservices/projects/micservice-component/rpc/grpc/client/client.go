/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-09-27
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
	"google.golang.org/grpc"
	"grpc-project/pb"
	"log"
)

func main() {

	// grpc 服务端地址
	serviceAddress := "127.0.0.1:9527"
	// 创建目标连接的客户端
	clientConn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc connction err : %s\n", err)
	}
	// 最后关闭连接
	defer clientConn.Close()

	// 创建一个客户端连接服务
	userClient := pb.NewUserServiceClient(clientConn)

	// 定义请求参数结构体
	loginRequest := &pb.LoginRequest{
		Username: "fufeng",
		Password: "123456",
	}

	// 调用远程服务
	loginResponse, err := userClient.CheckPassword(context.Background(), loginRequest)
	if err != nil {
		log.Printf("grpc erro : %s\n", err)
	}

	log.Printf("grpc call method CheckPassword return : %s\n", loginResponse.Ret)
}
