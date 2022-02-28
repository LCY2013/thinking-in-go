package main

import (
	"sync"
	"testing"
)

var cs = 0 // 模拟临界区要保护的数据
var mu sync.Mutex
var c = make(chan struct{}, 1)

/*
goos: darwin
goarch: amd64
pkg: github.com/lcy2013/sync_mutex_channel_test
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkCriticalSectionSyncByMutex-8                   79013780                15.36 ns/op
BenchmarkCriticalSectionSyncByMutexInParallel-8         19623278                60.60 ns/op
BenchmarkCriticalSectionSyncByChan-8                    23533351                47.93 ns/op
BenchmarkCriticalSectionSyncByChanInParallel-8           4129382               293.8 ns/op
PASS
ok      github.com/lcy2013/sync_mutex_channel_test      5.697s
*/

func criticalSectionSyncByMutex() {
	mu.Lock()
	cs++
	mu.Unlock()
}

func criticalSectionSyncByChan() {
	c <- struct{}{}
	cs++
	<-c
}

func BenchmarkCriticalSectionSyncByMutex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		criticalSectionSyncByMutex()
	}
}

func BenchmarkCriticalSectionSyncByMutexInParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			criticalSectionSyncByMutex()
		}
	})
}

func BenchmarkCriticalSectionSyncByChan(b *testing.B) {
	for n := 0; n < b.N; n++ {
		criticalSectionSyncByChan()
	}
}

func BenchmarkCriticalSectionSyncByChanInParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			criticalSectionSyncByChan()
		}
	})
}
