package main

import "fmt"

// main Detecting Race Conditions With Go
// Go memory model 提到过: 表示写入单个 machine word 将是原子的，但 interface 内部是是两个 machine word 的值。另一个goroutine 可能在更改接口值时观察到它的内容。
// 在这个例子中，Ben 和 Jerry 内存结构布局是相同的，因此它们在某种意义上是兼容的。想象一下，如果他们有不同的内存布局会发生什么混乱？
// 如果是一个普通的指针、map、slice 可以安全的更新吗？
// 没有安全的 data race(safe data race)。您的程序要么没有 data race，要么其操作未定义。
// 原子性
// 可见行
func main() {
	var ben = &Ben{1, "Ben"}
	var jerry = &Jerry{"Jerry"}
	var maker IceCreamMaker = ben

	var loop0, loop1 func()

	loop0 = func() {
		maker = ben
		go loop1()
	}

	loop1 = func() {
		maker = jerry
		go loop0()
	}

	// 员工开始操作
	go loop0()
	for {
		maker.Hello()
	}
}

// IceCreamMaker interface 包含两部分 Type 和 Data, 它的负值操作不是一个原子性的
// 这里代码运行会存在nil情况（如果实现该接口的结构体内部的基础数据占用的字节数不同就会存在问题）而产生panic
type IceCreamMaker interface {
	// Hello greets a customer
	Hello()
}

type Ben struct {
	id   int
	name string
}

func (maker *Ben) Hello() {
	fmt.Printf("Ben says, \"Hello my name is %s\"\n", maker.name)
}

type Jerry struct {
	name string
}

func (maker *Jerry) Hello() {
	fmt.Printf("Jerry says, \"Hello my name is %s\"\n", maker.name)
}
