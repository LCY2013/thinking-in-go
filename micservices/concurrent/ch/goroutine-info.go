/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2019 Ramostear.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-08-12
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

// 生产者
func Producer(begin, end int, ch chan int) {
	for i := begin; i < end; i++ {
		ch <- i
	}
}

// 消费者
func Consumer(indexName int, ch chan int) {
	for value := range ch {
		fmt.Printf("consumer %d recover : %v \n", indexName, value)
	}
}

func main() {
	// 创建一个chan实例
	ch := make(chan int)
	// 最后关闭通道
	defer close(ch)

	// 三个生产者
	for i := 0; i < 3; i++ {
		go Producer(i*7, (i+1)*7, ch)
	}
	// 两个消费者
	for i := 0; i < 2; i++ {
		go Consumer(i, ch)
	}

	// 主协程休眠一秒,等待其他协程结束
	time.Sleep(time.Second)
}
