/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2019 Ramostear.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-08-09
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

import "fmt"

func main() {
	// testCase01()
	// testCase02()
	testCase03()
}

func testCase01() {
	var a int = 17
	var p *int = &a

	a = 21
	fmt.Println("a = ", a)
	fmt.Println("*p = ", *p)

	case01(11)

	// 通过指针操作 变量 a 所在存储空间的值
	*p = 100
	fmt.Println("a = ", a)
	fmt.Println("*p = ", *p)

	a = 7

	fmt.Println("a = ", a)
	fmt.Println("*p = ", *p)
}

func case01(n int) {
	var b int = 1
	b += n
}

func newPoint() {
	var p *int = new(int)
	*p = 12
}

func testCase02() {
	var a int = 2
	fmt.Println("&a=", &a)
	fmt.Printf("%v\n", a)
	fmt.Printf("%T\n", a)

	var p *int

	// 在heap上申请一片内存
	fmt.Printf("%d\n", *p)
	// 打印go语言格式的字符串
	fmt.Printf("%v", *p)
}

// 测试两个交换
func testCase03() {
	// var one,two int = 1,2
	one, two := 1, 2
	swap1(one, two)
	fmt.Println("one = ", one, ", two = ", two)
	swap2(&one, &two)
	fmt.Println("one = ", one, ", two = ", two)
}

// 测试go语言是值传递
func swap1(one, two int) {
	// go语言语法糖
	one, two = two, one
}

func swap2(one, two *int) {
	// 通过取地址获取真是的存储
	*one, *two = *two, *one
}
