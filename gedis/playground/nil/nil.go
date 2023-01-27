package main

import "fmt"

// main src/builtin/builtin.go:111 nil
// var nil Type // Type must be a pointer, channel, func, interface, map, or slice type
func main() {
	var a *int
	fmt.Println(a == nil)
	var b map[string]int
	fmt.Println(b == nil)

}
