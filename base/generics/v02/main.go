package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Min go 范性实现
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println(Min(1, 2))
	fmt.Println(Min(1.0, 2.0))
}
