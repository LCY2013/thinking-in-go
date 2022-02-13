package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("vim-go")
}

type Resource struct {
	url        string
	polling    bool
	lastPolled int64
}

type Resources struct {
	data []*Resource
	lock *sync.Mutex
}

func Poller(res *Resources) {
	for {
		res.lock.Lock()
		var r *Resource
		for _, v := range res.data {
			if v.polling {
				continue
			}
			if r == nil || v.lastPolled < r.lastPolled {
				r = v
			}
		}
		if r != nil {
			r.polling = true
		}
		res.lock.Unlock()
		if r == nil {
			continue
		}

		// poll this url
		res.lock.Lock()
		r.polling = false
		//r.lastPolled = time.Nanoseconds()
		r.lastPolled = int64(time.Now().Nanosecond())
		res.lock.Unlock()
	}
}
