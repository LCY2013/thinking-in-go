/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-06-15
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package goroutines

import (
	"fmt"
	"time"
)

/*
Goroutines Go 程

我们称之为 Go 程 是因为现有的术语—线程、协程、进程等等—无法准确传达它的含义。
Go 程具有简单的模型:它是与其它 Go 程并发运行在同一地址空间的函数。
它是轻量级的，所有小号几乎就只有栈空间的分配。
而且栈最开始是非常小的，所以它们很廉价，仅在需要时才会随着堆空间的分配(和释放)而变化。

Go 程在多线程操作系统上可实现多路复用，因此若一个线程阻塞，比如说等待 I/O， 那么其它的线程就会运行。
Go 程的设计隐藏了线程创建和管理的诸多复杂性。

在函数或方法前添加 go 关键字能够在新的 Go 程中调用它。当调用完成后，该 Go 程也会安静地退出。
(效果有点像 Unix Shell 中的 & 符号，它能让命令在后台运行。)
go list.Sort() // 并发运行lsit.Sort()无需等待它结束

函数字面在 Go 程调用中非常有用。
func Announce(message string, delay time.Duration) {
    go func() {
        time.Sleep(delay)
		fmt.Println(message)
	}() // 注意括号 - 必须调用该函数。
}

在 Go 中，函数字面都是闭包:其实现在保证了函数内引用变量的生命周期与函数的活动时间相同。

这些函数没什么实用性，因为它们没有实现完成时的信号处理。因此，我们需要信道。
*/

func goroutines() {
	// 并发运行list.Sort()无需等待它结束
	// go list.Sort()
}

func Announce(message string, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		fmt.Println(message)
	}() // 注意括号 - 必须调用该函数
}
