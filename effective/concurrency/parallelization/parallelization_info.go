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
package parallelization

import "runtime"

/*
这些设计的另一个应用是在多 CPU 核心上实现并行计算。
如果计算过程能够被分为几块可独立执行的过程，它就可以在每块计算结束时向信道发送信号，从而实现并行处理。

让我们看看这个理想化的例子。我们在对一系列向量项进行极耗资源的操作， 而每个项的值计算是完全独立的。
type Vector []float64

//DoSome 将此操应用至 v[i], v[i+1] ... 直到 v[n-1]
func (v Vector) DoSome(i, n int, u Vector, c chan int) {
	for ; i < n; i++ {
		v[i] += u.Op(v[i])
	}
	c <- 1 // 发信号表示这一块计算完成。
}

我们在循环中启动了独立的处理块，每个 CPU 将执行一个处理。 它们有可能以乱序的形式完成并结束，但这没有关系;
我们只需在所有 Go 程开始后接收，并统计信道中的完成信号即可。
*/

type Vector []float64

// DoSome 将此操应用至 v[i], v[i+1] ... 直到 v[n-1]
func (v Vector) DoSome(i, n int, u Vector, c chan int) {
	for ; i < n; i++ {
		v[i] += u.Op(v[i])
	}
	c <- 1 // 发信号表示这一块计算完成。
}

func (v Vector) Op(num float64) float64 {
	return num
}

// const CPU = runtime.NumCPU() // CPU 核心数

func (v Vector) DoAll(u Vector) {
	CPU := runtime.NumCPU() // CPU 核心数
	// 告诉运行时希望同时有多少 Go 程能执行代码
	runtime.GOMAXPROCS(CPU)
	c := make(chan int, CPU) // 缓冲区是可选的，但明显用上更好
	for i := 0; i < CPU; i++ {
		go v.DoSome(i*len(v)/CPU, (i+1)*len(v)/CPU, u, c)
	}

	// 排空信道。
	for i := 0; i < CPU; i++ {
		<-c // 等待任务完成
	}
	// 一切完成。
}

/*
目前 Go 运行时的实现默认并不会并行执行代码，它只为用户层代码提供单一的处理核心。
任意数量的 Go 程都可能在系统调用中被阻塞，而在任意时刻默认只有一个会执行用户层代码。
它应当变得更智能，而且它将来肯定会变得更智能。
但现在，若你希望 CPU 并行执行， 就必须告诉运行时你希望同时有多少 Go 程能执行代码。
有两种途径可意识形态，要么在运行你的工作时将 GOMAXPROCS 环境变量设为你要使用的核心数， 要么导入 runtime 包并调用 runtime.GOMAXPROCS(NCPU)。
runtime.NumCPU() 的值可能很有用，它会返回 当前机器的逻辑 CPU 核心数。
当然，随着调度算法和运行时的改进，将来会不再需要这种方法。

注意不要混淆并发和并行的概念:并发是用可独立执行的组件构造程序的方法，而并行则是为了效率在多 CPU 上平行地进行计算。
尽管 Go 的并发特性能够让某些问题更易构造成并行计算， 但 Go 仍然是种并发而非并行的语言，且 Go 的模型并不适合所有的并行问题。
关于其中区别的讨论，见 此博文(https://blog.golang.org/2013/01/concurrency-is-not-parallelism.html)。
*/
