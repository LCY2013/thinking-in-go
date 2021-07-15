package main

import "fmt"

// main : go vet vet_info.go
// # command-line-arguments
//./vet_info.go:9:12: Printf call has arguments but no formatting directives
func main() {
	fmt.Printf("The quick brown fox jumped over lazy dogs: ", 3.14)
}
