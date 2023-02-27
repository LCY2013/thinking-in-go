package main

/*
#cgo windows CFLAGS: -DCGO_OS_WINDOWS=1
#cgo darwin CFLAGS: -DCGO_OS_DARWIN=1
#cgo linux CFLAGS: -DCGO_OS_LINUX=1

#if defined(CGO_OS_WINDOWS)
    const char* os = "windows";
#elif defined(CGO_OS_DARWIN)
    const char* os = "darwin";
#elif defined(CGO_OS_LINUX)
    const char* os = "linux";
#else
#	error(unknown os)
#endif
*/
import "C"

/*
#cgo 语句
在 import "C" 语句前的注释中可以通过 #cgo 语句设置编译阶段和链接阶段的相关参数。
编译阶段的参数主要用于定义相关宏和指定头文件检索路径。
链接阶段的参数主要是指定库文件检索路径和要链接的库文件。

// #cgo CFLAGS: -DPNG_DEBUG=1 -I./include
// #cgo LDFLAGS: -L/usr/local/lib -lpng
// #include <png.h>
import "C"

上面的代码中，CFLAGS 部分，-D 部分定义了宏 PNG_DEBUG，值为 1；-I 定义了头文件包含的检索目录。
LDFLAGS 部分，-L 指定了链接时库文件检索目录，-l 指定了链接时需要链接 png 库。

因为 C/C++ 遗留的问题，C 头文件检索目录可以是相对目录，但是库文件检索目录则需要绝对路径。
在库文件的检索目录中可以通过 ${SRCDIR} 变量表示当前包目录的绝对路径：
// #cgo LDFLAGS: -L${SRCDIR}/libs -lfoo

上面的代码在链接时将被展开为：
// #cgo LDFLAGS: -L/go/src/foo/libs -lfoo

#cgo 语句主要影响 CFLAGS、CPPFLAGS、CXXFLAGS、FFLAGS 和 LDFLAGS 几个编译器环境变量。LDFLAGS 用于设置链接时的参数，除此之外的几个变量用于改变编译阶段的构建参数 (CFLAGS 用于针对 C 语言代码设置编译参数)。

对于在 cgo 环境混合使用 C 和 C++ 的用户来说，可能有三种不同的编译选项：其中 CFLAGS 对应 C 语言特有的编译选项、CXXFLAGS 对应是 C++ 特有的编译选项、CPPFLAGS 则对应 C 和 C++ 共有的编译选项。但是在链接阶段，C 和 C++ 的链接选项是通用的，因此这个时候已经不再有 C 和 C++ 语言的区别，它们的目标文件的类型是相同的。

#cgo 指令还支持条件选择，当满足某个操作系统或某个 CPU 架构类型时后面的编译或链接选项生效。比如下面是分别针对 windows 和非 windows 下平台的编译和链接选项：
// #cgo windows CFLAGS: -DX86=1
// #cgo !windows LDFLAGS: -lm

其中在 windows 平台下，编译前会预定义 X86 宏为 1；在非 windows 平台下，在链接阶段会要求链接 math 数学库。这种用法对于在不同平台下只有少数编译选项差异的场景比较适用。

如果在不同的系统下 cgo 对应着不同的 c 代码，我们可以先使用 #cgo 指令定义不同的 C 语言的宏，然后通过宏来区分不同的代码：
package main
/*
#cgo windows CFLAGS: -DCGO_OS_WINDOWS=1
#cgo darwin CFLAGS: -DCGO_OS_DARWIN=1
#cgo linux CFLAGS: -DCGO_OS_LINUX=1

#if defined(CGO_OS_WINDOWS)
    const char* os = "windows";
#elif defined(CGO_OS_DARWIN)
    const char* os = "darwin";
#elif defined(CGO_OS_LINUX)
    const char* os = "linux";
#else
#	error(unknown os)
#endif
*
import "C"

func main() {
	print(C.GoString(C.os))
}

这样我们就可以用 C 语言中常用的技术来处理不同平台之间的差异代码。

*/

func main() {
	print(C.GoString(C.os))
}
