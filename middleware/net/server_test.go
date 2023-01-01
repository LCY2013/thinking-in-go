package net

import "testing"

func TestServe(t *testing.T) {
	_ = Serve(":8080")
}

func TestServer_StartAndServe(t *testing.T) {
	serve := &Server{
		addr: ":8080",
	}
	_ = serve.StartAndServe()
}
