package main

import "fmt"

// Car 查看对应的iface
// 原始类型 src/runtime/runtime2.go:202 iface
type Car interface {
	Drive()
}

type TrafficTool interface {
	Drive()
}

type Truck struct {
	Model string
}

func (t Truck) Drive() {
	fmt.Println(t.Model)
}

//func (t *Truck) Drive() {
//	fmt.Println(t.Model)
//}

func main() {
	var c Car = Truck{}
	t := c.(Truck)
	fmt.Println(t)
	tt := c.(TrafficTool)
	fmt.Println(tt)

	switch c.(type) {
	case Car:
		fmt.Println("car")
	case TrafficTool:
		fmt.Println("traffic tool")
	}
}
