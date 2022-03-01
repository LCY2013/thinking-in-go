package main

import (
	"fmt"
	"github.com/lcy2013/workerpool"
	"time"
)

// main 不过，由于 Goroutine 调度的不确定性，这个结果仅仅是很多种结果的一种。仅仅 002 这个 worker 收到了 task，其余的 worker 都因为 worker 尚未创建完毕，而返回了错误，而不是像 demo1 那样阻塞在 Schedule 调用上。
func main() {
	pool := workerpool.New(5, workerpool.WithBlock(false), workerpool.WithPreAllocWorkers(false))
	time.Sleep(time.Second * 2)
	for i := 0; i < 10; i++ {
		err := pool.Schedule(func() {
			time.Sleep(time.Second * 3)
		})
		if err != nil {
			fmt.Printf("task[%d]: error: %s\n", i, err)
		}
	}

	pool.Free()
}
