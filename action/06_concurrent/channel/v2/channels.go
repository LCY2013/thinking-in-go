package main

import (
	"fmt"
	"sync"
	"time"
)

// 使用无缓冲的通道模拟 4个 goroutine 间的接力赛

// wg 用来等待程序的结束
var wg sync.WaitGroup

// main go语言程序的入口
func main() {
	// 创建一个无缓冲的通道
	baton := make(chan int)

	// 为最后一个位跑步者将计数加一
	wg.Add(1)

	// 第一位跑步者持有接力棒
	go Runner(baton)

	// 开始比赛
	baton <- 1

	// 等待比赛结束
	wg.Wait()
}

// Runner 模拟跑步比赛中的一位跑步者
func Runner(baton chan int) {
	var newRunner int

	// 等待接力棒
	runner := <-baton

	// 开始绕着跑道跑步
	fmt.Printf("Runner %d Running With Baton\n", runner)

	// 创建下一位跑步者
	if runner != 4 {
		newRunner = runner + 1
		fmt.Printf("Runner %d To The Line\n", newRunner)
		go Runner(baton)
	}

	// 围绕跑到跑
	time.Sleep(100 * time.Millisecond)

	// 比赛结束
	if runner == 4 {
		fmt.Printf("Runner %d finished,Race Over\n", runner)
		wg.Done()
		return
	}

	// 将接力棒交给下一位跑步者
	fmt.Printf("Runner %d Exchange With Runner %d\n", runner, newRunner)

	baton <- newRunner
}
