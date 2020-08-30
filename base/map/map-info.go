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
	/*	var m1 map[int]string			// 声明map ，没有空间，不能直接存储key -- value
		//m1[100] = "Green"
		if m1 == nil {
			fmt.Println("map is nil ")
		}

		m2 :=  map[int]string{}			//
		fmt.Println(len(m2))
		fmt.Println("m2 = ", m2)
		m2[4] = "red"
		fmt.Println("m2 = ", m2)

		m3 := make(map[int]string)
		fmt.Println(len(m3))
		fmt.Println("m3 = ", m3)
		m3[400] = "red"
		fmt.Println("m3 = ", m3)

		m4 := make(map[int]string, 5)		// len
		fmt.Println("len(m4) = ", len(m4))
		//fmt.Println("len(m4) = ", cap(m4))		// 不能在map中使用 cap（）
		fmt.Println("m4 = ", m4)*/

	/*	// 初始化map
		var m5 map[int]string = map[int]string{1:"Luffy", 130:"Sanji", 1301:"Zoro"}

		fmt.Println("m5 = ", m5)

		m6 := map[int]string{1:"Luffy", 130:"Sanji", 1303:"Zoro"}
		fmt.Println("m6 = ", m6)
	*/

	/*	m7 := make(map[int]string, 1)
		m7[100] = "Nami"
		m7[20] = "Hello"
		m7[3] = "world"
		fmt.Println("m7=", m7)

		m7[3] = "yellow"			// 成功！ 将原map中 key 值为 3 的map元素，替换。
		fmt.Println("m7=", m7)*/

	// 遍历map
	/*	var m8 map[int]string = map[int]string{1:"Luffy", 130:"Sanji", 1301:"Zoro"}
		for k, v := range m8 {
			fmt.Printf("key:%d --- value:%q\n", k, v)
		}

		// range返回的key/ value 。 省略value打印。
		for _, K := range m8 {
			fmt.Printf("key:%s\n", K)
		}*/

	// 判断 map 中的key 是否存在
	/*var m9 map[int]string = map[int]string{1:"Luffy", 130:"Sanji", 1301:"Zoro"}

	if v, has := m9[12]; has {	// m9[下标] 返回两个值，第一个是value，第二个是bool 代表key是否存在。
		fmt.Println("value=", v, "has=", has)
	} else {
		fmt.Println("false value=", v, "has=", has)
	}*/

	/*m10 := map[int]string{1:"Luffy", 130:"Sanji", 1301:"Zoro"}
	fmt.Println("before delete m :", m10)

	mapDelete(m10, 130)

	fmt.Println("after delete m :", m10)*/

	str := "I love my work and I I I I love love love my family too"
	mRet := wordCountFunc(str)

	//mRet := wordCountFunc2(str)

	// 遍历map ，展示每个word 出现的次数：
	for k, v := range mRet {
		fmt.Printf("%q:%d\n", k, v)
	}
}

func wordCountFunc(str string) map[string]int {
	s := strings.Fields(str)  // 将字符串，拆分成 字符串切片s
	m := make(map[string]int) // 创建一个用于存储 word 出现次数的 map

	// 遍历拆分后的字符串切片
	for i := 0; i < len(s); i++ {
		if _, ok := m[s[i]]; ok { // ok == ture 说明 s[i] 这个key存在
			m[s[i]] = m[s[i]] + 1 // m[s[i]]++
		} else { // 说明 s[i] 这个key不存在， 第一次出现。添加到map中
			m[s[i]] = 1
		}
	}
	return m
}

func wordCountFunc2(str string) (m map[string]int) {
	m = make(map[string]int)
	arr := strings.Fields(str)
	for _, v := range arr {
		m[v]++
	}
	return
}

// map做函数参数、返回值，传引用
func mapDelete(m map[int]string, key int) {
	delete(m, key) // 删除 m 中 键值为 key的 map 元素
}
