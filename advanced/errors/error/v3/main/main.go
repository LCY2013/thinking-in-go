package main

import "fmt"

// Positive returns true if the number is positive , false if it is negative
func Positive(num int) (bool, bool) {
	if num == 0 {
		return false, false
	}
	return num > -1, true
}

func Check(num int) {
	/*if pos, ok := Positive(num); !ok {
		fmt.Println(num, "is neither")
	} else {
		if pos {
			fmt.Println(num, "is Positive")
		}else {
			fmt.Println(num, "is negative")
		}
	}*/
	pos, ok := Positive(num)
	if !ok {
		fmt.Println(num, "is neither")
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
