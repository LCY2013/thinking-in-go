package benchmark

import (
	"fmt"
	"strconv"
	"testing"
)

// 基准测试

// BenchmarkSprintf 关于 Sprintf 相关性能的基准测试
// go test -v -run="none" -bench="BenchmarkSprintf"
// none 表示在运行定制的基准测试之前没有其他的单元测试运行，上面两个参数都支持正则表达式
// 基准测试结果如下：
//goos: darwin
//goarch: amd64
//pkg: fufeng.org/test/benchmark
//cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
//BenchmarkSprintf
//BenchmarkSprintf-8      14903394(本次调用的次数)                75.00 ns/op(每次操作的耗时)
//PASS
//ok      fufeng.org/test/benchmark       1.794s(本次基准测试的耗时，如果想让时间更长一点可以使用-benchtime)
// go test -v -run="none" -bench="BenchmarkSprintf" -benchtime="2s"
func BenchmarkSprintf(b *testing.B) {
	number := 10

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%d", number)
	}
}

// BenchmarkFormat 对 strconv.FormatInt 函数进行基准测试
func BenchmarkFormat(b *testing.B) {
	number := int64(10)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		strconv.FormatInt(number, 10)
	}
}

// BenchmarkItoa 对 strconv.Itoa 函数进行基准测试
func BenchmarkItoa(b *testing.B) {
	number := 10

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		strconv.Itoa(number)
	}
}

// 三个基准测试一起运行
// go test -v -run="none" -bench=. -benchtime="2s"
/*
goos: darwin
goarch: amd64
pkg: fufeng.org/test/benchmark
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkSprintf
BenchmarkSprintf-8      30535633                78.86 ns/op
BenchmarkFormat
BenchmarkFormat-8       841094868                2.763 ns/op
BenchmarkItoa
BenchmarkItoa-8         826659973                2.868 ns/op
PASS
ok      fufeng.org/test/benchmark       8.229s
*/
// go test -v -run="none" -bench=. -benchtime="2s" -benchmem
/*
goos: darwin
goarch: amd64
pkg: fufeng.org/test/benchmark
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkSprintf
BenchmarkSprintf-8      30282272                75.38 ns/op            2 B/op(每次操作分配的字节数)          1 allocs/op(表示每次操作从堆上分配内存的次数)
BenchmarkFormat
BenchmarkFormat-8       839742630                2.865 ns/op           0 B/op          0 allocs/op
BenchmarkItoa
BenchmarkItoa-8         830595536                2.897 ns/op           0 B/op          0 allocs/op
PASS
ok      fufeng.org/test/benchmark       7.866s
*/
