//go:build linux && !386 && !arm
// +build linux,!386,!arm

package main

import "fmt"

func main() {
	fmt.Println("hello, go build!")
}
