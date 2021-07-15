package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 竞争状态 race condition
// 本例展示如何在程序里造成竞争状态，然后通过 atomic 包下的提供的原子操作来解决程序里的竞争状态

// go build -race race.go  竞争状态检测
// go run -race race.go  竞争状态检测

var (
	// shutdown 是通知正在执行的所有 goroutine 停止工作的标识
	shutdown int64

	// wg 用来等待程序结束
	wg sync.WaitGroup
)

// main 程序的主入口
func main() {
	// 计数加2表示要等待两个 goroutine
	wg.Add(2)

	// 创建两个 goroutine 执行工作任务
	go doWork("A")
	go doWork("B")

	// 给定 goroutine 的执行时间
	time.Sleep(1 * time.Second)

	// 该工作停止了，安全设置 shutdown 标识符
	fmt.Println("Shutdown Now")
	atomic.StoreInt64(&shutdown, 1)

	// 等待 goroutine 程序结束
	wg.Wait()
}

// doWork 用来模拟执行工作的 goroutine，检测之前的 shutdown 标识符来决定是否提前终止
func doWork(name string) {
	// 在函数退出的时候执行 Done 用于通知 main 函数
	defer wg.Done()

	for {
		fmt.Printf("Doing %s Work\n", name)
		time.Sleep(200 * time.Millisecond)

		// 停止工作
		if atomic.LoadInt64(&shutdown) == 1 {
			fmt.Printf("Shutting %s Down\n", name)
			break
		}
	}
}
