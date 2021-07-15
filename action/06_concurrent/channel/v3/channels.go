package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 有缓冲通道（buffered channel）
// 使用有缓冲的通道和固定数目的 goroutine 执行一堆工作

const (
	numberGoroutines = 4  // 要使用的 goroutine 的数量
	taskLoad         = 10 // 需要处理的工作数量
)

// wg 用来等待程序的结束
var wg sync.WaitGroup

// init 初始化包，Go语言运行时会优先执行在其他代码之前
func init() {
	// 初始化随机种子
	rand.Seed(time.Now().Unix())
}

// main go程序的主入口
func main() {
	// 创建一个有缓冲的通道来管理工作
	tasks := make(chan string, taskLoad)

	// 启动 goroutine 来处理工作
	wg.Add(numberGoroutines)
	for gr := 1; gr <= numberGoroutines; gr++ {
		go worker(tasks, gr)
	}

	// 增加一组要完成的工作
	for post := 1; post <= taskLoad; post++ {
		tasks <- fmt.Sprintf("Task : %d", post)
	}

	// 当所有工作处理完成时关闭通道，以便所有的 goroutine 安全退出
	close(tasks)

	// 等待所有工作完成
	wg.Wait()
}

// worker 作为 goroutine 启动来处理，从有缓冲的通道传入的工作
func worker(tasks chan string, worker int) {
	// 通知 main 函数任务已经完成
	defer wg.Done()

	for {
		// 等待分配工作
		task, ok := <-tasks
		if !ok {
			// 意味着通道不为空，且已经关闭
			fmt.Printf("Worker: %d : Shutting Down\n", worker)
			return
		}

		// 显示开始工作
		fmt.Printf("Worker : %d : Started %s\n", worker, task)

		// 随机一段时间来模拟工作
		sleep := rand.Int63n(100)
		time.Sleep(time.Duration(sleep) * time.Millisecond)

		// 显示工作已经完成
		fmt.Printf("Worker : %d : Completed %s\n", worker, task)
	}
}
