/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-06-09
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

import "fmt"

func main() {
	fmt.Println("myAppend001")
	myAppend001()
	fmt.Println("myAppend002")
	myAppend002()
}

/*
append 函数的签名不同于前面我们自 定义的Append 函数。大致来说，它就像这样:
 func append(slice []T, elements ...T) []T

其中的 T 为任意给定类型的占位符。实际上，你无法在 Go 中编写一个类型 T 由调用者决定 的函数。这也就是为何 append 为内建函数的原因:它需要编译器的支持。

append 会在切片末尾追加元素并返回结果。我们必须返回结果， 原因与我们手写的 Append 一样，即底层数组可能会被改变。以下简单的例子
*/
func myAppend001() {
	x := []int{1, 2, 3}
	x = append(x, 4, 5, 6)
	fmt.Println(x)
}

/*
但如果我们要像 Append 那样将一个切片追加到另一个切片中呢? 很简单:在调用的地方使 用 ...，就像我们在上面调用 Output 那样。以下代码片段的输出与上一个相同。
*/
func myAppend002() {
	x := []int{4, 5, 6}
	y := []int{7, 8, 9}
	// 这里的y必须加上...
	// 如果没有 ... 它就会由于类型错误而无法编译，因为 y 不是 int 类型的。
	z := append(x, y...)
	fmt.Println(z)
}
