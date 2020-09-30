/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 Ramostear.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-09-17
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package main

import (
	"fmt"
	"os"
	"strings"
)

/**
  命令行参数使用情况
	运行时程序的 命令行参数可以通过os包中一个叫Args的变量来获取;当在os包外部使用该变量时，需要用os.Args来访问。
	$ go run echo.go hello fufeng
	$ go build echo.go
	$ ./echo hello fufeng

	-n用于忽略行尾的换行符，-s sep用于指定分隔字符(默认是空格)
*/
func main() {
	// echo1()
	// echo2()
	echo3()
}

func echo1() {
	var echo string
	for i := 0; i < len(os.Args); i++ {
		echo += " " + os.Args[i]
	}
	fmt.Println(echo)
}

func echo2() {
	var echo string
	for _, args := range os.Args[1:] {
		echo += " " + args
	}
	fmt.Println(echo)
}

func echo3() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}
