/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-06-10
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
	"math"
)

/*
Go 中的常量就是不变量。它们在编译时创建，即便它们可能是函数中定义的局部变量。
常量 只能是数字、字符(符文)、字符串或布尔值。由于编译时的限制， 定义它们的表达式必须 也是可被编译器求值的常量表达式。
例如 1<<3 就是一个常量表达式，而 math.Sin(math.Pi/4) 则不是，因为对 math.Sin 的函数调用在运行时才会发生。

在 Go 中，枚举常量使用枚举器 iota 创建。由于 iota 可为表达式的一部分，而表达式可以被 隐式地重复，这样也就更容易构建复杂的值的集合了。
*/
func main() {
	fmt.Println("---------------------", "move")
	move()
	fmt.Println("---------------------", "ByteSizeFmt")
	ByteSizeFmt()
}

func move() {
	fmt.Println(1 << 3)
	fmt.Println(1 >> 3)
	fmt.Println(math.Sin(math.Pi / 4))
	fmt.Println(math.Sin(math.Pi))
}

func ByteSizeFmt() {
	// 表达式 YB 会打印出 1.00YB，而 ByteSize(1e13) 则会打印出 9.09。
	b := ByteSize(1e13)
	fmt.Println(b)
}

// ByteSize 字节大小定义
type ByteSize float64

const (
	// 通过赋予空白标志符来忽略第一个值
	_           = iota
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

// String 字节大小转换
func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}
