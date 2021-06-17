/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-06-17
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package errors

import (
	"os"
	"syscall"
)

/*
Errors  错误

库例程通常需要向调用者返回某种类型的错误提示。
之前提到过，Go 语言的多值返回特性， 使得它在返回常规的值时，还能轻松地返回详细的错误描述。
按照约定，错误的类型通常为 error，这是一个内建的简单接口。
type error interface {
	Error() string
}

库的编写者通过更丰富的底层模型可以轻松实现这个接口，这样不仅能看见错误， 还能提供 一些上下文。
例如，os.Open 可返回一个 os.PathError。
// PathError 记录一个错误以及产生该错误的路径和操作。
type PathError struct {
	Op string // "open"、"unlink" 等等。
	Path string // 相关联的文件。
	Err error // 由系统调用返回。
}
func (e *PathError) Error() string {
    return e.Op + " " + e.Path + ": " + e.Err.Error()
}

PathError 的 Error 会生成如下错误信息:
open /etc/passwx: no such file or directory

这种错误包含了出错的文件名、操作和触发的操作系统错误，即便在产生该错误的调用和输出的错误信息相距甚远时，它也会非常有用，这比苍白的 “不存在该文件或目录” 更具说明性。

错误字符串应尽可能地指明它们的来源，例如产生该错误的包名前缀。
例如在 image 包中， 由于未知格式导致解码错误的字符串为 “image: unknown format”。

若调用者关心错误的完整细节，可使用类型选择或者类型断言来查看特定错误，并抽取其细节。
对于 PathErrors，它应该还包含检查内部的 Err 字段以进行可能的错误恢复。
for try := 0; try < 2; try++ {
	file, err := os.Create(filename)
	if err == nil {
		return
	}
	if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOSPC {
		deleteTempFiles() // 恢复一些空间。
		continue
	}
	return
}

这里的第二条 if 是另一种类型断言。
若它失败，ok 将为 false，而 e 则为 nil.
若它成功，ok 将为 true，这意味着该错误属于 *os.PathError 类型，而 e 能够检测关于该错误的更多信息。

*/

func errShow(filename string) {
	for try := 0; try < 2; try++ {
		file, err := os.Create(filename)
		if err == nil {
			_ = file
			return
		}
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOSPC {
			//deleteTempFiles() // 恢复一些空间。
			continue
		}
		return
	}
}
