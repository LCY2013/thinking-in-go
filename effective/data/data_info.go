/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-04-20
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package main

import (
	"fmt"
	"os"
)

/**
new

Go 提供了两种分配原语，即内建函数 new 和 make。 它们所做的事情不同，所应用的类型 也不同。它们可能会引起混淆，但规则却很简单。
让我们先来看看 new。这是个用来分配内 存的内建函数， 但与其它语言中的同名函数不同，它不会初始化内存，只会将内存置零。
也 就是说，new(T) 会为类型为 T 的新项分配已置零的内存空间， 并返回它的地址，也就是一个 类型为 *T 的值。
用 Go 的术语来说，它返回一个指针， 该指针指向新分配的，类型为 T 的 零值。

既然 new 返回的内存已置零，那么当你设计数据结构时， 每种类型的零值就不必进一步初始 化了，这意味着该数据结构的使用者只需用 new 创建一个新的对象就能正常工作。
例如， bytes.Buffer 的文档中提到 “零值的 Buffer 就是已准备就绪的缓冲区。
同样，sync.Mutex 并 没有显式的构造函数或 Init 方法， 而是零值的 sync.Mutex 就已经被定义为已解锁的互斥锁 了。

“零值属性” 可以带来各种好处。考虑以下类型声明。
type SyncedBuffer struct {
    lock    sync.Mutex
    buffer  bytes.Buffer
}

SyncedBuffer 类型的值也是在声明时就分配好内存就绪了。后续代码中， p 和 v 无需进一步 处理即可正确工作。

构造函数与复合字面

有时零值还不够好，这时就需要一个初始化构造函数，如来自 os 包中的这段代码所示。

func NewFile(fd int, name string) *File {
    if fd < 0 {
return nil }
    f := new(File)
    f.fd = fd
    f.name = name
    f.dirinfo = nil
    f.nepipe = 0
    return f
}

这里显得代码过于冗长。我们可通过复合字面来简化它， 该表达式在每次求值时都会创建新 的实例。

func NewFile(fd int, name string) *File {
    if fd < 0 {
return nil }
    f := File{fd, name, nil, 0}
return &f
}
*/

/**
make

make 分配

再回到内存分配上来。内建函数 make(T, args) 的目的不同于 new(T)。它只用于创建切片、 映射和信道，并返回类型为 T(而非 *T )的一个已初始化 (而非置零)的值。
出现这种用 差异的原因在于，这三种类型本质上为引用数据类型，它们在使用前必须初始化。
例如，切 片是一个具有三项内容的描述符，包含一个指向(数组内部)数据的指针、长度以及容量， 在这三项被初始化之前，该切片为 nil。
对于切片、映射和信道，make 用于初始化其内部的 数据结构并准备好将要使用的值。例如，
make([]int,10,100)

会分配一个具有 100 个 int 的数组空间，接着创建一个长度为 10， 容量为 100 并指向该数组 中前 10 个元素的切片结构。
与 此相反，new([]int) 会返回一个指向新分配的，已置零的切片结构， 即一个指向 nil 切片值的 指针。
*/

/**
下面的例子阐明了 new 和 make 之间的区别:

var p *[]int = new([]int) // 分配切片结构;*p == nil;基本没用
var v []int = make([]int, 100) // 切片 v 现在引用了一个具有 100 个 int 元素的新数组
// 没必要的复杂:
var p *[]int = new([]int) *p = make([]int, 100, 100)
// 习惯用法:
v := make([]int, 100)

请记住，make 只适用于映射、切片和信道且不返回指针。若要获得明确的指针， 请使用 new 分配内存。
*/

func main() {
	fmt.Printf("Hello %d\n", 23)
	fmt.Fprint(os.Stdout, "Hello ", 23, "\n")
	fmt.Println("Hello", 23)
	fmt.Println(fmt.Sprint("Hello ", 23))

	/*
		若你只想要默认的转换，如使用十进制的整数，你可以使用通用的格式 %v(对应 “值”);
		其 结果与 Print 和 Println 的输出完全相同。
		此外，这种格式还能打印任意值，甚至包括数组、 结构体和映射。 以下是打印上一节中定义的时区映射的语句。
	*/
	var x uint64 = 1<<64 - 1
	fmt.Printf("%d %x; %d %x\n", x, x, int64(x), int64(x))

	//fmt.Printf("%v\n",timeZone)

	/*
		当然，映射中的键可能按任意顺序输出。当打印结构体时，改进的格式 %+v 会为结构体的每 个字段添上字段名，而另一种格式 %#v 将完全按照 Go 的语法打印值。
	*/
	t := &T{7, -2.35, "abc\tdef"}
	fmt.Printf("%v\n", t)
	fmt.Printf("%+v\n", t)
	fmt.Printf("%#v\n", t)
	//fmt.Printf("%#v\n", timeZone)

	/*
		(请注意其中的 & 符号)当遇到 string 或 []byte 值时， 可使用 %q 产生带引号的字符串;而 格式 %#q 会尽可能使用反引号。
		(%q 格式也可用于整数和符文，它会产生一个带单引号的 符文常量。)
		此外，%x 还可用于字符串、字节数组以及整数，并生成一个很长的十六进制 字符串， 而带空格的格式(% x)还会在字节之间插入空格。
	*/

	// 另一种实用的格式是 %T，它会打印某个值的类型.

	//y := []int{1, 2, 3}
	//y = append(y, 4, 5, 6)
	//fmt.Println(y)

	k := []int{1, 2, 3}
	y := []int{4, 5, 6}
	k = append(k, y...)
	fmt.Println(k)

	/*
		常量
		Go 中的常量就是不变量。它们在编译时创建，即便它们可能是函数中定义的局部变量。
		常量 只能是数字、字符(符文)、字符串或布尔值。由于编译时的限制， 定义它们的表达式必须 也是可被编译器求值的常量表达式。
		例如 1<<3 就是一个常量表达式，而 math.Sin(math.Pi/4) 则不是，因为对 math.Sin 的函数调用在运行时才会发生。

		在 Go 中，枚举常量使用枚举器 iota 创建。由于 iota 可为表达式的一部分，而表达式可以被 隐式地重复，这样也就更容易构建复杂的值的集合了。


	*/
}

type T struct {
	a int
	b float64
	c string
}
