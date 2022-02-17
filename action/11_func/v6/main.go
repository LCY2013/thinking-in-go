package main

import "fmt"

/*
健壮性的“三不要”原则
原则一：不要相信任何外部输入的参数。

原则二：不要忽略任何一个错误。

原则三：不要假定异常不会发生。
*/
func foo() {
	println("call foo")
	//bar()
	barUp()
	println("exit foo")
}

func bar() {
	println("call bar")
	panic("panic occurs in bar")
	zoo()
	println("exit bar")
}

func barUp() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("recover the panic:", e)
		}
	}()
	println("call bar")
	panic("panic occurs in bar")
	zoo()
	println("exit bar")
}

func zoo() {
	println("call zoo")
	println("exit zoo")
}

func main() {
	println("call main")
	foo()
	println("exit main")
}
