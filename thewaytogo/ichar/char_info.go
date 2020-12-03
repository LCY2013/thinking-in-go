/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-12-03
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

/*
格式化说明符 %c 用于表示字符；当和字符配合使用时，%v 或 %d 会输出用于表示该字符的整数；%U 输出格式为 U+hhhh 的字符串（另一个示例见第 5.4.4 节）。

包 unicode 包含了一些针对测试字符的非常有用的函数（其中 ch 代表字符）：

判断是否为字母：unicode.IsLetter(ch)
判断是否为数字：unicode.IsDigit(ch)
判断是否为空白符号：unicode.IsSpace(ch)
这些函数返回一个布尔值。包 utf8 拥有更多与 rune 相关的函数。
*/
func main() {
	var ch1 byte = 65
	var ch2 byte = '\x41'
	fmt.Printf("%c %c \n", ch1, ch2)

	charInfo()
}

func charInfo() {
	var ch int = '\u0041'
	var ch2 int = '\u03B2'
	var ch3 int = '\U00101234'
	fmt.Printf("%d - %d - %d\n", ch, ch2, ch3) // integer
	fmt.Printf("%c - %c - %c\n", ch, ch2, ch3) // character
	fmt.Printf("%X - %X - %X\n", ch, ch2, ch3) // UTF-8 bytes
	fmt.Printf("%U - %U - %U\n", ch, ch2, ch3) // UTF-8 code point
}
