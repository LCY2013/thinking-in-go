/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2019 Ramostear.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-08-13
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
	"fmt"
	"time"
)

func sender(begin int, ch chan int) {
	// 开始循环向通道发送消息
	for i := begin; i < begin+10; i++ {
		ch <- i
	}
}

// 接受一次通道消息
func receiver(ch chan int) {
	for i := 0; i < 2; i++ {
		value := <-ch
		fmt.Println("receiver : ", value)
	}
}

func main() {
	// 定义一个int类型通道
	ch1 := make(chan int)
	ch2 := make(chan int)
	// 开始发送
	go sender(10, ch1)
	// 接受者
	go receiver(ch2)

	// 主goroutine休眠一秒，保证调度成功
	time.Sleep(time.Second)

	// 一直执行，等待超时
	for {
		// select 多路选择器
		select {
		case v1 := <-ch1: // 使用ch1接受数据
			fmt.Println("select receiver : ", v1)
		case ch2 <- 2: // 使用ch2发送数据
			fmt.Println("send to ch2")
		case <-time.After(time.Second * 2): // 超时设置
			fmt.Println("timeout")
			return
		}
	}
}
