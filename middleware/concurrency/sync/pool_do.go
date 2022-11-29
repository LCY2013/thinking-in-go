package sync

import (
	"fmt"
	"sync"
)

type PoolCache struct {
	pool sync.Pool
}

func NewPoolCache() *PoolCache {
	return &PoolCache{
		pool: sync.Pool{
			New: func() any {
				fmt.Println("get cache...")
				return []byte{}
			},
		},
	}
}
