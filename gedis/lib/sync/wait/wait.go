package wait

import (
	"fmt"
	"github.com/LCY2013/thinking-in-go/gedis/lib/logger"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

// Wait is similar with sync.WaitGroup which can wait with timeout
type Wait struct {
	wg sync.WaitGroup
}

// Add adds delta, which may be negative, to the WaitGroup counter.
func (w *Wait) Add(delta int) {
	w.wg.Add(delta)
}

// Done decrements the WaitGroup counter by one
func (w *Wait) Done() {
	w.wg.Done()
}

// Wait blocks until the WaitGroup counter is zero.
func (w *Wait) Wait() {
	w.wg.Wait()
}

// WaitWithTimeout blocks until the WaitGroup counter is zero or timeout
// returns true if timeout
func (w *Wait) WaitWithTimeout(timeout time.Duration) bool {
	c := make(chan bool, 1)
	go func() {
		defer close(c)
		w.wg.Wait()
		c <- true
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

// AsyncDo async do
func (w *Wait) AsyncDo(f func()) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.Try(f)
	}()
}

// SyncDo sync do
func (w *Wait) SyncDo(f func()) {
	w.wg.Add(1)
	defer w.wg.Done()
	w.Try(f)
}

// Try executes f, catching any panic it might spawn. It is safe
// to call from multiple goroutines simultaneously.
func (w *Wait) Try(f func()) {
	defer w.tryRecover()
	f()
}

func (w *Wait) tryRecover() {
	if val := recover(); val != nil {
		var callers [64]uintptr
		n := runtime.Callers(2, callers[:])
		logger.Error(fmt.Sprintf("case: [%s], \ncallers: [%d], \nstack: [%s]", val, n, debug.Stack()))
	}
}
