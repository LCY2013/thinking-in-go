package main

import (
	"context"
	"fmt"
)

func main() {
	gen := func(ctx context.Context) <-chan int {
		ch := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Done")
					return // 防止goroutine 泄漏
				case ch <- n:
					n++
				}
			}
		}()
		return ch
	}

	// 定义新的context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for v := range gen(ctx) {
		fmt.Println(v)
		if v == 5 {
			break
		}
	}
}
