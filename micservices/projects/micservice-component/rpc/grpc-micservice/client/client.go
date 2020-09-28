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
	"fmt"
	"google.golang.org/grpc"
	"grpc-mic/pb"
	"time"
)

func main() {
	for i := 0; i < 5; i++ {
		go startClientRpc()
	}
	// 阻塞等待应用响应
	time.Sleep(time.Second * 1)
}

func startClientRpc() {
	// 定义RPC服务地址
	serviceAddress := "localhost:9527"
	// 创建客户端连接
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("grpc client connect err")
	}
	defer conn.Close()

	userServiceClient := pb.NewUserServiceClient(conn)
	ret, err := userServiceClient.CheckPassword(context.Background(), &pb.LoginRequest{
		Username: "fufeng",
		Password: "123456",
	})
	if err != nil {
		// panic(err)
		fmt.Printf("call rpc service err : %s\n", err)
	} else {
		fmt.Println("check password status : ", ret.Ret)
	}
}
