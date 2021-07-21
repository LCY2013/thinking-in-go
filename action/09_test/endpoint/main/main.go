package main

import (
	"fufeng.org/test/endpoint"
	"log"
	"net/http"
)

func init() {
	endpoint.Routes()
}

func main() {
	// endpoint.Routes()
	log.Println("listener : Started : Listening on : 5000")
	http.ListenAndServe(":5000", nil)
}
