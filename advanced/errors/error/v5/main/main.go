package main

import "fmt"

// Positive returns true if the number is positive , false if it is negative
func Positive(num int) *bool {
	if num == 0 {
		return nil
	}
	r := num > -1
	return &r
}

func Check(num int) {
	pos := Positive(num)
	if pos == nil {
		fmt.Println(num, "is neither")
		return
	}
	if *pos {
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
