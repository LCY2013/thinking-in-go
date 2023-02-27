package main

/*
类型转换

最初 CGO 是为了达到方便从 Go 语言函数调用 C 语言函数（用 C 语言实现 Go 语言声明的函数）以复用 C 语言资源这一目的而出现的（因为 C 语言还会涉及回调函数，自然也会涉及到从 C 语言函数调用 Go 语言函数（用 Go 语言实现 C 语言声明的函数））。
现在，它已经演变为 C 语言和 Go 语言双向通讯的桥梁。要想利用好 CGO 特性，自然需要了解此二语言类型之间的转换规则。

数值类型

在 Go 语言中访问 C 语言的符号时，一般是通过虚拟的 “C” 包访问，比如 C.int 对应 C 语言的 int 类型。
有些 C 语言的类型是由多个关键字组成，但通过虚拟的 “C” 包访问 C 语言类型时名称部分不能有空格字符，比如 unsigned int 不能直接通过 C.unsigned int 访问。
因此 CGO 为 C 语言的基础数值类型都提供了相应转换规则，比如 C.uint 对应 C 语言的 unsigned int。

Go 语言中数值类型和 C 语言数据类型基本上是相似的，以下是它们的对应关系表 2-1 所示。

C 语言类型	CGO 类型	Go 语言类型
char	C.char	byte
singed char	C.schar	int8
unsigned char	C.uchar	uint8
short	C.short	int16
unsigned short	C.ushort	uint16
int	C.int	int32
unsigned int	C.uint	uint32
long	C.long	int32
unsigned long	C.ulong	uint32
long long int	C.longlong	int64
unsigned long long int	C.ulonglong	uint64
float	C.float	float32
double	C.double	float64
size_t	C.size_t	uint

需要注意的是，虽然在 C 语言中 int、short 等类型没有明确定义内存大小，但是在 CGO 中它们的内存大小是确定的。
在 CGO 中，C 语言的 int 和 long 类型都是对应 4 个字节的内存大小，size_t 类型可以当作 Go 语言 uint 无符号整数类型对待。

CGO 中，虽然 C 语言的 int 固定为 4 字节的大小，但是 Go 语言自己的 int 和 uint 却在 32 位和 64 位系统下分别对应 4 个字节和 8 个字节大小。
如果需要在 C 语言中访问 Go 语言的 int 类型，可以通过 GoInt 类型访问，GoInt 类型在 CGO 工具生成的 _cgo_export.h 头文件中定义。
其实在 _cgo_export.h 头文件中，每个基本的 Go 数值类型都定义了对应的 C 语言类型，它们一般都是以单词 Go 为前缀。
下面是 64 位环境下，_cgo_export.h 头文件生成的 Go 数值类型的定义，其中 GoInt 和 GoUint 类型分别对应 GoInt64 和 GoUint64：

typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef float GoFloat32;
typedef double GoFloat64;

除了 GoInt 和 GoUint 之外，我们并不推荐直接访问 GoInt32、GoInt64 等类型。更好的做法是通过 C 语言的 C99 标准引入的 <stdint.h> 头文件。
为了提高 C 语言的可移植性，在 <stdint.h> 文件中，不但每个数值类型都提供了明确内存大小，而且和 Go 语言的类型命名更加一致。
Go 语言类型 <stdint.h> 头文件类型对比如表 2-2 所示。

C 语言类型	CGO 类型	Go 语言类型
int8_t	C.int8_t	int8
uint8_t	C.uint8_t	uint8
int16_t	C.int16_t	int16
uint16_t	C.uint16_t	uint16
int32_t	C.int32_t	int32
uint32_t	C.uint32_t	uint32
int64_t	C.int64_t	int64
uint64_t	C.uint64_t	uint64

前文说过，如果 C 语言的类型是由多个关键字组成，则无法通过虚拟的 “C” 包直接访问(比如 C 语言的 unsigned short 不能直接通过 C.unsigned short 访问)。
但是，在 <stdint.h> 中通过使用 C 语言的 typedef 关键字将 unsigned short 重新定义为 uint16_t 这样一个单词的类型后，
我们就可以通过 C.uint16_t 访问原来的 unsigned short 类型了。
对于比较复杂的 C 语言类型，推荐使用 typedef 关键字提供一个规则的类型命名，这样更利于在 CGO 中访问。
*/

func main() {

}
