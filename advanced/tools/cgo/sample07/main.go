package main

//void SayHello(char* s);
import "C"
import "fmt"

/*
面向 C 接口的 Go 编程

在前面的例子中，我们的全部 CGO 代码都在一个 Go 文件中。
然后，通过面向 C 接口编程的技术将 SayHello 分别拆分到不同的 C 文件，而 main 依然是 Go 文件。
再然后，是用 Go 函数重新实现了 C 语言接口的 SayHello 函数。
但是对于目前的例子来说只有一个函数，要拆分到三个不同的文件确实有些繁琐了。

正所谓合久必分、分久必合，我们现在尝试将例子中的几个文件重新合并到一个 Go 文件。下面是合并后的成果：
package main

//void SayHello(char* s);
import "C"

import (
    "fmt"
)

func main() {
    C.SayHello(C.CString("Hello, World\n"))
}

//export SayHello
func SayHello(s *C.char) {
    fmt.Print(C.GoString(s))
}

*/

//export SayHello
func SayHello(s *C.char) {
	fmt.Print(C.GoString(s))
}

func main() {
	C.SayHello(C.CString("hello, world!"))
}
