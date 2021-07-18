package main

import (
	"errors"
	"fmt"
)

// Positive returns true if the number is positive , false if it is negative
func Positive(num int) (bool, error) {
	if num == 0 {
		return false, errors.New("undefined")
	}
	return num > -1, nil
}

func Check(num int) {
	pos, err := Positive(num)
	if err != nil {
		fmt.Println(num, err)
		return
	}
	if pos {
		fmt.Println(num, "is Positive")
	} else {
		fmt.Println(num, "is negative")
	}

}

func main() {
	Check(1)
	Check(0)
	Check(-1)
}
