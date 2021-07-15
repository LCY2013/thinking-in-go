package main

import (
	"fmt"
	"runtime"
	"sync"
)

// 展示 goroutine 调度器如何在单线程上切分时间片

// wg 声明等待组
var wg sync.WaitGroup

// main 程序主入口
func main() {
	// 分配一个逻辑处理器给调度器使用
	runtime.GOMAXPROCS(1)

	// 给每个可用的核心分配一个逻辑处理器
	// runtime.GOMAXPROCS(runtime.NumCPU())

	// 计数器加 2 表示要等待两个 goroutine
	wg.Add(2)

	// 打印前置
	fmt.Println("Create Goroutines")

	// 创建两个 goroutine , 用于模拟单线程切分时间片的过程
	go printPrime("A")
	go printPrime("B")

	// 等待 goroutine 结束
	fmt.Println("Waiting To Finish")

	// 等待组完成
	wg.Wait()

	// 打印完成
	fmt.Println("Terminating Goroutines")
}

// printPrime 显示 5000 以内的素数
func printPrime(prefix string) {
	// 在函数退出时调用 Done 来通知 main 函数
	defer wg.Done()

next:
	for outer := 2; outer < 5000; outer++ {
		for inner := 2; inner < outer; inner++ {
			if outer%inner == 0 {
				continue next
			}
		}
		fmt.Printf("%s:%d\n", prefix, outer)
	}

	fmt.Println("Completed ", prefix)
}
