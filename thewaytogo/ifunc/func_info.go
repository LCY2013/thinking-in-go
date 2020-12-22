/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-12-07
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package ifunc

import "errors"

type Stack []int

// 出栈操作
// 没有参数的函数通常被称为 niladic 函数（niladic function），就像 main.main()。
func (st *Stack) Pop() int {
	v := 0
	for ix := len(*st) - 1; ix >= 0; ix-- {
		if v = (*st)[ix]; v != 0 {
			(*st)[ix] = 0
			return v
		}
	}
	return v
}

// 获取并且移除
func (st *Stack) Remove() (value int, err error) {
	if *st != nil {
		value = st.Pop()
		*st = (*st)[1:]
		err = nil
	} else {
		err = errors.New("fail")
	}
	return
}

// 入栈操作
func (st *Stack) Push(value int) (success bool, err error) {
	*st = append(*st, value)
	success = true
	err = nil
	return
}
