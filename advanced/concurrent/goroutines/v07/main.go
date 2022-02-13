package main

import (
	"fmt"
	"net/http"
)

func serveApp() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "v7 - hello, magic!")
	})

	// app trace
	_ = http.ListenAndServe(":8080", mux)
}

func serveDebug() {
	_ = http.ListenAndServe(":8081", http.DefaultServeMux)
}

func main() {
	go serveDebug()
	serveApp()
}
