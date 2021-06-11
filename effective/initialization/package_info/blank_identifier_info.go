/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-06-11
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package package_info

import (
	"fmt"
	"os"
)

/*
The blank identifier
空白标识符

在 for-range 循环和 映射 中提过几次空白标识符。
空白标识符可被赋予或声明为任何类型的任何值，而其值会被无害地丢弃。
它有点像 Unix 中的 /dev/null 文件:它表示只写的值，在需要变量但不需要实际值的地方用作占位符。

多重赋值中的空白标识符

for range 循环中对空白标识符的用法是一种具体情况，更一般的情况即为多重赋值。

若某次赋值需要匹配多个左值，但其中某个变量不会被程序使用，
那么用空白标识符来代替该变量可避免创建无用的变量，并能清楚地表明该值将被丢弃。
例如，当调用某个函数时， 它会返回一个值和一个错误，但只有错误很重要， 那么可使用空白标识符来丢弃无关的值。
if _, err := os.Stat(path); os.IsNotExist(err) {
    fmt.Printf("%s does not exist\n", path)
}

你偶尔会看见为忽略错误而丢弃错误值的代码，这是种糟糕的实践。请务必检查错误返回，它们会提供错误的理由。
// 烂代码!若路径不存在，它就会崩溃。
fi, _ := os.Stat(path)
if fi.IsDir() {
    fmt.Printf("%s is a directory\n", path)
}
*/

func StatPathIsNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("%s does not exist\n", path)
	}
}
