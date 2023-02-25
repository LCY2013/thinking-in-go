package common

import (
	"runtime/debug"
)

func GO(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()

		f()
	}()
}
