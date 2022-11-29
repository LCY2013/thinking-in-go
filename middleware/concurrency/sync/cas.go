package sync

import (
	"sync"
	"sync/atomic"
)

type CasLock struct {
	state uint32
	lock  sync.RWMutex
}

func (cl *CasLock) Lock() bool {
	if atomic.LoadUint32(&cl.state) == 0 {
		return true
	}
	return cl.doSlow()
}

func (cl *CasLock) doSlow() bool {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	if cl.state == 0 {
		defer atomic.StoreUint32(&cl.state, 1)
	}
	return true
}
