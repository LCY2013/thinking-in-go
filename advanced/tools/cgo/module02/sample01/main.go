package main

//static const char* cs = "hello world";
import "C"
import "fufeng.org/advanced/tools/cgo/module02/sample01/cgo_helper"

//import "./cgo_helper"

/*
CGO 基础

要使用 CGO 特性，需要安装 C/C++ 构建工具链，
在 macOS 和 Linux 下是要安装 GCC，
在 windows 下是需要安装 MinGW 工具。

同时需要保证环境变量 CGO_ENABLED 被设置为 1，这表示 CGO 是被启用的状态。

在本地构建时 CGO_ENABLED 默认是启用的，当交叉构建时 CGO 默认是禁止的。

比如要交叉构建 ARM 环境运行的 Go 程序，需要手工设置好 C/C++ 交叉构建的工具链，同时开启 CGO_ENABLED 环境变量。
然后通过 import "C" 语句启用 CGO 特性。

import "C" 语句

如果在 Go 代码中出现了 import "C" 语句则表示使用了 CGO 特性，紧跟在这行语句前面的注释是一种特殊语法，
里面包含的是正常的 C 语言代码。当确保 CGO 启用的情况下，还可以在当前目录中包含 C/C++ 对应的源文件。

举个最简单的例子：
package main


/*#include <stdio.h>

void printint(int v) {
    printf("printint: %d\n", v);
}*
import "C"

func main() {
	v := 42
	C.printint(C.int(v))
}

这个例子展示了 cgo 的基本使用方法。
开头的注释中写了要调用的 C 函数和相关的头文件，头文件被 include 之后里面的所有的 C 语言元素都会被加入到”C” 这个虚拟的包中。
需要注意的是，import "C" 导入语句需要单独一行，不能与其他包一同 import。
向 C 函数传递参数也很简单，就直接转化成对应 C 语言类型传递就可以。
如上例中 C.int(v) 用于将一个 Go 中的 int 类型值强制类型转换转化为 C 语言中的 int 类型值，然后调用 C 语言定义的 printint 函数进行打印。

需要注意的是，Go 是强类型语言，所以 cgo 中传递的参数类型必须与声明的类型完全一致，而且传递前必须用”C” 中的转化函数转换成对应的 C 类型，不能直接传入 Go 中类型的变量。
同时通过虚拟的 C 包导入的 C 语言符号并不需要是大写字母开头，它们不受 Go 语言的导出规则约束。

cgo 将当前包引用的 C 语言符号都放到了虚拟的 C 包中，同时当前包依赖的其它 Go 语言包内部可能也通过 cgo 引入了相似的虚拟 C 包，
但是不同的 Go 语言包引入的虚拟的 C 包之间的类型是不能通用的。这个约束对于要自己构造一些 cgo 辅助函数时有可能会造成一点的影响。

比如我们希望在 Go 中定义一个 C 语言字符指针对应的 CChar 类型，然后增加一个 GoString 方法返回 Go 语言字符串：
package cgo_helper

//#include <stdio.h>
import "C"

type CChar C.char

func (p *CChar) GoString() string {
    return C.GoString((*C.char)(p))
}

func PrintCString(cs *C.char) {
    C.puts(cs)
}

现在我们可能会想在其它的 Go 语言包中也使用这个辅助函数：
func main() {
	cgo_helper.PrintCString(C.cs)
}


这段代码是不能正常工作的，因为当前 main 包引入的 C.cs 变量的类型是当前 main 包的 cgo 构造的虚拟的 C 包下的 *char 类型（具体点是 *C.char，更具体点是 *main.C.char），
它和 cgo_helper 包引入的 *C.char 类型（具体点是 *cgo_helper.C.char）是不同的。
在 Go 语言中方法是依附于类型存在的，不同 Go 包中引入的虚拟的 C 包的类型却是不同的（main.C 不等 cgo_helper.C），
这导致从它们延伸出来的 Go 类型也是不同的类型（*main.C.char 不等 *cgo_helper.C.char），这最终导致了前面代码不能正常工作。

有 Go 语言使用经验的用户可能会建议参数转型后再传入。
但是这个方法似乎也是不可行的，因为 cgo_helper.PrintCString 的参数是它自身包引入的 *C.char 类型，在外部是无法直接获取这个类型的。
换言之，一个包如果在公开的接口中直接使用了 *C.char 等类似的虚拟 C 包的类型，其它的 Go 包是无法直接使用这些类型的，
除非这个 Go 包同时也提供了 *C.char 类型的构造函数。因为这些诸多因素，如果想在 go test 环境直接测试这些 cgo 导出的类型也会有相同的限制。

*/

// main cannot use (*_Cvar_cs) (variable of type *_Ctype_char) as type *cgo_helper._Ctype_char in argument to cgo_helper.PrintCString
func main() {
	//cgo_helper.PrintCString(C.cs)
	cgo_helper.PrintCString(cgo_helper.GenCChar("hello world\n"))
}
