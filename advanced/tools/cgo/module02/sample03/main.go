//go:build debug

// +build: darwin,!cgo

package main

var buildMode = "debug"

/*
build tag 条件编译

build tag 是在 Go 或 cgo 环境下的 C/C++ 文件开头的一种特殊的注释。
条件编译类似于前面通过 #cgo 指令针对不同平台定义的宏，只有在对应平台的宏被定义之后才会构建对应的代码。
但是通过 #cgo 指令定义宏有个限制，它只能是基于 Go 语言支持的 windows、darwin 和 linux 等已经支持的操作系统。
如果我们希望定义一个 DEBUG 标志的宏，#cgo 指令就无能为力了。而 Go 语言提供的 build tag 条件编译特性则可以简单做到。

比如下面的源文件只有在设置 debug 构建标志时才会被构建：
// +build debug

package main

var buildMode = "debug"

可以用以下命令构建：
go build -tags="debug"
go build -tags="windows debug"

我们可以通过 -tags 命令行参数同时指定多个 build 标志，它们之间用空格分隔。

当有多个 build tag 时，我们将多个标志通过逻辑操作的规则来组合使用。比如以下的构建标志表示只有在”linux/386“或”darwin 平台下非 cgo 环境 “才进行构建。
// +build linux,386 darwin,!cgo

其中 linux,386 中 linux 和 386 用逗号链接表示 AND 的意思；而 linux,386 和 darwin,!cgo 之间通过空白分割来表示 OR 的意思。
*/

func main() {

}
