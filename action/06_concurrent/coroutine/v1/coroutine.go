// 展示构建 goroutine 和以及调度器行为
package main

import (
	"fmt"
	"runtime"
	"sync"
)

// main 程序的入口
func main() {
	// 分配一个逻辑处理器给调度器使用
	// runtime.GOMAXPROCS(1)

	// 分配两个逻辑处理器给调度器使用
	runtime.GOMAXPROCS(2)

	// wg 等待组，用来等待程序的完成，计数器加2，表示等待两个 goroutine 的完成
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Start Goroutines")

	// 声明一个匿名函数，创建一个 goroutine
	go func() {
		// 在函数退出时调用 Done 来通知 main 函数工作已完成
		defer wg.Done()

		// 显示字母表3次
		for counter := 0; counter < 3; counter++ {
			for char := 'a'; char < 'a'+26; char++ {
				fmt.Printf("%c ", char)
			}
		}
	}()

	// 声明一个匿名函数，创建一个 goroutine
	go func() {
		// 在函数退出时调用 Done 来通知 main 函数工作已完成
		defer wg.Done()

		// 显示字母表3次
		for counter := 0; counter < 3; counter++ {
			for char := 'A'; char < 'A'+26; char++ {
				fmt.Printf("%c ", char)
			}
		}
	}()

	// 等待 goroutine 结束
	fmt.Println("Waiting To Finish")
	wg.Wait()

	fmt.Println("\n Terminating Program")
}
