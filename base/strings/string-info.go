package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

/*
strings.Builder、strings.Join、fmt.Sprintf

两个视角来看待 Go 字符串的组成，一种是字节视角。Go 字符串是由一个可空的字节序列组成，字节的个数称为字符串的长度；另外一种是字符视角。Go 字符串是由一个可空的字符序列构成。Go 字符串中的每个字符都是一个 Unicode 字符。
Go 使用 rune 类型来表示一个 Unicode 字符的码点。为了传输和存储 Unicode 字符，Go 还使用了 UTF-8 编码方案，UTF-8 编码方案使用变长字节的编码方式，码点小的字符用较少的字节编码，码点大的字符用较多字节编码，这种编码方式兼容 ASCII 字符集，并且拥有很高的空间利用率。
Go 语言在运行时层面通过一个二元组结构（Data, Len）来表示一个 string 类型变量，其中 Data 是一个指向存储字符串数据内容区域的指针值，Len 是字符串的长度。因此，本质上，一个 string 变量仅仅是一个“描述符”，并不真正包含字符串数据。因此，我们即便直接将 string 类型变量作为函数参数，其传递的开销也是恒定的，不会随着字符串大小的变化而变化。
Go 为其原生支持的 string 类型提供了许多原生操作类型，在进行字符串操作时你要注意以下几点：
通过常规 for 迭代与 for range 迭代所得到的结果不同，常规 for 迭代采用的是字节视角；而 for range 迭代采用的是字符视角；
基于 +/+= 操作符的字符串连接是对开发者体验最好的字符串连接方式，但却不是性能最好的方式；
无论是字符串转切片，还是切片转字符串，都会有内存分配的开销，这缘于 Go 字符串数据内容不可变的性质。
*/
func dumpBytesArray(arr []byte) {
	fmt.Printf("[")
	for _, b := range arr {
		fmt.Printf("%c ", b)
	}
	fmt.Printf("]\n")
}

func main() {
	stringType()
	fmt.Println("----------------------------")
	iterator()
	fmt.Println("----------------------------")
	join()
	fmt.Println("----------------------------")
	compare()
	fmt.Println("----------------------------")
	replace()
}

// replace Go 支持字符串与字节切片、字符串与 rune 切片的双向转换，并且这种转换无需调用任何函数，只需使用显式类型转换
func replace() {
	var s string = "中国人"

	// string -> []rune
	rs := []rune(s)
	fmt.Printf("%x\n", rs) // [4e2d 56fd 4eba]

	// string -> []byte
	bs := []byte(s)
	fmt.Printf("%x\n", bs) // e4b8ade59bbde4baba

	// []rune -> string
	s1 := string(rs)
	fmt.Println(s1) // 中国人

	// []byte -> string
	s2 := string(bs)
	fmt.Println(s2) // 中国人
}

func stringType() {
	var s = "magic"
	// 显示将string类型地址转换成reflect.StringHeader
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	fmt.Printf("0x%x\n", hdr)
	// 取出data字段指向的数组的指针
	p := (*[5]byte)(unsafe.Pointer(hdr.Data))
	// 输出底层数组的内容
	dumpBytesArray((*p)[:])
}

func iterator() {
	var china = "中国人"
	fmt.Printf("0x%x\n", china[0]) // 0xe4：字符“中” utf-8编码的第一个字节

	for i := 0; i < len(china); i++ {
		fmt.Printf("index: %d, value: 0x%x\n", i, china[i])
	}

	for i, v := range china {
		fmt.Printf("index: %d, value: 0x%x\n", i, v)
	}
}

func join() {
	rob := "Rob Pike, "
	rob = rob + "Robert Griesemer, "
	rob += " Ken Thompson"
	fmt.Println(rob) // Rob Pike, Robert Griesemer, Ken Thompson
}

func compare() {
	// ==
	s1 := "世界和平"
	s2 := "世界" + "和平"
	fmt.Println(s1 == s2) // true
	// !=
	s1 = "Go"
	s2 = "C"
	fmt.Println(s1 != s2) // true
	// < and <=
	s1 = "12345"
	s2 = "23456"
	fmt.Println(s1 < s2)  // true
	fmt.Println(s1 <= s2) // true
	// > and >=
	s1 = "12345"
	s2 = "123"
	fmt.Println(s1 > s2)  // true
	fmt.Println(s1 >= s2) // true
}
