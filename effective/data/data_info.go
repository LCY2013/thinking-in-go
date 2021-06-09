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
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
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

type SyncedBuffer struct {
	lock   sync.Mutex
	buffer bytes.Buffer
}

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
	fmt.Printf("%#v\n", timeZone)

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

	// p := new(SyncedBuffer)
	// var v SyncedBuffer

	//a := [...]string {Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
	//s := []string {Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
	//m := map[int]string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}

	// make
	intArrays := make([]int, 0, 10)
	fmt.Println(intArrays)

	// 分配切片结构;*p == nil;基本没用
	// var p *[]int = new([]int)

	// 切片 v 现在引用了一个具有 100 个 int 元素的新数组
	// var v []int = make([]int, 10)

	// 没必要的复杂:
	// var p *[]int = new([]int)
	// *p = make([]int, 100, 100)

	// 习惯用法
	// v := make([]int, 10)

	// ===== 注意  =====
	// 请记住，make 只适用于映射、切片和信道且不返回指针。若要获得明确的指针， 请使用 new 分配内存。

	//数组为值的属性很有用，但代价高昂;若你想要 C 那样的行为和效率，你可以传递一个指向 该数组的指针。
	array := [...]float64{1.1, 2.2, 3.3}
	f := Sum(&array)
	fmt.Println("array: ", f)
	//但这并不是 Go 的习惯用法，切片才是。

	//切片通过对数组进行封装，为数据序列提供了更通用、强大而方便的接口。
	//除了矩阵变换这 类需要明确维度的情况外，Go 中的大部分数组编程都是通过切片来完成的。

	//切片保存了对底层数组的引用，若你将某个切片赋予另一个切片，它们会引用同一个数组。
	//若某个函数将一个切片作为参数传入，则它对该切片元素的修改对调用者而言同样可见，这可以理解为传递了底层数组的指针。
	//因此，Read 函数可接受一个切片实参 而非一个指针和 一个计数;切片的长度决定了可读取数据的上限。
	//以下为 os 包中 File 类型的 Read 方法签 名:
	//  func (file *File) Read(buf []byte) (n int, err error)
	// 该方法返回读取的字节数和一个错误值(若有的话)。若要从更大的缓冲区 b 中读取前 32 个 字节，只需对其进行切片即可。
	//  n, err := f.Read(buf[0:32])

	//这种切片的方法常用且高效。若不谈效率，以下片段同样能读取该缓冲区的前 32 个字节。
	/*
		var n int
		var err error
		for i := 0; i < 32; i++ {
		    nbytes, e := f.Read(buf[i:i+1])  // Read one byte.
		    if nbytes == 0 || e != nil {
				err = e
				break
			}
			n += nbytes
		}
	*/

	//只要切片不超出底层数组的限制，它的长度就是可变的，只需将它赋予其自身的切片即可。
	//切片的容量可通过内建函数 cap 获得，它将给出该切片可取得的最大长度。
	//以下是将数据追 加到切片的函数。若数据超出其容量，则会重新分配该切片。返回值即为所得的切片。
	//该函 数中所使用的 len 和 cap 在应用于 nil 切片时是合法的，它会返回 0.
	/*
		func Append(slice, data []byte) []byte {
		    l := len(slice)
		    if l + len(data) > cap(slice) {  // reallocate
		        // Allocate double what's needed, for future growth.
		        newSlice := make([]byte, (l+len(data))*2)
		        // The copy function is predeclared and works for any slice type.
		        copy(newSlice, slice)
		        slice = newSlice
		    }
		    slice = slice[0:l+len(data)]
		    for i, c := range data {
		        slice[l+i] = c
		    }
		    return slice
		}
	*/

	// 最终我们必须返回切片，因为尽管 Append 可修改 slice 的元素，但切片自身(其运行时数据 结构包含指针、长度和容量)是通过值传递的。

	// 向切片追加东西的想法非常有用，因此有专门的内建函数 append。 要理解该函数的设计，我 们还需要一些额外的信息。

	// Two-dimensional slices
	// Go 的数组和切片都是一维的。要创建等价的二维数组或切片，就必须定义一个数组的数组， 或切片的切片，就像这样:
	// type Transform [3][3]float64  // A 3x3 array, really an array of arrays.
	// type LinesOfText [][]byte     // A slice of byte slices.

	// 由于切片长度是可变的，因此其内部可能拥有多个不同长度的切片。在下面例子中，这是种常见的情况:每行都有其自己的长度。
	/*
		text := LinesOfText{
		    []byte("Now is the time"),
		    []byte("for all good gophers"),
		    []byte("to bring some fun to the party."),
		}
	*/

	//有时必须分配一个二维数组，例如在处理像素的扫描行时，这种情况就会发生。
	//有两种方式来达到这个目的。一种就是独立地分配每一个切片;而另一种就是只分配一个数组， 将 各个切片都指向它。
	//采用哪种方式取决于你的应用。若切片会增长或收缩， 就应该通过独立 分配来避免覆盖下一行;
	//若不会，用单次分配来构造对象会更加高效。 以下是这两种方法的 大概代码，仅供参考。首先是一次一行的:
	/*
		// Allocate the top-level slice.
		picture := make([][]uint8, YSize) // One row per unit of y.
		// Loop over the rows, allocating the slice for each row.
		for i := range picture {
		    picture[i] = make([]uint8, XSize)
		}
	*/
	/*
		// 分配顶层切片，和前面一样。
		picture := make([][]uint8, YSize) // 每 y 个单元一行。
		// 分配一个大的切片来保存所有像素
		pixels := make([]uint8, XSize*YSize) // 拥有类型 []uint8，尽管图片是 [][]uint8.
		// 遍历行，从剩余像素切片的前面切出每行来。
		for i := range picture {
		    picture[i], pixels = pixels[:XSize], pixels[XSize:]
		}
	*/

	// Maps
	// 映射是方便而强大的内建数据结构，它可以关联不同类型的值。
	// 其键可以是任何相等性操作符支持的类型， 如整数、浮点数、复数、字符串、指针、接口(只要其动态类型支持相等性 判断)、结构以及数组。
	// 切片不能用作映射键，因为它们的相等性还未定义。与切片一样， 映射也是引用类型。
	// 若将映射传入函数中，并更改了该映射的内容，则此修改对调用者同样 可见。

	// 映射可使用一般的复合字面语法进行构建，其键 - 值对使用逗号分隔，因此可在初始化时很容 易地构建它们。
	var timeZone = map[string]int{
		"UTC": 0 * 60 * 60,
		"EST": -5 * 60 * 60,
		"CST": -6 * 60 * 60,
		"MST": -7 * 60 * 60,
		"PST": -8 * 60 * 60,
	}
	fmt.Println(timeZone)

	// 赋值和获取映射值的语法类似于数组，不同的是映射的索引不必为整数。
	// offset := timeZone["EST"]

	//若试图通过映射中不存在的键来取值，就会返回与该映射中项的类型对应的零值。
	//例如，若 某个映射包含整数，当查找一个不存在的键时会返回 0。
	//集合可实现成一个值类型为 bool 的 映射。将该映射中的项置为 true 可将该值放入集合中，此后通过简单的索引操作即可判断是 否存在。
	attended := map[string]bool{
		"Ann": true,
		"Joe": true,
		"Mac": true,
	}
	if attended["Mac"] { // will be false if person is not in the map
		fmt.Println("Mac", "was at the meeting")
	}

	//有时你需要区分某项是不存在还是其值为零值。如对于一个值本应为零的 "UTC" 条目，也可 能是由于不存在该项而得到零值。你可以使用多重赋值的形式来分辨这种情况。
	// var seconds int
	// var ok bool
	// seconds , ok = timeZone["UTC"]

	//显然，我们可称之为 “逗号 ok” 惯用法。
	//在下面的例子中，若 tz 存在， seconds 就会被赋予 适当的值，且 ok 会被置为 true;
	//若不存在，seconds 则会被置为零，而 ok 会被置为 false。
	offset("UTC")

	// 若仅需判断映射中是否存在某项而不关心实际的值，可使用 空白标识符 ( _ )来代替该值的 一般变量。
	//_,present := timeZone["UTC"]

	// 要删除映射中的某项，可使用内建函数 delete，它以映射及要被删除的键为实参。 即便对应 的键不在该映射中，此操作也是安全的。
	delete(timeZone, "PDT")

	// Printing
	// Go 采用的格式化打印风格和 C 的 printf 族类似，但却更加丰富而通用。
	// 这些函数位于 fmt 包 中，且函数名首字母均为大写:如 fmt.Printf、fmt.Fprintf，fmt.Sprintf 等。
	// 字符串函数 (Sprintf 等)会返回一个字符串，而非填充给定的缓冲区。

	//无需提供一个格式字符串。
	//每个 Printf、Fprintf 和 Sprintf 都分别对应另外的函数，如 Print 与 Println。
	//这些函数并不接受格式字符串，而是为每个实参生成一种默认格式。
	//Println 系列 的函数还会在实参中插入空格，并在输出时追加一个换行符，而 Print 版本仅在操作数两侧都 没有字符串时才添加空白。
	//以下示例中各行产生的输出都是一样的。
	fmt.Printf("Hello %d\n", 23)
	_, _ = fmt.Fprint(os.Stdout, "Hello ", 23, "\n")
	fmt.Println("Hello", 23)
	fmt.Println(fmt.Sprint("Hello ", 23))

	//fmt.Fprint 一类的格式化打印函数可接受任何实现了 io.Writer 接口的对象作为第一个实参;
	//变 量 os.Stdout 与 os.Stderr 都是人们熟知的例子。

	//从这里开始，就与 C 有些不同了。首先，像%d 这样的数值格式并不接受表示符号或大小的标记， 打印例程会根据实参的类型来决定这些属性。
	fmt.Printf("%d %x; %d %x\n", x, x, int64(x), int64(x))

	//若你只想要默认的转换，如使用十进制的整数，你可以使用通用的格式 %v(对应 “值”);
	//其 结果与 Print 和 Println 的输出完全相同。
	//此外，这种格式还能打印任意值，甚至包括数组、 结构体和映射。 以下是打印上一节中定义的时区映射的语句。
	fmt.Printf("%v\n", timeZone)

	//当然，映射中的键可能按任意顺序输出。
	//当打印结构体时，改进的格式 %+v 会为结构体的每 个字段添上字段名，
	//而另一种格式 %#v 将完全按照 Go 的语法打印值。
	fmt.Printf("%v\n", t)
	// 如果你需要像指向 T 的指针那样打印类型 T 的值， String 的接收者就必须是值类型的;上 面的例子中接收者是一个指针， 因为这对结构来说更高效而通用。
	// 我们的 String 方法也可调用 Sprintf， 因为打印例程可以完全重入并按这种方式封装。不过要 理解这种方式，还有一个重要的细节: 请勿通过调用 Sprintf 来构造 String 方法，因为它会无 限递归你的的 String 方法。

}

