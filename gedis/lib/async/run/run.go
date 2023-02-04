package run

import (
	"github.com/LCY2013/thinking-in-go/gedis/lib/logger"
	"runtime/debug"
)

func GO(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(string(debug.Stack()))
			}
		}()

		f()
	}()
}
