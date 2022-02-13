package main

import (
	"fmt"
	"log"
	"net/http"
)

func serveApp() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "v8 - hello, magic!")
	})

	// app trace
	if err := http.ListenAndServe(":8080", mux); err != nil {
		// os.Exit会直接无条件终止程序，defer 不会被调用到
		log.Fatal(err)
	}
}

func serveDebug() {
	if err := http.ListenAndServe(":8081", http.DefaultServeMux); err != nil {
		log.Fatal(err)
	}
}

func main() {
	go serveApp()
	go serveDebug()
	select {}
}