//Sum
/*
以下为数组在 Go 和 C 中的主要区别。
在 Go 中，数组是值。将一个数组赋予另一个数组会复制其所有元素。
特别地，若将某个数组传入某个函数，它将接收到该数组的一份副本而非指针。
数组的大小是其类型的一部分。类型 [10]int 和 [20]int 是不同的。
*/
func Sum(a *[3]float64) (sum float64) {
	for _, v := range *a {
		sum += v
	}
	return
}

var timeZone = map[string]int{
	"UTC": 0 * 60 * 60,
	"EST": -5 * 60 * 60,
	"CST": -6 * 60 * 60,
	"MST": -7 * 60 * 60,
	"PST": -8 * 60 * 60,
}

//offset
/*
显然，我们可称之为 “逗号 ok” 惯用法。
在下面的例子中，若 tz 存在， seconds 就会被赋予 适当的值，且 ok 会被置为 true;
若不存在，seconds 则会被置为零，而 ok 会被置为 false。
*/
func offset(tz string) int {
	if seconds, ok := timeZone[tz]; ok {
		return seconds
	}
	log.Println("unknown time zone:", tz)
	return 0
}

type T struct {
	a int
	b float64
	c string
}

func (t *T) String() string {
	return fmt.Sprintf("%d/%g/%q", t.a, t.b, t.c)
}

// MyString 下列演示错误无限递归,Sprintf会无限调用某个类型的String方法
type MyString string

/*func (ms MyString) String() string {
	return fmt.Sprintf("MyString=%s",ms)
}*/
// 要解决这个问题也很简单:将该实参转换为基本的字符串类型，它没有这个方法。
func (ms MyString) String() string {
	return fmt.Sprintf("MyString=%s", string(ms))
}

// Min 获取最大的值
func Min(array ...int) int {
	// 最大的int
	min := int(^uint(0) >> 1)
	for _, i := range array {
		if i < min {
			min = i
		}
	}
	return min
}
