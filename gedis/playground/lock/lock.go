package main

import (
	"fmt"
	"sync"
)

type Person struct {
	mu     sync.Mutex
	level  int
	salary int
}

// main 监测bug或者可疑的结构 go vet lock.go
// ➜  lock git:(master) ✗ go vet lock.go
// # command-line-arguments
// ./lock.go:22:8: assignment copies lock value to p1: command-line-arguments.Person contains sync.Mutex
// ./lock.go:24:14: call of fmt.Println copies lock value: command-line-arguments.Person contains sync.Mutex
// go build -race lock.go  race监测锁竞争
// 死锁监测 go-deadlock项目 https://github.com/sasha-s/go-deadlock
func main() {
	p := Person{
		level:  1,
		salary: 1000000,
	}

	p1 := p

	fmt.Println(p1)
}
