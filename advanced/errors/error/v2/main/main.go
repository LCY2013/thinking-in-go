package main

import (
	"fmt"
	"math"
)

// Positive returns true if the number is positive , false if it is negative
func Positive(num int) bool {
	return num > -1
}

func Check(num int) {
	if Positive(num) {
		fmt.Println(num, "is Positive")
	} else {
		fmt.Println(num, "is negative")
	}
}

func main() {
	Check(1)
	Check(0)
	Check(-1)
	num := math.Exp2(32) + 1000000000000000000
	var n int
	n = int(num)
	fmt.Println(num, n)
}
