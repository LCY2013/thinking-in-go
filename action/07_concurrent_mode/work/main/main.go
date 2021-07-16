package main

import (
	"fufeng.org/concurrent_mode/work"
	"log"
	"sync"
	"time"
)

// 展示使用 work 包
// 创建一个 goroutine 池并完成工作

// names 提供一组用来显示的名字
var names = []string{
	"Steve",
	"Bob",
	"Mary",
	"Therese",
	"Jason",
}

// namePrinter 使用特定方式打印名字
type namePrinter struct {
	name string
}

// Task 实现 Worker 接口
func (name *namePrinter) Task() {
	log.Println(name.name)
	time.Sleep(time.Second)
}

// main 程序的主入口
func main() {
	// 使用两个 goroutine 来创建工作池
	p := work.New(2)

	var wg sync.WaitGroup
	wg.Add(10 * len(names))

	for i := 0; i < 10; i++ {
		// 迭代 names 切片
		for _, name := range names {
			// 创建一个 namePrinter 并指定名字
			np := namePrinter{
				name: name,
			}
			go func() {
				// 将任务提交执行，当 Run 返回时就是任务已经处理完成了
				p.Run(&np)
				wg.Done()
			}()
		}
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 让工作池停止工作，等待所有的工作完成
	p.Shutdown()
}
