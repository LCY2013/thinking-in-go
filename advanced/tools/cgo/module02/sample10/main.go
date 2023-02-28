package main

//#include <stdio.h>
import "C"
import "unsafe"

/*
C++ 类包装

CGO 是 C 语言和 Go 语言之间的桥梁，原则上无法直接支持 C++ 的类。
CGO 不支持 C++ 语法的根本原因是 C++ 至今为止还没有一个二进制接口规范 (ABI)。
一个 C++ 类的构造函数在编译为目标文件时如何生成链接符号名称、方法在不同平台甚至是 C++ 的不同版本之间都是不一样的。
但是 C++ 是兼容 C 语言，所以我们可以通过增加一组 C 语言函数接口作为 C++ 类和 CGO 之间的桥梁，这样就可以间接地实现 C++ 和 Go 之间的互联。
当然，因为 CGO 只支持 C 语言中值类型的数据类型，所以我们是无法直接使用 C++ 的引用参数等特性的。

C++ 类到 Go 语言对象

实现 C++ 类到 Go 语言对象的包装需要经过以下几个步骤：
首先是用纯 C 函数接口包装该 C++ 类；
其次是通过 CGO 将纯 C 函数接口映射到 Go 函数；
最后是做一个 Go 包装对象，将 C++ 类到方法用 Go 对象的方法实现。

准备一个 C++ 类

为了演示简单，我们基于 std::string 做一个最简单的缓存类 MyBuffer。
除了构造函数和析构函数之外，只有两个成员函数分别是返回底层的数据指针和缓存的大小。
因为是二进制缓存，所以我们可以在里面中放置任意数据。

#include <string>

struct MyBuffer {
    std::string* s_;

    MyBuffer(int size) {
        this->s_ = new std::string(size, char('\0'));
    }

    ~MyBuffer() {
        delete this->s_;
    }

    int size() const {
        return this->s_->size();
    }

    char* Data() {
        return (char*)this->s_->data();
    }
}

我们在构造函数中指定缓存的大小并分配空间，在使用完之后通过析构函数释放内部分配的内存空间。下面是简单的使用方式：

int main() {
    auto pBuf = new MyBuffer(1024);

    auto data = pBuf->Data();
    auto size = pBuf->Size();

    delete pBuf;
}

为了方便向 C 语言接口过渡，在此处我们故意没有定义 C++ 的拷贝构造函数。
我们必须以 new 和 delete 来分配和释放缓存对象，而不能以值风格的方式来使用。

用纯 C 函数接口封装 C++ 类

如果要将上面的 C++ 类用 C 语言函数接口封装，我们可以从使用方式入手。我们可以将 new 和 delete 映射为 C 语言函数，将对象的方法也映射为 C 语言函数。

在 C 语言中我们期望 MyBuffer 类可以这样使用：
int main() {
    MyBuffer* pBuf = NewMyBuffer(1024);

    char* data = MyBuffer_Data(pBuf);
    auto size = MyBuffer_Size(pBuf);

    DeleteMyBuffer(pBuf);
}

先从 C 语言接口用户的角度思考需要什么样的接口，然后创建 my_buffer_capi.h 头文件接口规范：
// my_buffer_capi.h
typedef struct MyBuffer_T MyBuffer_T;

MyBuffer_T* NewMyBuffer(int size);
void DeleteMyBuffer(MyBuffer_T* p);

char* MyBuffer_Data(MyBuffer_T* p);
int MyBuffer_Size(MyBuffer_T* p);

然后就可以基于 C++ 的 MyBuffer 类定义这些 C 语言包装函数。我们创建对应的 my_buffer_capi.cc 文件如下：
// my_buffer_capi.cc

#include "./my_buffer.h"

extern "C" {
    #include "./my_buffer_capi.h"
}

struct MyBuffer_T: MyBuffer {
    MyBuffer_T(int size): MyBuffer(size) {}
    ~MyBuffer_T() {}
};

MyBuffer_T* NewMyBuffer(int size) {
    auto p = new MyBuffer_T(size);
    return p;
}
void DeleteMyBuffer(MyBuffer_T* p) {
    delete p;
}

char* MyBuffer_Data(MyBuffer_T* p) {
    return p->Data();
}
int MyBuffer_Size(MyBuffer_T* p) {
    return p->Size();
}

因为头文件 my_buffer_capi.h 是用于 CGO，必须是采用 C 语言规范的名字修饰规则。
在 C++ 源文件包含时需要用 extern "C" 语句说明。
另外 MyBuffer_T 的实现只是从 MyBuffer 继承的类，这样可以简化包装代码的实现。
同时和 CGO 通信时必须通过 MyBuffer_T 指针，我们无法将具体的实现暴露给 CGO，因为实现中包含了 C++ 特有的语法，CGO 无法识别 C++ 特性。

将 C++ 类包装为纯 C 接口之后，下一步的工作就是将 C 函数转为 Go 函数。

将纯 C 接口函数转为 Go 函数

将纯 C 函数包装为对应的 Go 函数的过程比较简单。需要注意的是，因为我们的包中包含 C++11 的语法，因此需要通过 #cgo CXXFLAGS: -std=c++11 打开 C++11 的选项。

// my_buffer_capi.go

package main

/*
#cgo CXXFLAGS: -std=c++11

#include "my_buffer_capi.h"
*
import "C"

type cgo_MyBuffer_T C.MyBuffer_T

func cgo_NewMyBuffer(size int) *cgo_MyBuffer_T {
	p := C.NewMyBuffer(C.int(size))
	return (*cgo_MyBuffer_T)(p)
}

func cgo_DeleteMyBuffer(p *cgo_MyBuffer_T) {
	C.DeleteMyBuffer((*C.MyBuffer_T)(p))
}

func cgo_MyBuffer_Data(p *cgo_MyBuffer_T) *C.char {
	return C.MyBuffer_Data((*C.MyBuffer_T)(p))
}

func cgo_MyBuffer_Size(p *cgo_MyBuffer_T) C.int {
	return C.MyBuffer_Size((*C.MyBuffer_T)(p))
}

为了区分，我们在 Go 中的每个类型和函数名称前面增加了 cgo_ 前缀，比如 cgo_MyBuffer_T 是对应 C 中的 MyBuffer_T 类型。

为了处理简单，在包装纯 C 函数到 Go 函数时，除了 cgo_MyBuffer_T 类型外，对输入参数和返回值的基础类型，我们依然是用的 C 语言的类型。

包装为 Go 对象

在将纯 C 接口包装为 Go 函数之后，我们就可以很容易地基于包装的 Go 函数构造出 Go 对象来。
因为 cgo_MyBuffer_T 是从 C 语言空间导入的类型，它无法定义自己的方法，因此我们构造了一个新的 MyBuffer 类型，里面的成员持有 cgo_MyBuffer_T 指向的 C 语言缓存对象。

// my_buffer.go

package main

import "unsafe"

type MyBuffer struct {
    cptr *cgo_MyBuffer_T
}

func NewMyBuffer(size int) *MyBuffer {
    return &MyBuffer{
        cptr: cgo_NewMyBuffer(size),
    }
}

func (p *MyBuffer) Delete() {
    cgo_DeleteMyBuffer(p.cptr)
}

func (p *MyBuffer) Data() []byte {
    data := cgo_MyBuffer_Data(p.cptr)
    size := cgo_MyBuffer_Size(p.cptr)
    return ((*[1 << 31]byte)(unsafe.Pointer(data)))[0:int(size):int(size)]
}

同时，因为 Go 语言的切片本身含有长度信息，我们将 cgo_MyBuffer_Data 和 cgo_MyBuffer_Size 两个函数合并为 MyBuffer.Data 方法，它返回一个对应底层 C 语言缓存空间的切片。

现在我们就可以很容易在 Go 语言中使用包装后的缓存对象了（底层是基于 C++ 的 std::string 实现）：

package main

//#include <stdio.h>
import "C"
import "unsafe"

func main() {
    buf := NewMyBuffer(1024)
    defer buf.Delete()

    copy(buf.Data(), []byte("hello\x00"))
    C.puts((*C.char)(unsafe.Pointer(&(buf.Data()[0]))))
}

例子中，我们创建了一个 1024 字节大小的缓存，然后通过 copy 函数向缓存填充了一个字符串。
为了方便 C 语言字符串函数处理，我们在填充字符串的默认用'\0'表示字符串结束。最后我们直接获取缓存的底层数据指针，用 C 语言的 puts 函数打印缓存的内容。
*/

func main() {
	buf := NewMyBuffer(1024)
	defer buf.Delete()

	copy(buf.Data(), []byte("hello\x00"))
	C.puts((*C.char)(unsafe.Pointer(&(buf.Data()[0]))))
}
