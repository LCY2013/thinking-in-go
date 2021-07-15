package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

// 竞争状态 race condition
// 本例展示如何在程序里造成竞争状态，然后通过 atomic 包下的提供的原子操作来解决程序里的竞争状态

// go build -race race.go  竞争状态检测
// go run -race race.go  竞争状态检测

var (
	// counter 是所有 goroutine 都有增加其值的变量
	counter int64

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
func incCounter(id int64) {
	// 在函数退出是调用 Done 来通知 main 函数
	defer wg.Done()

	for count := 0; count < 2; count++ {

		// 安全读
		// atomic.LoadInt64(&counter)
		// 安全写
		// atomic.StoreInt64(&counter,id)

		// 安全的对 counter 加一
		atomic.AddInt64(&counter, id)

		// 当前 goroutine 从线程退出，并放回队列中
		runtime.Gosched()

	}
}
