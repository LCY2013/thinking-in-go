package sync

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type TaskPool struct {
	ch       chan func()
	closeNow chan int
	closed   uint32
	lock     sync.Mutex
}

func NewTaskPool(runSize, queueSize int) *TaskPool {
	taskPool := &TaskPool{
		ch:       make(chan func(), queueSize),
		closeNow: make(chan int, 0),
	}

	for i := 0; i < runSize; i++ {
		go func(i int) {
			for {
				select {
				case task, ok := <-taskPool.ch:
					if !ok {
						fmt.Printf("i = [%d] closed.\n", i)
						return
					}
					fmt.Printf("i = [%d] start.\n", i)
					task()
					fmt.Printf("i = [%d] end.\n", i)
					if atomic.LoadUint32(&taskPool.closed) != 0 && len(taskPool.ch) == 0 {
						// 探测是否符合预期
						taskPool.closeNow <- len(taskPool.ch)
					}
				}
			}
		}(i)
	}

	return taskPool
}

func (t *TaskPool) Run(fn func()) error {
	if t == nil {
		return fmt.Errorf("task pool, need to initialize")
	}
	if atomic.LoadUint32(&t.closed) != 0 {
		return fmt.Errorf("task pool, already closed")
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	t.ch <- fn
	return nil
}

func (t *TaskPool) Stop(ctx context.Context) {
	if atomic.LoadUint32(&t.closed) != 0 {
		return
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	atomic.AddUint32(&t.closed, 1)
	if len(t.ch) == 0 {
		close(t.ch)
		return
	}
	var breakFor bool
	timeOut := time.After(time.Second * 10)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("上下文关闭")
			breakFor = true
		case <-timeOut:
			fmt.Println("默认超时关闭")
			breakFor = true
		case <-t.closeNow:
			fmt.Println("已有任务执行完成")
			breakFor = true
		}
		if breakFor {
			break
		}
	}
	close(t.ch)
}
