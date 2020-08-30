/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2019 Ramostear.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-08-16
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
	"fmt"
	"time"
)

// 定义数据库信息
const DB_ADDRESS = "127.0.0.1"
const CALCULATE = "CALCULATE_VALUE"

// 读数据库操作
func readDB(ctx context.Context, cost time.Duration) {
	fmt.Println("db address is", ctx.Value(DB_ADDRESS))
	select {
	case <-time.After(cost): // 模拟数据库读取数据
		fmt.Println("read data from db")
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // 任务取消原因
		// clear resource 清理资源
	}
}

// 模拟一些计算任务
func calculate(ctx context.Context, cost time.Duration) {
	fmt.Println("calculate value is : ", ctx.Value(CALCULATE))
	select {
	case <-time.After(cost): // 模拟数据计算的时间
		fmt.Println("calculate finish")
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // 任务取消原因
		// clear resource 清理资源
	}
}

func main() {
	// 创建一个应用上下文
	ctx := context.Background()
	// 添加上下文信息
	ctx = context.WithValue(ctx, DB_ADDRESS, "127.0.0.1:3306")
	ctx = context.WithValue(ctx, CALCULATE, 9527)
	// 设定子context 2s后执行超时
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*2)
	defer cancelFunc()
	// 设定执行时间为3s
	go readDB(ctx, time.Second*3)
	go calculate(ctx, time.Second*2)

	// 主goroutine休眠5秒
	time.Sleep(time.Second * 5)
}
