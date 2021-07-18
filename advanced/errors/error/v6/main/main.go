package main

import "fmt"

// Positive returns true if the number is positive , false if it is negative
func Positive(num int) bool {
	if num == 0 {
		panic("undefined")
	}
	return num > -1
}

func Check(num int) {
	defer func() {
		if recover() != nil {
			fmt.Println(num, "is neither")
		}
	}()
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
}
