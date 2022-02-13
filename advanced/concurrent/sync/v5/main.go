package main

import "fmt"

// main Detecting Race Conditions With Go
// 修复v4版本
func main() {
	var ben = &Ben{"Ben"}
	var jerry = &Jerry{field2: 2}
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
// 这里代码运行会存在nil情况而产生panic
type IceCreamMaker interface {
	// Hello greets a customer
	Hello()
}

type Ben struct {
	//id   int
	name string
}

func (maker *Ben) Hello() {
	fmt.Printf("Ben says, \"Hello my name is %s\"\n", maker.name)
}

type Jerry struct {
	field2 int
	//name string
	field1 *[5]byte
}

func (maker *Jerry) Hello() {
	fmt.Printf("Jerry says, \"Hello my name is %s\"\n", maker.field1)
}
