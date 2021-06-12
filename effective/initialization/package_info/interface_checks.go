/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-06-11
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package package_info

import (
	"encoding/json"
	"fmt"
)

/*
Interface checks
接口检查

就像我们在前面 接口 中讨论的那样， 一个类型无需显式地声明它实现了某个接口。
取而代之，该类型只要实现了某个接口的方法，其实就实现了该接口。
在实践中，大部分接口转换都是静态的，因此会在编译时检测。
例如，将一个 *os.File 传入一个预期的 io.Reader 函数将不会被编译， 除非 *os.File 实现了 io.Reader 接口。

尽管有些接口检查会在运行时进行。encoding/json 包中就有个实例它定义了一个 Marshaler 接口。
当 JSON 编码器接收到一个实现了该接口的值，那么该编码器就会调用该值的编组方法，将其转换为JSON，而非进行标准的类型转换。
编码器在运行时通过类型断言 检查其属性，就像这样:
	m,ok := val.(json.Marshaler)
*/

func MarshalerFunc(param interface{}) {
	m, ok := param.(json.Marshaler)
	_ = m
	_ = ok
}

/*
若只需要判断某个类型是否是实现了某个接口，而不需要实际使用接口本身 (可能是错误检 查部分)，
就使用空白标识符来忽略类型断言的值:
*/

func MarshalerInterfaceFunc(param interface{}) {
	if _, ok := param.(json.Marshaler); ok {
		fmt.Printf("value %v of type %T implements json.Marshaler\n", param, param)
	}
}

/*
当需要确保某个包中实现的类型一定满足该接口时，就会遇到这种情况。
若某个类型(例如 json.RawMessage)需要一种定制的 JSON 表现时，它应当实现 json.Marshaler，不过现在没有静态转换可以让编译器去自动验证它。
若该类型通过忽略转换失败来满足该接口，那么 JSON 编码器仍可工作，但它却不会使用定制的实现。
为确保其实现正确，可在该包中用空白标识符声明一个全局变量:
var _ json.Marshaler = (*RawMessage)(nil)

在此声明中调用了一个 *RawMessage 转换并将其赋予了 Marshaler，
以此来要求 *RawMessage 实现 Marshaler，这时其属性就会在编译时被检测。
若 json.Marshaler 接口被更改，此包将无法通过编译，而我们则会注意到它需要更新。

在这种结构中出现空白标识符，即表示该声明的存在只是为了类型检查。不过请不要为满足接口就将它用于任何类型。
作为约定，仅当代码中不存在静态类型转换时才能这种声明，毕竟这是种罕见的情况。
*/

//var _ json.Marshaler = (*RawMessage)(nil)

func MarshalerJson() {

}
