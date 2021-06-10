/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-06-10
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
	"sort"
)

/*
Conversions

Sequence 的 String 方法重新实现了 Sprint 为切片实现的功能。
若我们在调用 Sprint 之前将 Sequence 转换为纯粹的 []int，就能共享已实现的功能。
*/

// Sequence 自定义排序接口实现
/*
Interfaces and other types

Go 中的接口为指定对象的行为提供了一种方法:如果某样东西可以完成这个， 那么它就可以用在这里。
通过实现 String 方法，可以自定义打印函数，而通过 Write 方法，Fprintf 则能对任何对象产生输出。
在 Go 代码中， 仅包含一两种方法的接口很常见，且其名称通常来自于实现它的方法， 如 io.Writer 就是实现了 Write 的一类对象。

每种类型都能实现多个接口。例如一个实现了 sort.Interface 接口的集合就可通过 sort 包中的 例程进行排序。
该接口包括 Len()、Less(i, j int) bool 以及 Swap(i, j int)，另外，该集合仍然可 以有一个自定义的格式化器。
以下特意构建的例子 Sequence 就同时满足这两种情况。
*/
type Sequence []int

// Len
// Methods required by sort.Interface
// sort.Interface 所需的方法
func (s Sequence) Len() int {
	return len(s)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s Sequence) Less(i, j int) bool {
	return s[i] < s[j]
}

// Swap swaps the elements with indexes i and j.
func (s Sequence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Method for printing - sorts the elements before printing
// 用于在排序之前打印数组元素
/*
现在，不必让 Sequence 实现多个接口(排序和打印)，
可通过将数据条目转换为多种类型(Sequence、sort.IntSlice 和 []int)来使用相应的功能，每次转换都完成一部分工作。
这在实践中虽然有些不同寻常，但往往却很有效。
*/
func (s Sequence) String() string {
	sort.Sort(s)
	return fmt.Sprint([]int(s))
}

// Stringer
/*
Interface conversions and type assertions
接口转换与类型断言

类型选择 是类型转换的一种形式:它接受一个接口，在选择(switch)中根据其判断选择对应的情况(case)， 并在某种意义上将其转换为该种类型。
以下代码为 fmt.Printf 通过类型 选择将值转换为字符串的简化版。
若它已经为字符串，需要该接口中实际的字符串值; 若它有 String 方法，则需要调用该方法所得的结果。
*/
type Stringer interface {
	String() string
}

// StringerPrint
// 第一种情况获取具体的值，第二种将该接口转换为另一个接口。这种方式对于混合类型来说非常完美。
func StringerPrint() string {
	// 调用者提供的值
	var value interface{}
	switch str := value.(type) {
	case string:
		return str
	case Stringer:
		return str.String()
	}
	return ""
}

/*
若只关心一种类型呢?
若知道该值拥有一个 string 而想要提取它呢?
只需一种情况 的类型选择就行，但它需要类型断言。
类型断言接受一个接口值，并从中提取指定的明确类型的值。
其语法借鉴自类型选择开头的子句，但它需要一个明确的类型， 而非 type 关键字:
value.(type)

但若它所转换的值中不包含字符串，该程序就会以运行时错误崩溃。
为避免这种情况， 需使 用 “逗号, ok” 惯用测试它能安全地判断该值是否为字符串:
*/
func valueConverted(value interface{}) {
	str, ok := value.(string)
	if ok {
		fmt.Printf("String value is : %q\n", str)
	} else {
		fmt.Printf("value is not string\n")
	}
}

/*
若类型断言失败，str 将继续存在且为字符串类型，但它将拥有零值，即空字符串。

作为对能量的说明，这里有个 if-else 语句，它等价于本节开头的类型选择。
*/
func assertStr(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	} else if str, ok := value.(Stringer); ok {
		return str.String()
	}
	return ""
}

func main() {
	valueConverted([]int{})
	valueConverted("hello fufeng")
}
