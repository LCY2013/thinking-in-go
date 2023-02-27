package main

/*
C 代码的模块化

在编程过程中，抽象和模块化是将复杂问题简化的通用手段。
当代码语句变多时，我们可以将相似的代码封装到一个个函数中；
当程序中的函数变多时，我们将函数拆分到不同的文件或模块中。
而模块化编程的核心是面向程序接口编程（这里的接口并不是 Go 语言的 interface，而是 API 的概念）。

在前面的例子中，我们可以抽象一个名为 hello 的模块，模块的全部接口函数都在 hello.h 头文件定义：
// hello.h
void SayHello(const char* s);

其中只有一个 SayHello 函数的声明。
但是作为 hello 模块的用户来说，就可以放心地使用 SayHello 函数，而无需关心函数的具体实现。
而作为 SayHello 函数的实现者来说，函数的实现只要满足头文件中函数的声明的规范即可。
下面是 SayHello 函数的 C 语言实现，对应 hello.c 文件：
// hello.c

#include "hello.h"
#include <stdio.h>

	void SayHello(const char* s) {
		puts(s);
	}

在 hello.c 文件的开头，实现者通过 #include "hello.h" 语句包含 SayHello 函数的声明，
这样可以保证函数的实现满足模块对外公开的接口。

接口文件 hello.h 是 hello 模块的实现者和使用者共同的约定，但是该约定并没有要求必须使用 C 语言来实现 SayHello 函数。
我们也可以用 C++ 语言来重新实现这个 C 语言函数：
// hello.cpp

#include <iostream>

	extern "C" {
		#include "hello.h"
	}

	void SayHello(const char* s) {
		std::cout << s;
	}

在 C++ 版本的 SayHello 函数实现中，我们通过 C++ 特有的 std::cout 输出流输出字符串。
不过为了保证 C++ 语言实现的 SayHello 函数满足 C 语言头文件 hello.h 定义的函数规范，
我们需要通过 extern "C" 语句指示该函数的链接符号遵循 C 语言的规则。

在采用面向 C 语言 API 接口编程之后，我们彻底解放了模块实现者的语言枷锁：
实现者可以用任何编程语言实现模块，只要最终满足公开的 API 约定即可。
我们可以用 C 语言实现 SayHello 函数，也可以使用更复杂的 C++ 语言来实现 SayHello 函数，
当然我们也可以用汇编语言甚至 Go 语言来重新实现 SayHello 函数。
*/
func main() {

}
