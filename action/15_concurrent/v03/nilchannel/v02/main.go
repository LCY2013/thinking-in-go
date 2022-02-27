package main

import (
	"fmt"
	"time"
)

/*
这里，改进后的示例程序的最关键的一个变化，就是在判断 ch1 或 ch2 被关闭后，显式地将 ch1 或 ch2 置为 nil。
而前面已经知道了，对一个 nil channel 执行获取操作，这个操作将阻塞。
于是，这里已经被置为 nil 的 c1 或 c2 的分支，将再也不会被 select 选中执行。
*/
func main() {
	ch1, ch2 := make(chan int), make(chan int)
	go func() {
		time.Sleep(time.Second * 1)
		ch1 <- 5
		close(ch1)
	}()
	go func() {
		time.Sleep(time.Second * 2)
		ch2 <- 7
		close(ch2)
	}()
	for {
		select {
		case x, ok := <-ch1:
			if !ok {
				ch1 = nil
			} else {
				fmt.Println(x)
			}
		case x, ok := <-ch2:
			if !ok {
				ch2 = nil
			} else {
				fmt.Println(x)
			}
		}
		if ch1 == nil && ch2 == nil {
			break
		}
	}
	fmt.Println("program end")
}
