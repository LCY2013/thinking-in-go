package main

import (
	"fmt"
	"runtime"
	"sync"
)

// 互斥锁 mutex , 互斥(mutual exclusion)
// 保证一段代码属于一个临界区，在这临界区内的代码是保证同一时间只能有一个 goroutine 去执行

// 本例使用 互斥锁 保证一段需要同步的代码资源的同步

var (
	// counter 所有的 goroutine 都要增加其值的变量
	counter int

	// wg 用来等待程序的结束
	wg sync.WaitGroup

	// mutex 用来定义一段代码临界区
	mutex sync.Mutex
)

// main 程序的主入口
func main() {
	// 计数器增加 2 ， 表示有两个需要等待的 goroutine
	wg.Add(2)

	// 创建两个 goroutine 执行任务
	go incCounter(1)
	go incCounter(2)

	// 等待 goroutine 结束
	wg.Wait()
	fmt.Printf("Final counter: %d\n", counter)
}

// incCounter 使用互斥锁来同步并保证安全访问，增加包里 counter 变量的值
func incCounter(id int) {
	// 在函数退出时调用 Done 来通知 main 函数工作以及完成
	defer wg.Done()

	for count := 0; count < 2; count++ {
		// 同时只允许一个 goroutine 进入这个临界区
		mutex.Lock()
		{
			// 捕获counter的值
			value := counter

			// 当前 goroutine 从线程退出，并放回到队列中
			runtime.Gosched()

			// 增加本地 value 变量的值
			value++

			// 将该值保存回 counter
			counter = value
		}
		// 释放锁，允许其他的 goroutine 进入临界区
		mutex.Unlock()
	}
}
