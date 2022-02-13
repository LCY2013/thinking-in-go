package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "v6 - hello, magic!")
	})

	// debug 端点
	go http.ListenAndServe(":8081", http.DefaultServeMux)
	// app trace
	http.ListenAndServe(":8080", mux)
}
