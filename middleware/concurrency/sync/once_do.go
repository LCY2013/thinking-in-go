package sync

import (
	"fmt"
	"sync"
)

type OnceClose struct {
	close sync.Once
}

func (oc *OnceClose) Close() error {
	oc.close.Do(func() {
		fmt.Println("close")
	})
	return nil
}

// 这样就会执行多次
/*func (oc OnceClose) Close() error {
	oc.close.Do(func() {
		fmt.Println("close")
	})
	return nil
}*/
