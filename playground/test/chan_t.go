package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("second")
			case <-ch:
				break
			}
		}
	}()

	time.Sleep(5 * time.Second)
	close(ch)
}
