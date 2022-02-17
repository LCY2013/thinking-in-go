package main

import (
	"errors"
	"fmt"
)

type MyError struct {
	e string
}

func (e *MyError) Error() string {
	return e.e
}

func main() {
	var err = &MyError{"MyError error demo"}
	err1 := fmt.Errorf("wrap err1: %w", err)
	err2 := fmt.Errorf("wrap err2: %w", err1)
	var e *MyError
	// 开始不等，在As之后才会相等
	println(e == err)
	if errors.As(err2, &e) {
		println("as MyError is on the chain of err2")
		println(e == err)

	} else {
		println("as MyError is not on the chain of err2")
	}

	if errors.Is(err2, err) {
		println("is MyError is on the chain of err2")
		println(err2 == err)
		return
	} else {
		println("is MyError is not on the chain of err2")
	}
}
