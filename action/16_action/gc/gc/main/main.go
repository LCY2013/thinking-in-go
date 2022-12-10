package main

/*
# gc/main
./main.go:4:6: can inline main
./main.go:5:14: make([]int, 10240) escapes to heap
*/

// main go build -gcflags="-m" 监测内存逃逸
func main() {
	var m = make([]int, 10240)
	println(m[0])
}
