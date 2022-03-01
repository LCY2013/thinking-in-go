package workerpool

import (
	"errors"
	"fmt"
	"sync"
)

/*
workerpool 包对外主要提供三个 API，它们分别是：
1、workerpool.New：用于创建一个 pool 类型实例，并将 pool 池的 worker 管理机制运行起来；
2、workerpool.Free：用于销毁一个 pool 池，停掉所有 pool 池中的 worker；
3、Pool.Schedule：这是 Pool 类型的一个导出方法，workerpool 包的用户通过该方法向 pool 池提交待执行的任务（Task）。
*/

// ErrWorkerPoolFreed workerpool已终止运行
var ErrWorkerPoolFreed = errors.New("workerpool freed")

// ErrNoWorkerAvailInPool 协程池中无可用的Worker
var ErrNoWorkerAvailInPool = errors.New("no worker avail in pool")

const (
	defaultCapacity = 10
	maxCapacity     = 1000
)

// Task 是一个对用户提交的请求的抽象，它的本质就是一个函数类型
// 这样，用户通过 Schedule 方法实际上提交的是一个函数类型的实例。
type Task func()

// Pool 协程池定义结构体
type Pool struct {
	capacity int // workpool 大小

	active chan struct{} // active channel
	tasks  chan Task     // task channel

	wg   sync.WaitGroup // 用于在pool销毁时等待所有worker退出
	quit chan struct{}  // 用于通知各个worker退出的信号channel

	preAlloc bool //是否在创建pool的时候就预创建workers，默认值为：false

	// block 当pool满的情况下，新的Schedule调用是否阻塞当前goroutine。默认值：true
	// 如果block = false，则Schedule返回ErrNoWorkerAvailInPool
	block bool
}

// run 执行协程池信息
// Pool 类型实例变量 p 完成初始化后，创建了一个新的 Goroutine，用于对 workerpool 进行管理，这个 Goroutine 执行的是 Pool 类型的 run 方法。
// run 方法内是一个无限循环，循环体中使用 select 监视 Pool 类型实例的两个 channel：quit 和 active。
// 这种在 for 中使用 select 监视多个 channel 的实现，在 Go 代码中十分常见，是一种惯用法。
// 当接收到来自 quit channel 的退出“信号”时，这个 Goroutine 就会结束运行。
// 而当 active channel 可写时，run 方法就会创建一个新的 worker Goroutine。
// 此外，为了方便在程序中区分各个 worker 输出的日志，这里将一个从 1 开始的变量 idx 作为 worker 的编号，并把它以参数的形式传给创建 worker 的方法。
//
// 新版 run 方法在 preAlloc=false 时，会根据 tasks channel 的情况在适合的时候创建 worker（第 61 行~ 第 75 行)，直到 active channel 写满，才会进入到和第一版代码一样的调度逻辑中（第 77 行~ 第 86 行）。
func (p *Pool) run() {
	// 通过获取active通道内的数据长度作为goroutine的起始下标
	idx := len(p.active)

	if !p.preAlloc {
	loop:
		for t := range p.tasks {
			p.returnTask(t)
			select {
			case <-p.quit:
				return
			case p.active <- struct{}{}:
				idx++
				p.newWorker(idx)
			default:
				break loop
			}
		}
	}

	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}:
			// create new worker
			idx++
			p.newWorker(idx)
		}
	}
}

// newWorker 创建新的 worker goroutine
//
// 在创建一个新的 worker goroutine 之前，newWorker 方法会先调用 p.wg.Add 方法将 WaitGroup 的等待计数加一。
// 由于每个 worker 运行于一个独立的 Goroutine 中，newWorker 方法通过 go 关键字创建了一个新的 Goroutine 作为 worker。
//
// 新 worker 的核心，依然是一个基于 for-select 模式的循环语句，在循环体中，新 worker 通过 select 监视 quit 和 tasks 两个 channel。
// 和前面的 run 方法一样，当接收到来自 quit channel 的退出“信号”时，这个 worker 就会结束运行。
// tasks channel 中放置的是用户通过 Schedule 方法提交的请求，新 worker 会从这个 channel 中获取最新的 Task 并运行这个 Task。
func (p *Pool) newWorker(idx int) {
	p.wg.Add(1)
	go func() {
		// 在新 worker 中，为了防止用户提交的 task 抛出 panic，进而导致整个 workerpool 受到影响，
		// 在 worker 代码的开始处，使用了 defer+recover 对 panic 进行捕捉，捕捉后 worker 也是要退出的，
		// 于是我们还通过<-p.active更新了 worker 计数器。
		// 并且一旦 worker goroutine 退出，p.wg.Done 也需要被调用，这样可以减少 WaitGroup 的 Goroutine 等待数量。
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exit\n", idx, err)
				<-p.active
			}
			p.wg.Done()
		}()

		fmt.Printf("worker[%03d]: start\n", idx)

		for {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d]: exit\n", idx)
				<-p.active
				return
			case task := <-p.tasks:
				fmt.Printf("worker[%03d]: receive a task\n", idx)
				task()
			}
		}
	}()
}

// Schedule workerpool 提供给用户提交请求的导出方法 Schedule
// Schedule 方法的核心逻辑，是将传入的 Task 实例发送到 workerpool 的 tasks channel 中。
// 但考虑到现在 workerpool 已经被销毁的状态，这里通过一个 select，检视 quit channel 是否有“信号”可读，如果有，就返回一个哨兵错误 ErrWorkerPoolFreed。
// 如果没有，一旦 p.tasks 可写，提交的 Task 就会被写入 tasks channel，以供 pool 中的 worker 处理。
//
// 这里要注意的是，这里的 Pool 结构体中的 tasks 是一个无缓冲的 channel，
// 如果 pool 中 worker 数量已达上限，而且 worker 都在处理 task 的状态，那么 Schedule 方法就会阻塞，直到有 worker 变为 idle 状态来读取 tasks channel，schedule 的调用阻塞才会解除。
//
// 提供给用户的 Schedule 函数也因 WithBlock 选项，有了一些变化。
// Schedule 在 tasks chanel 无法写入的情况下，进入 default 分支。在 default 分支中，Schedule 根据 block 字段的值，决定究竟是继续阻塞在 tasks channel 上，还是返回 ErrNoIdleWorkerInPool 错误。
func (p *Pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- t:
		return nil
	default:
		if p.block {
			p.tasks <- t
			return nil
		}
		return ErrNoWorkerAvailInPool
	}
}

func (p *Pool) Free() {
	close(p.quit)
}

func (p *Pool) returnTask(t Task) {
	go func() {
		p.tasks <- t
	}()
}

// New workerpool.New 是如何创建一个 pool 实例的
// New 函数接受一个参数 capacity 用于指定 workerpool 池的容量，
// 这个参数用于控制 workerpool 最多只能有 capacity 个 worker，共同处理用户提交的任务请求。
// 函数开始处有一个对 capacity 参数的“防御性”校验，当用户传入不合理的值时，函数 New 会将它纠正为合理的值。
func New(capacity int, opts ...Option) *Pool {
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}

	p := &Pool{
		capacity: capacity,
		tasks:    make(chan Task),
		quit:     make(chan struct{}),
		active:   make(chan struct{}, capacity),
	}

	// 功能选项初始化
	for _, opt := range opts {
		opt(p)
	}

	fmt.Printf("workerpool start\n")

	// 如果需要初始化worker
	if p.preAlloc {
		// create all goroutines and send into works channel
		for i := 0; i < p.capacity; i++ {
			p.newWorker(i)
			p.active <- struct{}{}
		}
	}

	go p.run()

	return p
}
