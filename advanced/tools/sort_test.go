package tools

import (
	"math/rand"
	"testing"
	"time"
)

const N = 100_000

// https://github.com/golang/exp/blob/master/slices/sort_benchmark_test.go

func makeRandomFloat64s(n int) []float64 {
	rand.Seed(time.Now().Unix())
	floats := make([]float64, n)
	for i := 0; i < n; i++ {
		floats[i] = rand.Float64()
	}
	return floats
}

func BenchmarkSortFloat64FastV1(b *testing.B) {
	//b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		//b.StopTimer()
		SortFloat64FastV1(makeRandomFloat64s(N))
		//b.StartTimer()
	}
}

func BenchmarkSortFloat64FastV2(b *testing.B) {
	//b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		//b.StopTimer()
		SortFloat64FastV2(makeRandomFloat64s(N))
		//b.StartTimer()
	}
}
