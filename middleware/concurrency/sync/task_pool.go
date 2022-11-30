package sync

import (
	"fmt"
)

type TaskPool struct {
	ch chan func()
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

func (t *TaskPool) Run(fn func()) {
	t.ch <- fn
}
