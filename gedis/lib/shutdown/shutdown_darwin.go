//go:build darwin

/*
//go:build功能和// +build一样。

只不过在go 1.17这个版本才实现对//go:build的支持。

//go:build xxx后必须同时有// +build xxx，否则编译器就会报错。
*/

package shutdown

import (
	"os"
	"syscall"
)

var Signals = []os.Signal{
	os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGSTOP,
	syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
	syscall.SIGABRT, syscall.SIGSYS, syscall.SIGTERM,
}
