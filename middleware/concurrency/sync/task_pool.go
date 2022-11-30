package sync

import (
	"fmt"
	"sync/atomic"
)

type TaskPool struct {
	ch     chan func()
	closed uint32
}

func NewTaskPool(runSize, queueSize int) *TaskPool {
	taskPool := &TaskPool{
		ch: make(chan func(), queueSize),
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
	t.ch <- fn
	return nil
}

func (t *TaskPool) Stop() {
	if atomic.LoadUint32(&t.closed) != 0 {
		return
	}
	close(t.ch)
	atomic.AddUint32(&t.closed, 1)
}
