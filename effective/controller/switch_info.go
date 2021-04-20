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

import "fmt"

/*
switch 也可用于判断接口变量的动态类型。如 类型选择 通过圆括号中的关键字 type 使用类 型断言语法。
若 switch 在表达式中声明了一个变量，那么该变量的每个子句中都将有该变量 对应的类型。在这些 case 中重用一个名字也是符合语义的，
实际上是在每个 case 里声明了 一个不同类型但同名的新变量。
*/
func main() {
	switchInfo()
}

func switchInfo() {
	var t interface{}
	// t = functionOfSomeType()
	switch ttype := t.(type) {
	default:
		fmt.Printf("unexpected type %T", ttype)
	case bool:
		fmt.Printf("boolean %t\n", ttype)
	case int:
		fmt.Printf("integer %d\n", ttype)
	case *bool:
		fmt.Printf("pointer to boolean %t\n", *ttype)
	case *int:
		fmt.Printf("pointer to int %d\n", *ttype)
	}
}
