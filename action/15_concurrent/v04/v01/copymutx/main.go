package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	i := 0
	var mu sync.Mutex // 负责对i的同步访问

	wg.Add(1)
	// g1
	go func(mu1 sync.Mutex) {
		mu1.Lock()
		i = 2
		time.Sleep(3 * time.Second)
		fmt.Printf("g1: i = %d\n", i)
		mu1.Unlock()
		wg.Done()
	}(mu)

	time.Sleep(time.Second)

	mu.Lock()
	i = 1
	fmt.Printf("g0: i = %d\n", i)
	mu.Unlock()

	wg.Wait()
}
