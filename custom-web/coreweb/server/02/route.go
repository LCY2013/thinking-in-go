package main

import (
	"github.com/lcy2013/custom-web/coreweb/server/02/framework"
)

func registerRouter(core *framework.Core) {
	// core.Get("foo", framework.TimeoutHandler(FooControllerHandler, time.Second*1))
	core.Get("foo", FooControllerHandler)
}
