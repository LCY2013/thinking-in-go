package main

/*
int sum(int a, int b) {
	return a + b;
}
*/
import "C"

import "fmt"

// main : go tool cgo main.go
func main() {
	fmt.Println(C.sum(1, 2))
}
