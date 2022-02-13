package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var Counter int = 0

// main go build -race
func main() {
	for routine := 1; routine <= 2; routine++ {
		wg.Add(1)
		go Routine(routine)
	}
	wg.Wait()
	fmt.Println(Counter)
}

func Routine(id int) {
	for count := 0; count < 2; count++ {
		value := Counter
		time.Sleep(time.Nanosecond)
		value++
		Counter = value
	}

	wg.Done()
}
