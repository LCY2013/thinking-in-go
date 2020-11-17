/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-11-17
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
	"log"
	"net"
	"strings"
	"time"
)

// 利用Thread-Per-Message模式
func main() {
	// 监听本地端口号
	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Panicln("echo server 监听端口启动失败", err)
	}

	go start(listen)

	log.Println("echo server started port on ", 8080)

	// 服务主协程休眠等待
	time.Sleep(time.Minute * 10000000)
}

// 启动echo 服务
func start(listen net.Listener) {
	// 最后关闭监听
	defer listen.Close()

	for {
		// 处理连接请求
		conn, err := listen.Accept()

		if err != nil {
			log.Panicln(err)
			continue
		}

		// 处理连接成功的连接
		go handlerTcp(conn)
	}
}

// 利用协程处理tcp连接处理
func handlerTcp(conn net.Conn) {
	// 处理完成关闭连接
	defer conn.Close()

	for {
		// 创建一个用于接受客户端数据大小为1024的buffer
		buf := make([]byte, 1024)
		// 读取请求数据
		readIndex, err := conn.Read(buf)
		if err != nil {
			log.Panicln(err)
		}
		//log.Println(bufio.Reader{})
		var content = string(buf[:readIndex])
		log.Println(content)
		if strings.Compare(strings.TrimSpace(content), "quit") == 0 {
			// 原样返回客户端数据
			_, _ = conn.Write([]byte("bye"))
			break
		}
		// 原样返回客户端数据
		_, _ = conn.Write(buf[:readIndex])
	}
}
