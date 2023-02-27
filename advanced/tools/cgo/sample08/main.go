package main

//void SayHello(_GoString_ s);
//void SayHello(_GoString_ s);
import "C"
import "fmt"

/*
现在版本的 CGO 代码中 C 语言代码的比例已经很少了，但是我们依然可以进一步以 Go 语言的思维来提炼我们的 CGO 代码。
通过分析可以发现 SayHello 函数的参数如果可以直接使用 Go 字符串是最直接的。
在 Go1.10 中 CGO 新增加了一个 _GoString_ 预定义的 C 语言类型，用来表示 Go 语言字符串。下面是改进后的代码：
// +build go1.10

package main

//void SayHello(_GoString_ s);
import "C"

import (
	"fmt"
)

func main() {
	C.SayHello("Hello, World\n")
}

//export SayHello
func SayHello(s string) {
	fmt.Print(s)
}

虽然看起来全部是 Go 语言代码，但是执行的时候是先从 Go 语言的 main 函数，到 CGO 自动生成的 C 语言版本 SayHello 桥接函数，最后又回到了 Go 语言环境的 SayHello 函数。

思考题: main 函数和 SayHello 函数是否在同一个 Goroutine 里执行？
*/

//export SayHello
func SayHello(s string) {
	fmt.Print(s)
}

func main() {
	C.SayHello("hello, world!\n")
}
