package sync

import "sync"

type ConcurrentMap[K comparable, V any] struct {
	values map[K]V
	lock   sync.RWMutex
}

// 已经有 key，返回对应的值，然后 loaded = true
// 没有，则放进去，返回 loaded false
// goroutine 1 => ("key1", 1)
// goroutine 2 => ("key1", 2)

func (concurrentMap *ConcurrentMap[K, V]) LoadOrStoreV1(key K, newValue V) (V, bool) {
	concurrentMap.lock.RLock()
	oldVal, ok := concurrentMap.values[key]
	defer concurrentMap.lock.RUnlock()
	if ok {
		return oldVal, true
	}
	concurrentMap.lock.Lock()
	defer concurrentMap.lock.Unlock()
	oldVal, ok = concurrentMap.values[key]
	if ok {
		return oldVal, true
	}
	// goroutine1 先进来，那么这里就会变成 key1 => 1
	// goroutine2 进来，那么这里就会变成 key1 => 2
	concurrentMap.values[key] = newValue
	return newValue, false
}

func (concurrentMap *ConcurrentMap[K, V]) LoadOrStoreV2(key K, newValue V) (V, bool) {
	concurrentMap.lock.RLock()
	oldVal, ok := concurrentMap.values[key]
	concurrentMap.lock.RUnlock()
	if ok {
		return oldVal, true
	}
	concurrentMap.lock.Lock()
	defer concurrentMap.lock.Unlock()
	// goroutine1 先进来，那么这里就会变成 key1 => 1
	// goroutine2 进来，那么这里就会变成 key1 => 2
	concurrentMap.values[key] = newValue
	return newValue, false
}

func (concurrentMap *ConcurrentMap[K, V]) LoadOrStoreV3(key K, newValue V) (V, bool) {
	concurrentMap.lock.RLock()
	oldVal, ok := concurrentMap.values[key]
	concurrentMap.lock.RUnlock()
	if ok {
		return oldVal, true
	}
	concurrentMap.lock.Lock()
	defer concurrentMap.lock.Unlock()
	// goroutine1 先进来，那么这里就会变成 key1 => 1
	// goroutine2 进来，那么这里就还是会成为 key1 => 1
	oldVal, ok = concurrentMap.values[key]
	if ok {
		return oldVal, true
	}
	concurrentMap.values[key] = newValue
	return newValue, false
}
