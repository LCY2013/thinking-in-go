package main

import (
	"time"

	"github.com/lcy2013/workerpool"
)

// main 这个示例程序创建了一个 capacity 为 5 的 workerpool 实例，并连续向这个 workerpool 提交了 10 个 task，每个 task 的逻辑很简单，只是 Sleep 3 秒后就退出。
// main 函数在提交完任务后，调用 workerpool 的 Free 方法销毁 pool，pool 会等待所有 worker 执行完 task 后再退出。
func main() {
	p := workerpool.New(5)
	for i := 0; i < 10; i++ {
		err := p.Schedule(func() {
			time.Sleep(time.Second * 3)
		})
		if err != nil {
			println("task: ", i, "err: ", err)
		}
	}
	p.Free()

	time.Sleep(time.Second * 6)
}
