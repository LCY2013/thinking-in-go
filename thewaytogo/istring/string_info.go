/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-12-03
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
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	str1 = `This is a raw string \n`
)

func main() {
	fmt.Println(str1)
	str := "Beginning of the string " +
		"second part of the string"
	fmt.Println(str)

	s := "hel" + "lo,"
	s += "world!"
	fmt.Println(s) //输出 “hello, world!”

	runeInfo()

	stringsPrefixInfo()

	stringsContains()

	stringsIndex()

	stringsReplace()

	stringsCount()

	stringsRepeat()

	StringsToLowerToUpper()

	StringsTrimSpace()

	stringsFields()

	stringsJoin()

	stringsReader()

	strconvInfo()
}

// 字符串与其它类型的转换
// strconv.Itoa(i int) string 返回数字 i 所表示的字符串类型的十进制数。
// strconv.FormatFloat(f float64, fmt byte, prec int, bitSize int) string 将 64 位浮点型的数字转换为字符串，其中 fmt 表示格式（其值可以是 'b'、'e'、'f' 或 'g'），prec 表示精度，bitSize 则使用 32 表示 float32，用 64 表示 float64。
//
// strconv.Atoi(s string) (i int, err error) 将字符串转换为 int 型。
// strconv.ParseFloat(s string, bitSize int) (f float64, err error) 将字符串转换为 float64 型。
func strconvInfo() {
	fmt.Printf("%s - %T\n", strconv.Itoa(1), strconv.Itoa(1))

	var orig string = "666"
	var an int
	var newS string

	fmt.Printf("The size of ints is: %d\n", strconv.IntSize)

	an, _ = strconv.Atoi(orig)
	fmt.Printf("The integer is: %d\n", an)
	an = an + 5
	newS = strconv.Itoa(an)
	fmt.Printf("The new string is: %s\n", newS)
}

// 从字符串中读取内容
func stringsReader() {
	var reader *strings.Reader = strings.NewReader("The quick brown fox jumps over ths lazy dog")

	readByte, _ := reader.ReadByte()

	fmt.Println(readByte)

	readRune, size, _ := reader.ReadRune()

	fmt.Println(readRune, size)

}

// 拼接 slice 到字符串
func stringsJoin() {
	str := "The quick brown fox jumps over the lazy dog"
	sl := strings.Fields(str)
	fmt.Printf("Splitted in slice: %v\n", sl)
	for _, val := range sl {
		fmt.Printf("%s - ", val)
	}
	fmt.Println()
	str2 := "GO1|The ABC of Go|25"
	sl2 := strings.Split(str2, "|")
	fmt.Printf("Splitted in slice: %v\n", sl2)
	for _, val := range sl2 {
		fmt.Printf("%s - ", val)
	}
	fmt.Println()
	str3 := strings.Join(sl2, ";")
	fmt.Printf("sl2 joined by ;: %s\n", str3)
}

// 分割字符串
func stringsFields() {
	// strings.Fields(s) 将会利用 1 个或多个空白符号来作为动态长度的分隔符将字符串分割成若干小块
	var orig string = "Hey, how are you fufeng fufeng, i love this word!"
	fmt.Println(strings.Fields(orig)[:1])
	fmt.Println(strings.Split(orig, "fufeng"))
}

// 修剪字符串
// 剔除字符串开头和结尾
func StringsTrimSpace() {
	var orig string = " Hey, how are you George? "
	fmt.Println(strings.TrimSpace(orig))
	var origCut string = "Hey, how are you George? "
	fmt.Println(strings.Trim(origCut, "Hey"))
}

// 修改字符串大小写
func StringsToLowerToUpper() {
	var orig string = "Hey, how are you George?"
	var lower string
	var upper string

	fmt.Printf("The original string is: %s\n", orig)
	lower = strings.ToLower(orig)
	fmt.Printf("The lowercase string is: %s\n", lower)
	upper = strings.ToUpper(orig)
	fmt.Printf("The uppercase string is: %s\n", upper)
}

// 重复字符串
func stringsRepeat() {
	var origS string = "Hi there! "
	var newS string

	newS = strings.Repeat(origS, 3)
	fmt.Printf("The new repeated string is: %s\n", newS)
}

// 统计字符串出现次数
func stringsCount() {
	var str string = "Hello, how is it going, Hugo?"
	var manyG = "gggggggggg"

	fmt.Printf("Number of H's in %s is: ", str)
	fmt.Printf("%d\n", strings.Count(str, "H"))

	fmt.Printf("Number of double g's in %s is: ", manyG)
	fmt.Printf("%d\n", strings.Count(manyG, "gg"))
}

// 字符串替换
func stringsReplace() {
	var str string = "This is an example of a string"
	fmt.Println(strings.ReplaceAll(str, "is", "are"))
}

// 判断子字符串或字符在父字符串中出现的位置（索引）
func stringsIndex() {
	var str string = "This is an example of a string"
	fmt.Printf("T/F? Does the string \"%s\" have index %s? ", str, "h")
	fmt.Printf("%v\n", strings.Index(str, "h"))
	fmt.Printf("T/F? Does the string \"%s\" have lastIndex %s? ", str, "i")
	fmt.Printf("%v\n", strings.LastIndex(str, "i"))
	fmt.Printf("T/F? Does the string \"%s\" have indexRune %s? ", str, "i")
	r := []rune("i")
	fmt.Printf("%v\n", strings.IndexRune(str, r[0]))
}

// 字符串包含关系
// Contains 判断字符串 s 是否包含 substr
func stringsContains() {
	var str string = "This is an example of a string"
	fmt.Printf("T/F? Does the string \"%s\" have contain %s? ", str, "Th")
	fmt.Printf("%t\n", strings.Contains(str, "Th"))
}

// strings prefix 使用
// HasPrefix 判断字符串 s 是否以 prefix 开头：
// HasSuffix 判断字符串 s 是否以 suffix 结尾：
func stringsPrefixInfo() {
	var str string = "This is an example of a string"
	fmt.Printf("T/F? Does the string \"%s\" have prefix %s? ", str, "Th")
	fmt.Printf("%t\n", strings.HasPrefix(str, "Th"))
}

// rune 计算字符长度
func runeInfo() {
	// count number of characters:
	str1 := "asSASA ddd dsjkdsjs dk"
	fmt.Printf("The number of bytes in string str1 is %d\n", len(str1))
	fmt.Printf("The number of characters in string str1 is %d\n", utf8.RuneCountInString(str1))
	str2 := "asSASA ddd dsjkdsjsこん dk"
	fmt.Printf("The number of bytes in string str2 is %d\n", len(str2))
	fmt.Printf("The number of characters in string str2 is %d\n", utf8.RuneCountInString(str2))
}
