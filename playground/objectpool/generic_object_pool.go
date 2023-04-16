package main

import (
	"sync"
	"sync/atomic"
)

type GoroutinesWithResultGenericPool[T, O any] struct {
	pool    sync.Pool
	fn      func(T) (O, error)
	options []GoroutinesWithResultOption[T, O]
}

func NewGenericPool[T, O any](fn func(T) (O, error),
	options ...GoroutinesWithResultOption[T, O]) *GoroutinesWithResultGenericPool[T, O] {
	gwrg := &GoroutinesWithResultGenericPool[T, O]{
		fn:      fn,
		options: options,
	}

	gwrg.pool = sync.Pool{
		New: gwrg.genericObject(),
	}

	return gwrg
}

func (p *GoroutinesWithResultGenericPool[T, O]) genericObject() func() any {
	return func() any {
		gwr := &GoroutinesWithResult[T, O]{
			fn:  p.fn,
			in:  make(chan T),
			out: make(chan O),
			n:   1,
		}

		for _, option := range p.options {
			option(gwr)
		}

		return gwr
	}
}

func (p *GoroutinesWithResultGenericPool[T, O]) Put(gwr *GoroutinesWithResult[T, O]) {
	// TODO
	gwr.n = 1
	gwr.in = make(chan T)
	gwr.out = make(chan O)
	gwr.hasErr = atomic.Bool{}
	gwr.closed = atomic.Bool{}

	p.pool.Put(gwr)
}

func (p *GoroutinesWithResultGenericPool[T, O]) Get() *GoroutinesWithResult[T, O] {
	gwr := p.pool.Get()
	gwrRes, ok := gwr.(*GoroutinesWithResult[T, O])
	if !ok {
		return nil
	}
	return gwrRes
}
