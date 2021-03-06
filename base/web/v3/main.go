/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-09-29
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// http处理函数定义
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

// http handler 用于处理返回url的uri部分
func handler(w http.ResponseWriter, req *http.Request) {
	_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", req.Method, req.URL, req.Proto)
	for k, v := range req.Header {
		_, _ = fmt.Fprintf(w, "Head[%q]=%q\n", k, v)
	}
	_, _ = fmt.Fprintf(w, "Host=%s\n", req.Host)
	_, _ = fmt.Fprintf(w, "RemoteAddr=%q\n", req.RemoteAddr)
	if err := req.ParseForm(); err != nil {
		log.Print(err)
	}

	for k, v := range req.Form {
		_, _ = fmt.Fprintf(w, "Form[%q]=%q\n", k, v)
	}
}
