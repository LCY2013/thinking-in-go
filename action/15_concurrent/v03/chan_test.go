package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkWithNoBuffer(b *testing.B) {
	benchmarkWithBuffer(b, 0)
}

func BenchmarkWithBufferSizeOf1(b *testing.B) {
	benchmarkWithBuffer(b, 1)
}

func BenchmarkWithBufferSizeEqualsToNumberOfWorker(b *testing.B) {
	benchmarkWithBuffer(b, 5)
}

func BenchmarkWithBufferSizeExceedsNumberOfWorker(b *testing.B) {
	benchmarkWithBuffer(b, 25)
}

func BenchmarkWithBufferSizeMuchExceedsNumberOfWorker(b *testing.B) {
	benchmarkWithBuffer(b, 250)
}

func BenchmarkWithBufferSizeMostExceedsNumberOfWorker(b *testing.B) {
	benchmarkWithBuffer(b, 15000)
}

func benchmarkWithBuffer(b *testing.B, chanSize int) {
	for i := 0; i < b.N; i++ {
		ch := make(chan uint32, chanSize)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := uint32(0); j < 10000; j++ {
				ch <- j
			}
			close(ch)
		}()
		var total uint32
		for j := 0; j < 5; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					v, ok := <-ch
					if !ok {
						break
					}
					atomic.AddUint32(&total, v)
				}
			}()
		}
		wg.Wait()
	}
}
