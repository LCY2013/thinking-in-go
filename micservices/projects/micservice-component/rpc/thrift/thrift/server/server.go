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
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
	"thrift/thrift/gen-go/user_service"
)

const (
	HOST = "localhost"
	PORT = "9527"
)

func main() {
	// 创建结构体方法
	userServiceHandler := &user_service.UserService{}

	// 创建一个具体的处理器，将自定义实现Handler绑定到处理器上
	userProcessor := user_service.NewUserProcessor(userServiceHandler)

	// 创建ServerSocket
	serverSocket, err := thrift.NewTServerSocket(HOST + ":" + PORT)
	if err != nil {
		log.Panicf("create server socket error : %s \n", err)
	}

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	factoryDefault := thrift.NewTBinaryProtocolFactoryDefault()

	server4 := thrift.NewTSimpleServer4(userProcessor, serverSocket,
		transportFactory, factoryDefault)

	fmt.Println("thrift server running at : ", HOST, ":", PORT)

	server4.Serve()

}
