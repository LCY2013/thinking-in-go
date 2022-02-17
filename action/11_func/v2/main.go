package main

import "time"

func setup(task string) func() {
	println("do some setup stuff for", task)
	return func() {
		println("do some teardown stuff for", task)
	}
}

func main() {
	teardown := setup("demo")
	defer teardown()
	println("do some business stuff")

	// go
	time.AfterFunc(time.Second*2, func() {
		println("time after func")
	})
	time.Sleep(time.Second * 3)
}
