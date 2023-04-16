package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type GoroutinesWithResult[T, O any] struct {
	n      int                // 协程数
	fn     func(T) (O, error) // 执行函数
	in     chan T
	out    chan O
	hasErr atomic.Bool
	closed atomic.Bool
	wg     sync.WaitGroup
}

type GoroutinesWithResultOption[T, O any] func(*GoroutinesWithResult[T, O])

func WithMaxGoroutines[T, O any](n int) GoroutinesWithResultOption[T, O] {
	return func(result *GoroutinesWithResult[T, O]) {
		result.n = n
	}
}

func NewGoroutinesWithResult[T, O any](
	fn func(T) (O, error),
	options ...GoroutinesWithResultOption[T, O],
) *GoroutinesWithResult[T, O] {
	gwr := &GoroutinesWithResult[T, O]{
		fn:  fn,
		in:  make(chan T),
		out: make(chan O),
		n:   1,
	}

	for _, option := range options {
		option(gwr)
	}

	return gwr
}

func (s *GoroutinesWithResult[T, O]) Data(data ...T) *GoroutinesWithResult[T, O] {
	go func(data ...T) {
		for idx, d := range data {
			s.in <- d
			if idx == len(data)-1 {
				s.close()
			}
		}
	}(data...)
	return s
}

func (s *GoroutinesWithResult[T, O]) do() {
	defer s.wg.Done()
	for val := range s.in {
		result, err := s.fn(val)
		if err != nil && !s.closed.Load() {
			s.hasErr.Store(true)
			continue
		}
		s.out <- result
	}
}

func (s *GoroutinesWithResult[T, O]) Run() *GoroutinesWithResult[T, O] {
	s.wg.Add(s.n)
	for i := 1; i <= s.n; i++ {
		go s.do()
	}

	go func() {
		s.wg.Wait()
		close(s.out)
	}()

	return s
}

func (s *GoroutinesWithResult[T, O]) Future() (chan<- O, bool) {
	return s.out, s.hasErr.Load()
}

func (s *GoroutinesWithResult[T, O]) close() {
	if !s.closed.Load() {
		close(s.in)
		s.closed.Store(true)
	}
}

func (s *GoroutinesWithResult[T, O]) Finish() ([]O, bool) {
	val := make([]O, 0)
	for v := range s.out {
		val = append(val, v)
	}
	return val, s.hasErr.Load()
}

func main() {
	const numCount = 1001

	var data = make([]int, 0, numCount)
	for i := 0; i < numCount; i++ {
		data = append(data, i)
	}

	p := NewGenericPool[int, string](func(s int) (string, error) {
		var err error
		if s == 10 {
			err = fmt.Errorf("xxx")
		}
		return fmt.Sprintf("%d+++", s), err
	}, WithMaxGoroutines[int, string](10))

	s := p.Get()

	result, hasErr := s.Data(data...).Run().Finish()

	if hasErr {
		fmt.Println(len(result), result)
	}

	p.Put(s)

	s = p.Get()

	result, hasErr = s.Data(data...).Run().Finish()

	if hasErr {
		fmt.Println(len(result), result)
	}
}
