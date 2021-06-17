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
	"log"
)

/*
当 panic 被调用后(包括不明确的运行时错误，例如切片检索越界或类型断言失败)， 程序将立刻终止当前函数的执行，并开始回溯 Go 程的栈，运行任何被推迟的函数。
若回溯到达 Go 程栈的顶端，程序就会终止。
不过我们可以用内建的 recover 函数来重新或来取回 Go 程的控制权限并使其恢复正常执行。

调用 recover 将停止回溯过程，并返回传入 panic 的实参。
由于在回溯时只有被推迟函数中的代码在运行，因此 recover 只能在被推迟的函数中才有效。

recover 的一个应用就是在服务器中终止失败的 Go 程而无需杀死其它正在执行的 Go 程。

func server(workChan <-chan *Work) {
	for work := range workChan {
		go safelyDo(work)
	}
}
func safelyDo(work *Work) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed:", err)
		}
	}()
	do(work)
}

在此例中，若 do(work) 触发了 Panic，其结果就会被记录，而该 Go 程会被干净利落地结束，不会干扰到其它 Go 程。
我们无需在推迟的闭包中做任何事情， recover 会处理好这一切。

由于直接从被推迟函数中调用 recover 时不会返回 nil，因此被推迟的代码能够调用本身使用了 panic 和 recover 的库函数而不会失败。
例如在 safelyDo 中，被推迟的函数可能在调用 recover 前先调用记录函数，而该记录函数应当不受 Panic 状态的代码的影响。

通过恰当地使用恢复模式，do 函数(及其调用的任何代码)可通过调用 panic 来避免更坏的结果。
我们可以利用这种思想来简化复杂软件中的错误处理。
让我们看看 regexp 包的理想化版本，它会以局部的错误类型调用 panic 来报告解析错误。
以下是一个 error 类型的 Error 方 法和一个 Compile 函数的定义:
// Error 是解析错误的类型，它满足 error 接口。
type Error string
func (e Error) Error() string {
	return string(e)
}

// error 是 *Regexp 的方法，它通过用一个 Error 触发 Panic 来报告解析错误。
func (regexp *Regexp) error(err string) {
	panic(Error(err))
}
// go1.15.5/src/regexp/regexp.go:132
// Compile 返回该正则表达式解析后的表示。
func Compile(str string) (regexp *Regexp, err error) {
	regexp = new(Regexp)
	// doParse will panic if there is a parse error.
	defer func() {
		if e := recover(); e != nil {
			regexp = nil    // 清理返回值。
			err = e.(Error) // 若它不是解析错误，将重新触发 Panic。
		}
	}()
	return regexp.doParse(str), nil
}

若 doParse 触发了 Panic，恢复块会将返回值设为 nil —被推迟的函数能够修改已命名的返回值。
在 err 的赋值过程中，我们将通过断言它是否拥有局部类型 Error 来检查它。
若它没有，类型断言将会失败，此时会产生运行时错误，并继续栈的回溯，仿佛一切从未中断过一样。
该检查意味着若发生了一些像索引越界之类的意外，那么即便我们使用了 panic 和 recover 来处理解析错误，代码仍然会失败。

通过适当的错误处理，error 方法(由于它是个绑定到具体类型的方法，因此即便它与内建的 error 类型名字相同也没有关系)
能让报告解析错误变得更容易，而无需手动处理回溯的解析栈:
if pos == 0 {
    re.error("'*' illegal at start of expression")
}

尽管这种模式很有用，但它应当仅在包内使用。
Parse 会将其内部的 panic 调用转为 error 值，它并不会向调用者暴露出 panic。这是个值得遵守的良好规则。

顺便一提，这种重新触发Panic的惯用法会在产生实际错误时改变Panic的值。
然而，不管是原始的还是新的错误都会在崩溃报告中显示，因此问题的根源仍然是可见的。
这种简单的重新触发Panic的模型已经够用了，毕竟他只是一次崩溃。
但若你只想显示原始的值，也可以多写一点代码来过滤掉不需要的问题，然后用原始值再次触发Panic。
*/

type Work struct {
}

func Server(workChan chan *Work) {
	go func() {
		for work := range workChan {
			go safelyDo(work)
		}
	}()
}
func safelyDo(work *Work) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed:", err)
		}
	}()
	do(work)
}

func do(work *Work) {
	panic("panic for purpose...")
}

func SendMessage(workChan chan *Work) {
	go func() {
		workChan <- &Work{}
	}()
}

// Run 封装go程
func Run(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("run failed: ", err)
		}
	}()
	// 执行具体协程
	go fn()
}
