package main

import "os"

// main : go fmt fmt_info.go
// fmt 前:
// if _, err := os.Open(""); err != nil{}
// fmt 后:
// if _, err := os.Open(""); err != nil {
//	}
func main() {
	if _, err := os.Open(""); err != nil {
	}
}
