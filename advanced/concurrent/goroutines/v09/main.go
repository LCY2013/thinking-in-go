package main

import (
	"context"
	"fmt"
	"net/http"
)

func serve(addr string, handler http.Handler, stop <-chan struct{}) error {
	server := http.Server{Addr: addr, Handler: handler}

	go func() {
		// 等待停止信号
		<-stop
		_ = server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}

func serveApp(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "v8 - hello, magic!")
	})

	// app trace
	return serve(":8080", mux, stop)
}

func serveDebug(stop <-chan struct{}) error {
	// debug trace
	return serve(":8081", http.DefaultServeMux, stop)
}

func main() {
	done := make(chan error, 2)
	stop := make(chan struct{})
	go func() {
		done <- serveApp(stop)
	}()
	go func() {
		done <- serveDebug(stop)
	}()

	var stopped bool
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			fmt.Printf("error: %+v", err)
		}
		if !stopped {
			stopped = true
			close(stop)
		}
	}
}
