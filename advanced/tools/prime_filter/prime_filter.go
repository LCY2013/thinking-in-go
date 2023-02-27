package main

import (
	"context"
	"fmt"
	"sync"
)

// GenerateNatural 返回生成自然数序列的管道: 2, 3, 4, ...
func GenerateNatural(ctx context.Context, wg *sync.WaitGroup) chan int {
	ch := make(chan int)
	go func() {
		defer wg.Done()
		for i := 2; ; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}
	}()
	return ch
}

// PrimeFilter 管道过滤器: 删除能被素数整除的数
func PrimeFilter(ctx context.Context, in <-chan int, prime int, wg *sync.WaitGroup) chan int {
	out := make(chan int)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case i := <-in:
				if i%prime != 0 {
					select {
					case <-ctx.Done():
						return
					case out <- i:
					}
				}
			}

		}
	}()
	return out
}

func main() {
	wg := sync.WaitGroup{}
	// 通过 Context 控制后台 Goroutine 状态
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	ch := GenerateNatural(ctx, &wg) // 自然数序列: 2, 3, 4, ...
	for i := 0; i < 100; i++ {
		prime := <-ch // 新出现的素数
		fmt.Printf("%v: %v\n", i+1, prime)
		wg.Add(1)
		ch = PrimeFilter(ctx, ch, prime, &wg) // 基于新素数构造的过滤器
	}

	cancel()
	wg.Wait()
}
