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
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
	"net"
	"thrift/thrift/gen-go/user_service"
)

const (
	HOST = "localhost"
	PORT = "9527"
)

func main() {

	// 创建thrift的tSocket
	tSocket, err := thrift.NewTSocket(net.JoinHostPort(HOST, PORT))
	if err != nil {
		log.Panicf(" thrift listen error : %v \n", err)
	}

	// 创建thrift传输工厂
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())

	transport, err := transportFactory.GetTransport(tSocket)
	if err != nil {
		log.Fatalln("get transport error : ", err)
	}

	factoryDefault := thrift.NewTBinaryProtocolFactoryDefault()

	userClient := user_service.NewUserClientFactory(transport, factoryDefault)

	if err := transport.Open(); err != nil {
		log.Fatalln("error open : ", HOST, ":", PORT)
	}

	// 关闭传输通道
	defer transport.Close()

	// 构建请求
	loginRequest := &user_service.LoginRequest{
		Username: "fufeng",
		Password: "12345",
	}

	// 检测密码是否正常
	loginResponse, err := userClient.CheckPassword(context.Background(), loginRequest)

	if err != nil {
		log.Fatalf("user CheckPassword call error : %s \n", err)
	}

	fmt.Println(loginResponse)
}
