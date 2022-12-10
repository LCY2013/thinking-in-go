/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-04-20
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
	"io"
	"os"
)

/*
Go 的 defer 语句用于预设一个函数调用(即推迟执行函数)， 该函数会在执行 defer 的函数 返回之前立即执行。
它显得非比寻常， 但却是处理一些事情的有效方式，例如无论以何种路 径返回，都必须释放资源的函数。 典型的例子就是解锁互斥和关闭文件。

推迟诸如 Close 之类的函数调用有两点好处:
第一， 它能确保你不会忘记关闭文件。如果你 以后又为该函数添加了新的返回路径时， 这种情况往往就会发生。
第二，它意味着 “关闭” 离 “打开” 很近， 这总比将它放在函数结尾处要清晰明了。
*/

// Contents returns the file's contents as a string.
func Contents(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close() // f.Close will run when we're finished.
	var result []byte
	buf := make([]byte, 100)
	for {
		n, err := f.Read(buf[0:])
		result = append(result, buf[0:n]...) // append is discussed later.
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err // f will be closed if we return here.
		}
	}
	return string(result), nil // f will be closed if we return here.
}

/*
*
被推迟函数的实参(如果该函数为方法则还包括接收者)在推迟执行时就会求值， 而不是在 调用执行时才求值。
这样不仅无需担心变量值在函数执行时被改变， 同时还意味着单个已推 迟的调用可推迟多个函数的执行。下面是个简单的例子。
*/
func invokeOnCurrent() {
	for i := 0; i < 5; i++ {
		defer fmt.Printf("%d ", i)
	}
}

/*
被推迟的函数按照后进先出(LIFO)的顺序执行，因此以上代码在函数返回时会打印 4 3 2 1 0。
一个更具实际意义的例子是通过一种简单的方法， 用程序来跟踪函数的执行。我们可以编 写一对简单的跟踪例程:
*/
//func trace(s string) string   { fmt.Println("entering:", s) }
//func untrace(s string) string { fmt.Println("leaving:", s) }
//
//// 像这样使用它们:
//func a() {
//	trace("a")
//	defer untrace("a") // 做一些事情....
//}

/*
我们可以充分利用这个特点，即被推迟函数的实参在 defer 执行时才会被求值。 跟踪例程可 针对反跟踪例程设置实参。以下例子:
*/
func trace(s string) string {
	fmt.Println("entering:", s)
	return s
}
func un(s string) {
	fmt.Println("leaving:", s)
}
func a() {
	defer un(trace("a"))
	fmt.Println("in a")
}
func b() {
	defer un(trace("b"))
	fmt.Println("in b")
	a()
}

/*
对于习惯其它语言中块级资源管理的程序员，defer 似乎有点怪异， 但它最有趣而强大的应用 恰恰来自于其基于函数而非块的特点。
在 panic 和 recover 这两节中，我们将看到关于它可能 性的其它例子。
*/
func main() {
	//fmt.Println(Contents(""))
	//invokeOnCurrent()
	//a()
	b()
}
