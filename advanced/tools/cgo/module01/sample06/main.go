package main

//#include <hello.h>
import "C"

var _ = `
用 Go 重新实现 C 函数

其实 CGO 不仅仅用于 Go 语言中调用 C 语言函数，还可以用于导出 Go 语言函数给 C 语言函数调用。
在前面的例子中，我们已经抽象一个名为 hello 的模块，模块的全部接口函数都在 hello.h 头文件定义：
// hello.h
void SayHello(/*const*/ char* s);

现在我们创建一个 hello.go 文件，用 Go 语言重新实现 C 语言接口的 SayHello 函数:
// hello.go
package main

import "C"

import "fmt"

//export SayHello
func SayHello(s *C.char) {
    fmt.Print(C.GoString(s))
}

我们通过 CGO 的 //export SayHello 指令将 Go 语言实现的函数 SayHello 导出为 C 语言函数。
为了适配 CGO 导出的 C 语言函数，我们禁止了在函数的声明语句中的 const 修饰符。
需要注意的是，这里其实有两个版本的 SayHello 函数：一个 Go 语言环境的；另一个是 C 语言环境的。
cgo 生成的 C 语言版本 SayHello 函数最终会通过桥接代码调用 Go 语言版本的 SayHello 函数。

通过面向 C 语言接口的编程技术，我们不仅仅解放了函数的实现者，同时也简化的函数的使用者。
现在我们可以将 SayHello 当作一个标准库的函数使用（和 puts 函数的使用方式类似）：
package main

//#include <hello.h>
import "C"

func main() {
    C.SayHello(C.CString("Hello, World\n"))
}

`

func main() {
	C.SayHello(C.CString("hello, world!\n"))
}
