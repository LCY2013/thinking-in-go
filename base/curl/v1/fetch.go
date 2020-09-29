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
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

/*
这个程序将获取对应的 url，并将其源文本打印出来;
这个例子的灵感来源于curl工具(译注:unix下的一个网络相关的工 具)。
当然了，curl提供的功能更为复杂丰富，这里我们只编写最简单的样例。
*/
func main() {
	for _, url := range os.Args[1:] {

		if !strings.HasPrefix(url, "http://") {
			continue
		}

		resp, err := http.Get(url)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "fetch err : %s\n", err)
			os.Exit(7)
		}
		// io.Copy(resp.Body,os.Stdout)

		fmt.Printf("http response status %d\n", resp.StatusCode)

		urlContent, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "read http response error : %s\n", err)
			os.Exit(7)
		}

		fmt.Printf("%s\n", urlContent)
	}
}
