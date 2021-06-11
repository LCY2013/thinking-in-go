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
package blank_identifier

import (
	"log"
	"os"
)

/*
Unused imports and variables
未使用的导入和变量

若导入某个包或声明某个变量而不使用它就会产生错误。
未使用的包会让程序膨胀并拖慢编 译速度，而已初始化但未使用的变量不仅会浪费计算能力，还有可能暗藏着更大的 Bug。
然而在程序开发过程中，经常会产生未使用的导入和变量。
虽然以后会用到它们， 但为了完成编译又不得不删除它们才行，这很让人烦恼。空白标识符就能提供一个工作空间。

这个写了一半的程序有两个未使用的导入(fmt 和 io)以及一个未使用的变量(fd)，因此它不能编译，但若到目前为止代码还是正确的，我们还是很乐意看到它们的。
*/

func UnusedImports() {
	fd, err := os.Open("test.go")
	if err != nil {
		log.Fatal(err)
	}
	//TODO: use fd
	fd.Close()
}

/*
要让编译器停止关于未使用导入的抱怨，需要空白标识符来引用已导入包中的符号。
同样，将未使用的变量 fd 赋予空白标识符也能关闭未使用变量错误。
该程序的以下版本可以编译。

按照惯例，我们应在导入并加以注释后，再使全局声明导入错误静默，这样可以让它们更易找到，并作为以后清理它的提醒。
*/

func UnusedImportsError() {
	fd, err := os.Open("test.go")
	if err != nil {
		log.Fatal(err)
	}
	// 通过 _ 丢弃异常信息
	_ = fd
}
