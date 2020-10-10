/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-10-10
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

type Celsius float64

type Fahrenheit float64

const (
	AbsoluteZeroC Celsius = -273.15 // 绝对零度
	FreezingC     Celsius = 0       // 结冰点温度
	BoilingC      Celsius = 100     // 水沸点
)

/*
每一个类型T，都有一个对应的类型转换操作T(x)，用于将x转为T类型(译注:如果T是指针类型，可能会需要用小括弧包装T，比如(*int)(0))。
*/
func CToF(c Celsius) Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

func FToC(f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

// Celsius 类型含有一个方法叫String，每个类型只要实现了String方法，在打印的时候%v，%s都会走这里
func (c Celsius) String() string {
	return fmt.Sprintf("%gºC", c)
}

func main() {
	fmt.Printf("%g\n", BoilingC-FreezingC)
	boilingC := CToF(BoilingC)
	fmt.Printf("%g\n", boilingC-CToF(FreezingC))
	// fmt.Printf("%g\n",boilingC - FreezingC) compile error
	fmt.Printf("%b\n", 100)

	fmt.Println()
	var c Celsius
	var f Fahrenheit
	fmt.Println(c == 0)
	fmt.Println(f >= 0)
	// fmt.Println(c == f) type mismatch
	fmt.Println(c == FToC(f))
	fmt.Println(c == Celsius(f))

	fmt.Println()
	fmt.Println(c.String())

	fmt.Println()
	toC := FToC(212)
	fmt.Printf("%v \n", toC)
	fmt.Printf("%s \n", toC)
	fmt.Println(c)
	fmt.Printf("%g \n", toC)
	fmt.Println(float64(toC))
}
