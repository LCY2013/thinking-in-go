package main

import "testing"

// 关于defer的性能基准测试（Benchmark）
// go test -bench . defer_test.go
/*
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkFooWithDefer-8         181568036                6.652 ns/op
BenchmarkFooWithoutDefer-8      247294932                4.602 ns/op
PASS
ok      command-line-arguments  3.813s

go1.13后差距不大了，go1.17后带有 defer 的函数执行开销，仅是不带有 defer 的函数的执行开销的 1.45 倍左右，已经达到了几乎可以忽略不计的程度
*/

// sum 	求和
func sum(max int) int {
	total := 0
	for i := 0; i < max; i++ {
		total += i
	}
	return total
}

func fooWithDefer() {
	defer func() {
		sum(10)
	}()
}

func fooWithoutDefer() {
	sum(10)
}

func BenchmarkFooWithDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fooWithDefer()
	}
}

func BenchmarkFooWithoutDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fooWithoutDefer()
	}
}
