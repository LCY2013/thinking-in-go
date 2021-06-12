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
package embedding

import "log"

/*
Embedding
内嵌

Go 并不提供典型的，类型驱动的子类化概念，但通过将类型 <内嵌到结构体或接口中， 它就 能 “借鉴” 部分实现。

接口内嵌非常简单。我们之前提到过 io.Reader 和 io.Writer 接口，这里是它们的定义。
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Writer interface {
    Write(p []byte) (n int, err error)
}

io 包也导出了一些其它接口，以此来阐明对象所需实现的方法。
例如 io.ReadWriter 就是个包含 Read 和 Write 的接口。
我们可以通过显示地列出这两个方法来指明io.ReadWriter，
但通过将这两个接口内嵌到新的接口中显然更容易且更具启发性，就像这样:
// ReadWriter is the interface that combines the Reader and Writer interfaces.
type ReadWriter interface {
    Reader
	Writer
}
正如它看起来那样:ReadWriter 能够做任何 Reader 和 Writer 可以做到的事情，它是内嵌接口的联合体 (它们必须是不相交的方法集)。只有接口能被嵌入到接口中。

同样的基本想法可以应用在结构体中，但其意义更加深远。
bufio 包中有 bufio.Reader 和 bufio.Writer 这两个结构体类型，它们每一个都实现了与 io 包中相同意义的接口。
此外， bufio 还通过结合 reader/writer 并将其内嵌到结构体中，实现了带缓冲的 reader/writer:它列出了结构体中的类型，但并未给予它们字段名。
// ReadWriter 存储了指向 Reader 和 Writer 的指针。 // 它实现了 io.ReadWriter。
type ReadWriter struct {
    *Reader  // *bufio.Reader
    *Writer  // *bufio.Writer
}

内嵌的元素为指向结构体的指针，当然它们在使用前必须被初始化为指向有效结构体的指针。
ReadWriter 结构体和通过如下方式定义:
type ReadWriter struct {
    reader *Reader
    writer *Writer
}

但为了提升该字段的方法并满足 io 接口，我们同样需要提供转发的方法，就像这样:
func (rw *ReadWriter) Read(p []byte) (n int, err error) {
    return rw.reader.Read(p)
}

而通过直接内嵌结构体，就能避免如此繁琐。
内嵌类型的方法可以直接引用，这意味着 bufio.ReadWriter 不仅包括 bufio.Reader 和 bufio.Writer 的方法，它还同时满足下列三个接口: io.Reader、io.Writer 以及 io.ReadWriter。

还有种区分内嵌与子类的重要手段。
当内嵌一个类型时，该类型的方法会成为外部类型的方法，但当它们被调用时，该方法的接收者是内部类型，而非外部的。
在我们的例子中，当 bufio.ReadWriter 的 Read 方法被调用时，它与之前写的转发方法具有同样的效果;
接收者 是 ReadWriter 的 reader 字段，而非 ReadWriter 本身。

内嵌同样可以提供便利。这个例子展示了一个内嵌字段和一个常规的命名字段。
type Job struct {
    Command string
    *log.Logger
}

Job 类型现在有了 Log、Logf 和 *log.Logger 的其它方法。
我们当然可以为 Logger 提供一个字段名，但完全不必这么做。
现在，一旦初始化后，我们就能记录 Job 了:
job.Log("starting now...")

Logger 是 Job 结构体的常规字段， 因此我们可在 Job 的构造函数中，通过一般的方式来初 始化它，就像这样:

func NewJob(command string, logger *log.Logger) *Job {
    return &Job{command, logger}
}

*/

type Job struct {
	Command string
	*log.Logger
}

func NewJob(command string, logger *log.Logger) *Job {
	return &Job{command, logger}
}

func JobUseFunc() {
	job := &Job{}
	job.Println("JobUseFunc")
}

//type Reader interface {
//	Read(p []byte) (n int, err error)
//}
//
//type Writer interface {
//	Write(p []byte) (n int, err error)
//}
