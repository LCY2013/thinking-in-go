package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 展示利用无缓冲通道来模拟 2个 goroutine 间的网球比赛

// wg 用来等待程序的结束
var wg sync.WaitGroup

// init 初始化应用程序
func init() {
	rand.Seed(time.Now().UnixNano())
}

// main 是所有 go 程序的入口
func main() {
	// 创建一个无缓冲的通道
	court := make(chan int)

	// 计数加2，表示要等待 2 个goroutine
	wg.Add(2)

	// 启动两个选手
	go player("Nadal", court)
	go player("Djokovic", court)

	// 发球
	court <- 1

	// 等待游戏结束
	wg.Wait()
}

// player 模拟一个选手在打网球
func player(name string, court chan int) {
	// 在函数退出时调用 Done 来通知 main 函数工作以及完成
	defer wg.Done()

	for {
		// 等待求被击打过来
		ball, ok := <-court
		if !ok {
			// 如果通道关闭，我们就赢了
			fmt.Printf("Player %s Won\n", name)
			return
		}

		// 选择一个随机数，然后用这个随机数判断是否需要丢球
		n := rand.Intn(100)
		if n%13 == 0 {
			fmt.Printf("Player %s Missed\n", name)

			// 关闭通道表示该选手已经输了
			close(court)
			return
		}

		// 显示击球数，并将击球数加一
		fmt.Printf("Player %s Hit %d\n", name, ball)
		ball++

		// 将球打向对方
		court <- ball
	}
}
