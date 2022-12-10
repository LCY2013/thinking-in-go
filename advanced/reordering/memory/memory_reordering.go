package main

import (
	"fmt"
	"sync"
)

/*
https://cch123.github.io/ooo/

从 Memory Reordering 说起

下面这段代码会有怎样的输出
var x, y int

	go func() {
	    x = 1 // A1
	    fmt.Print("y:", y, " ") // A2
	}()

	go func() {
	    y = 1                   // B1
	    fmt.Print("x:", x, " ") // B2
	}()

显而易见的几种结果:
y:0 x:1
x:0 y:1
x:1 y:1
y:1 x:1
令人意外的结果
x:0 y:0
y:0 x:0

这种令人意外的结果被称为内存重排： Memory Reordering

什么是内存重排

内存重排指的是内存的读/写指令重排。  — Xargin

软件或硬件系统可以根据其对代码的分析结果,一定程度上打乱代码的执行顺序，以达到其不可告人的目的。

为什么会发生内存重排？

1、编译器重排

2、CPU重排
*/
func main() {
	var wg sync.WaitGroup
	for {
		var x, y int
		wg.Add(2)
		go func() {
			defer wg.Done()
			x = 1
			fmt.Print("y:", y, " ")
		}()

		go func() {
			defer wg.Done()
			y = 1
			fmt.Print("x:", x, " ")
		}()
		fmt.Println()
		wg.Wait()
	}
}
