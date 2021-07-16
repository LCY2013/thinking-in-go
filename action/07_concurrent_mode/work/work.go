package work

import "sync"

// work 包管理了一个 goroutine 池来完成工作

// Worker 必须满足接口类型，才能使用工作池
type Worker interface {
	Task()
}

// Pool 提供一个 goroutine 池，这个池完成任何已提交的 Worker 工作
type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

// New 创建一个新的工作池
func New(maxGoroutine int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}

	p.wg.Add(maxGoroutine)
	for i := 0; i < maxGoroutine; i++ {
		go func() {
			for w := range p.work {
				w.Task()
			}
			p.wg.Done()
		}()
	}

	return &p
}

// Run 提交工作到工作池
func (p *Pool) Run(w Worker) {
	p.work <- w
}

// Shutdown 等待所有 goroutine 完成任务
func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
