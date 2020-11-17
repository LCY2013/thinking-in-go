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
	"fmt"
	"runtime"
	"time"
)

// 计算1-100亿相加的结果
// CSP 模型与生产者 - 消费者模式
// CSP 理论书籍 http://www.usingcsp.com/cspbook.pdf
func main() {
	v1(10000000000)
	v2(10000000000)
}

// 计算cpu核数，启动对应的线程
func v1(num uint64) {
	now := time.Now()
	// 获取cpu核数
	cpuNum := runtime.NumCPU()
	// 创建一个chan用于接受计算结果
	sumChan := make(chan uint64, cpuNum)
	// 计算数据的中间间隔
	seed := num / uint64(cpuNum)
	fmt.Println(seed, " - ", cpuNum)
	cur := uint64(0)
	for i := 0; i < cpuNum; i++ {
		go calc(cur, cur+seed, sumChan)
		cur += seed
	}
	// 统计结果集
	sum := uint64(0)

	for i := 0; i < cpuNum; i++ {
		sum += <-sumChan
	}

	fmt.Println(time.Since(now), " v1 计算结果: ", sum)
}

func v2(num uint64) {
	now := time.Now()
	// 创建一个chan用于接受计算结果
	sumChan := make(chan uint64)
	go calc(0, num, sumChan)

	fmt.Println(time.Since(now), " v2 计算结果: ", <-sumChan)
}

// 计算uint64 from -> to 相加
func calc(from, to uint64, sumChan chan<- uint64) {
	result := uint64(0)
	for ; from < to; from++ {
		result += from
	}
	// 将结果写入chan中
	sumChan <- result
}

/*
CSP 模型与 Actor 模型的区别

第一个区别是：
		Actor 模型中没有 channel。虽然 Actor 模型中的 mailbox 和 channel 非常像，看上去都像个 FIFO 队列，但是区别还是很大的。
	Actor 模型中的 mailbox 对于程序员来说是“透明”的，mailbox 明确归属于一个特定的 Actor，是 Actor 模型中的内部机制；
	而且 Actor 之间是可以直接通信的，不需要通信中介。
	但 CSP 模型中的 channel 就不一样了，它对于程序员来说是“可见”的，是通信的中介，传递的消息都是直接发送到 channel 中的。

第二个区别是：
	Actor 模型中发送消息是非阻塞的，而 CSP 模型中是阻塞的。
	Golang 实现的 CSP 模型，channel 是一个阻塞队列，当阻塞队列已满的时候，向 channel 中发送数据，会导致发送消息的协程阻塞。
*/
