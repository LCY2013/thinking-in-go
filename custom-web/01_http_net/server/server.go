package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func fooHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	// 创建一个Foo路由和处理函数
	http.Handle("/foo", http.HandlerFunc(fooHandler))

	// 创建一个bar路由和处理函数
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		if err != nil {
			log.Fatal(err)
			return
		}
	})

	// 监听8080端口
	log.Fatal(http.ListenAndServe(":8080", nil))
}
