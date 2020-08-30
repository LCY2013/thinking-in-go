/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2019 Ramostear.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-08-09
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package main

import (
	"fmt"
	"strings"
)

func main() {
	// testCase01()
	// testCase02()
	// testCase03()
	// testCase04()
	// testCase05()
	// testCase06()
	testCase07()
}

// copy 使用
func testCase07() {
	num := []int{1, 2, 3, 4, 5, 6}
	sli1 := num[3:]  // {4,5,6}
	sli2 := num[0:4] // {1,2,3,4} -> {4,5,6,4}
	_ = copy(sli2, sli1)
	fmt.Println("sli2 = ", sli2)
}

// 去掉slice中的重复字符
func testCase06() {
	str := []string{"fufeng", "magic", "", "", " ", "magic", "lcy"}
	fmt.Println("str = ", cast06NotSame(str))
	fmt.Println("len(str) = ", len(cast06NotSame(str)))
	fmt.Println("cap(str) = ", cap(cast06NotSame(str)))
}

func cast06NotSame(str []string) []string {
	var ret []string
	for _, v := range str {
		// 计算ret这里面是否存在同样的字符串
		i := 0
		for _, vRet := range ret {
			if v == vRet {
				break
			} else {
				i++
			}
		}
		if i == len(ret) {
			ret = append(ret, v)
		}
	}
	return ret
}

// 去掉slice中的空字符串
func testCase05() {
	str := []string{"fufeng", "magic", "", "", " ", "magic", "lcy"}
	fmt.Println("empty(str)=", case05NotEmpty(str))
	//str1 := []string{"fufeng","magic","magic","lcy"}
	//fmt.Println("emptyOldSlice(str)=", case05NotEmptyUseOldSlice(str1))
	fmt.Println("emptyOldSlice(str)=", case05NotEmptyUseOldSlice(str))
}

func case05NotEmptyUseOldSlice(str []string) []string {
	// 计算当前所到节点
	i := 0
	for _, v := range str {
		if v != "" && strings.TrimSpace(v) != "" {
			str[i] = v
			i++
		}
	}
	// 截取到当前计算到的节点位置
	return str[:i]
}

func case05NotEmpty(str []string) []string {
	// 这里这句会影响原有切片的值
	//ret := str[:0]
	ret := new([]string)
	for _, v := range str {
		// 不为空就追加
		if v != "" && strings.TrimSpace(v) != "" {
			*ret = append(*ret, v)
		}
	}
	return *ret
}

func testCase04() {
	// 切片添加
	sli := []int{1, 2, 3, 4, 5, 6, 7}
	sli = append(sli, 8)
	sli = append(sli, 9)
	sli = append(sli, 10)
	sli = append(sli, 11)
	fmt.Println(sli)
}

func testCase03() {
	// 自动推导赋值
	sli := [5]int{1, 3, 5, 7, 9}
	fmt.Println("sli = ", sli)
	fmt.Printf("%T\n", sli)
	fmt.Printf("%T\n", sli[:])
	arr := [1]int{7}
	fmt.Printf("%T\n", arr)

	sli2 := make([]int, 3, 7)
	fmt.Println("sli2 = ", sli2)
	fmt.Println("len(sli2) = ", len(sli2))
	fmt.Println("cap(sli2) = ", cap(sli2))
	fmt.Printf("%T\n", sli2)

	sli3 := make([]int, 3)
	fmt.Println("sli3 = ", sli3)
	fmt.Println("len(sli3) = ", len(sli3))
	fmt.Println("cap(sli3) = ", cap(sli3))
}

func testCase02() {
	arr := [7]int{1, 2, 3, 4, 5, 6, 7}
	sli := arr[2:5:5] // {3,4,5}
	fmt.Println(sli)
	fmt.Println("len(sli) = ", len(sli)) // 5 - 2
	fmt.Println("cap(sli) = ", cap(sli)) // 5 - 2
}

func testCase01() {
	arr := [7]int{1, 2, 3, 4, 5, 6, 7}

	sli := arr[1:3]
	fmt.Println(sli)
	fmt.Println("len(sli) = ", len(sli)) // 3 - 1
	fmt.Println("cap(sli) = ", cap(sli)) // 6 - 0

	sli1 := arr[1:4:5]
	fmt.Println(sli1)
	fmt.Println("len(sli1) = ", len(sli1)) // 3 - 1
	fmt.Println("cap(sli1) = ", cap(sli1)) // 5 - 1

	sli2 := arr[0:7]
	fmt.Println(sli2)
	fmt.Println("len(sli2) = ", len(sli2)) // 7 - 0
	fmt.Println("cap(sli2) = ", cap(sli2)) // 7 - 0

	sli3 := arr[:]
	fmt.Println(sli3)
	fmt.Println("len(sli3) = ", len(sli3)) // 7 - 0
	fmt.Println("cap(sli3) = ", cap(sli3)) // 7 - 0
}
