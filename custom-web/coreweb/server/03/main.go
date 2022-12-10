package main

import (
	"net/http"

	"github.com/lcy2013/custom-web/coreweb/server/03/framework"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}
	server.ListenAndServe()
}
