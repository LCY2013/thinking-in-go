package main

import (
	"fmt"
	"sync"
)

type goRes[O any] struct {
	out chan O
}

// StableGoroutinesWithResult 固定协程数并发执行器 并接受处理结果
type StableGoroutinesWithResult[T, O any] struct {
	n   int       // 协程数
	fn  func(T) O // 执行函数
	res goRes[O]
	wg  sync.WaitGroup
	in  chan T
}

func NewStableGoroutinesWithResult[T, O any](n int, fn func(T) O, in chan T) *StableGoroutinesWithResult[T, O] {
	return &StableGoroutinesWithResult[T, O]{
		n:  n,
		fn: fn,
		res: goRes[O]{
			// 增加缓冲区大小
			out: make(chan O, n),
		},
		wg: sync.WaitGroup{},
		in: in,
	}
}

func (s *StableGoroutinesWithResult[T, O]) do() {
	defer s.wg.Done()
	for val := range s.in {
		s.res.out <- s.fn(val)
	}
}

func (s *StableGoroutinesWithResult[T, O]) Run() {
	s.wg.Add(s.n)
	for i := 1; i <= s.n; i++ {
		go func() { // 匿名函数
			s.do()
		}()
	}

	go func() {
		s.wg.Wait()
		close(s.res.out)
	}()

	return
}

func (s *StableGoroutinesWithResult[T, O]) End() []O {
	val := make([]O, 0)
	for v := range s.res.out {
		val = append(val, v)
	}
	return val
}

func (s *StableGoroutinesWithResult[T, O]) SendData(data ...T) {
	for _, d := range data {
		s.in <- d
	}
}

func (s *StableGoroutinesWithResult[T, O]) AsyncSendData(data ...T) {
	for _, d := range data {
		s.in <- d
	}
}

func main() {
	in := make(chan int)
	const numWorkers = 2

	// 使用new函数初始化泛型类型的结构体
	s := NewStableGoroutinesWithResult[int, string](numWorkers, func(s int) string {
		return fmt.Sprintf("%d+++", s)
	}, in)

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s.SendData(data...)

	s.Run()

	go func() {
		result := s.End()

		fmt.Println(result)
	}()

	close(in)

}
