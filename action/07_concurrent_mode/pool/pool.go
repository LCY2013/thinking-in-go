package pool

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
)

// pool包管理用户定义的一组资源

// Pool 管理一组可以安全的在多个 goroutine 之间共享的资源
// 被管理的资源必须实现 io.Closer 接口
type Pool struct {
	m         sync.Mutex
	resources chan io.Closer
	factory   func() (io.Closer, error)
	closed    bool
}

// ErrPoolClosed 表示请求（Acquire） 了一个以及关闭的池
var ErrPoolClosed = errors.New("Pool has been closed.")

// New 用来创建一个用来管理资源的池
// 这个池需要一个可以用来分配资源的函数，以及可以规定池的大小
func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size < 0 {
		return nil, errors.New("Pool size value too small.")
	}

	//return &Pool{
	//	factory:   fn,
	//	resources: make(chan io.Closer, size),
	//}, nil

	p := &Pool{
		factory:   fn,
		resources: make(chan io.Closer, size),
	}

	for index := uint(0); index < size; index++ {
		r, err := fn()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		p.resources <- r
	}

	return p, nil
}

// Acquires 从池中获取一个资源
func (p *Pool) Acquires() (io.Closer, error) {
	select {
	// 检测是否有可用的资源
	case r, ok := <-p.resources:
		log.Println("Acquires", " Shared Resource")
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
		// 没有空的资源新建一个资源
		//default:
		//	log.Println("Acquires", " New Resource")
		//	return p.factory()
	}
}

// Release 将一个使用后的资源放回到池中
func (p *Pool) Release(r io.Closer) {
	// 保证本操作和Close操作的安全
	p.m.Lock()
	defer p.m.Unlock()

	// 如果池已经关闭，就关闭资源
	if p.closed {
		r.Close()
		return
	}

	select {
	// 试图将这个资源放入到池中
	case p.resources <- r:
		log.Println("Release: ", "In Queue")

	// 如果队列已经满了，就关闭这个资源
	default:
		log.Println("Release: ", "Closing")
		r.Close()
	}
}

// Close 会让资源池停止工作，并关闭所有现有的资源
func (p *Pool) Close() {
	// 保证本操作与Release操作的安全性
	p.m.Lock()
	defer p.m.Unlock()

	// 如果 pool 已经被关闭，什么也不做
	if p.closed {
		return
	}

	// 将池关闭
	p.closed = true

	// 在清空通道里的资源之前，将通道关闭
	// 如果不这么做会发生死锁
	close(p.resources)

	// 关闭资源
	for r := range p.resources {
		r.Close()
	}
}
