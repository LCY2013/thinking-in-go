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
	"log"
	discoverys "user/discovery"
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

	flag.Parse()

	// 创建DiscoveryClient
	client := discoverys.NewDiscoveryClient(*consulAddr, *consulPort)

	ctx := context.Background()
	// 错误通知的channel
	errChan := make(chan error)

	errorChan := <-errChan
	log.Println("error from chan  : ", errorChan)
}
