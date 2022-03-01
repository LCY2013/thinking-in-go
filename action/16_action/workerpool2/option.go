package workerpool

// Option 功能选项（functional option）设计
// 为 workerpool 添加两个功能选项：Schedule 调用是否阻塞，以及是否预创建所有的 worker。
type Option func(*Pool)

// WithBlock Schedule 调用是否阻塞
func WithBlock(block bool) Option {
	return func(pool *Pool) {
		pool.block = block
	}
}

// WithPreAllocWorkers 是否预创建所有的 worker
func WithPreAllocWorkers(preAlloc bool) Option {
	return func(pool *Pool) {
		pool.preAlloc = preAlloc
	}
}
