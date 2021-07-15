package main

import (
	"fmt"
	"runtime"
	"sync"
)

// 竞争状态 race condition
// 本例展示如何在程序里造成竞争状态，实际上不希望出现这种情况

// go build -race race.go  竞争状态检测
// go run -race race.go  竞争状态检测

var (
	// counter 是所有 goroutine 都有增加其值的变量
	counter int

	// wg 用来等待程序结束
	wg sync.WaitGroup
)

// main 程序的主入口
func main() {
	// 计数加2表示要等待两个 goroutine
	wg.Add(2)

	// 创建两个用于新增 counter 值的 goroutine
	go incCounter(1)
	go incCounter(2)

	// 等待 goroutine 结束
	wg.Wait()
	fmt.Println("Final Counter: ", counter)
}

// incCounter 增加包里 counter 的值
func incCounter(id int) {
	// 在函数退出是调用 Done 来通知 main 函数
	defer wg.Done()

	for count := 0; count < 2; count++ {
		// 捕获 counter 的值
		value := counter

		// 当前 goroutine 从线程退出，并放回队列中
		runtime.Gosched()

		// 增加本地 value 变量的值
		value += id

		// 将该值放回 counter
		counter = value
	}
}
