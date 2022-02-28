package atomic

import (
	"sync"
	"sync/atomic"
	"testing"
)

/*
$ go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/lcy2013/sync_mutex_channel_test/atomic
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkAddSyncByAtomic-8      53621286                21.93 ns/op
BenchmarkReadSyncByAtomic-8     1000000000               0.2719 ns/op
BenchmarkAddSyncByRWMutex-8     14705517                76.90 ns/op
BenchmarkReadSyncByRWMutex-8    34780426                35.30 ns/op
PASS
ok      github.com/lcy2013/sync_mutex_channel_test/atomic       4.459s

$ go test -bench . -cpu=1,4,6,8
goos: darwin
goarch: amd64
pkg: github.com/lcy2013/sync_mutex_channel_test/atomic
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkAddSyncByAtomic        170321643                7.046 ns/op
BenchmarkAddSyncByAtomic-4      53024709                22.45 ns/op
BenchmarkAddSyncByAtomic-6      54196406                22.28 ns/op
BenchmarkAddSyncByAtomic-8      53721466                22.14 ns/op
BenchmarkReadSyncByAtomic       757046976                1.561 ns/op
BenchmarkReadSyncByAtomic-4     1000000000               0.4064 ns/op
BenchmarkReadSyncByAtomic-6     1000000000               0.3074 ns/op
BenchmarkReadSyncByAtomic-8     1000000000               0.2505 ns/op
BenchmarkAddSyncByRWMutex       43943548                25.53 ns/op
BenchmarkAddSyncByRWMutex-4     18505788                58.03 ns/op
BenchmarkAddSyncByRWMutex-6     16704163                68.87 ns/op
BenchmarkAddSyncByRWMutex-8     15579458                73.76 ns/op
BenchmarkReadSyncByRWMutex      79716135                13.81 ns/op
BenchmarkReadSyncByRWMutex-4    27196687                44.60 ns/op
BenchmarkReadSyncByRWMutex-6    27239397                44.27 ns/op
BenchmarkReadSyncByRWMutex-8    27226197                44.01 ns/op
PASS
ok      github.com/lcy2013/sync_mutex_channel_test/atomic       18.049s

*/

var n1 int64

func addSyncByAtomic(delta int64) int64 {
	return atomic.AddInt64(&n1, delta)
}

func readSyncByAtomic() int64 {
	return atomic.LoadInt64(&n1)
}

var n2 int64
var rwmu sync.RWMutex

func addSyncByRWMutex(delta int64) {
	rwmu.Lock()
	n2 += delta
	rwmu.Unlock()
}

func readSyncByRWMutex() int64 {
	var n int64
	rwmu.RLock()
	n = n2
	rwmu.RUnlock()
	return n
}

func BenchmarkAddSyncByAtomic(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			addSyncByAtomic(1)
		}
	})
}

func BenchmarkReadSyncByAtomic(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			readSyncByAtomic()
		}
	})
}

func BenchmarkAddSyncByRWMutex(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			addSyncByRWMutex(1)
		}
	})
}

func BenchmarkReadSyncByRWMutex(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			readSyncByRWMutex()
		}
	})
}
