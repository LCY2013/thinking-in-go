package rwmutx

import (
	"sync"
	"testing"
)

/*
goos: darwin
goarch: amd64
pkg: github.com/lcy2013/sync_mutex_channel_test/rwmutx
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkWriteSyncByMutex-8             20367465                59.72 ns/op
BenchmarkReadSyncByMutex-8              19002603                65.84 ns/op
BenchmarkReadSyncByRWMutex-8            30787990                39.09 ns/op
BenchmarkWriteSyncByRWMutex-8           15851838                73.22 ns/op
PASS
ok      github.com/lcy2013/sync_mutex_channel_test/rwmutx       5.531s
*/

var cs1 = 0 // 模拟临界区要保护的数据
var mu1 sync.Mutex
var cs2 = 0 // 模拟临界区要保护的数据
var mu2 sync.RWMutex

func BenchmarkWriteSyncByMutex(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu1.Lock()
			cs1++
			mu1.Unlock()
		}
	})
}
func BenchmarkReadSyncByMutex(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu1.Lock()
			_ = cs1
			mu1.Unlock()
		}
	})
}
func BenchmarkReadSyncByRWMutex(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu2.RLock()
			_ = cs2
			mu2.RUnlock()
		}
	})
}
func BenchmarkWriteSyncByRWMutex(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu2.Lock()
			cs2++
			mu2.Unlock()
		}
	})
}
