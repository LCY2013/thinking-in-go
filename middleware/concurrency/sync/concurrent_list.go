package sync

import "sync"

// ConcurrentList 封装一个协程安全的List
type ConcurrentList[T any] struct {
	List[T]
	lock sync.RWMutex
}

func (cl *ConcurrentList[T]) Get(index int) (T, error) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.List.Get(index)
}

func (cl *ConcurrentList[T]) Append(t T) error {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.List.Append(t)
}
