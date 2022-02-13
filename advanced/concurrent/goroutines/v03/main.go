package main

import "fmt"

// goroutine 泄漏

func leak() {
	ch := make(chan int)
	go func() {
		chanVal := <-ch
		fmt.Printf("We received a value: %d\n", chanVal)
	}()
}
