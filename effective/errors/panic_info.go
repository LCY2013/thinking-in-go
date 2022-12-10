/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-06-17
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package errors

import (
	"fmt"
	"os"
)

/*
Panic

向调用者报告错误的一般方式就是将 error 作为额外的值返回。
标准的 Read 方法就是个众所周知的实例，它返回一个字节计数和一个 error。
但如果错误时不可恢复的呢? 有时程序就是不能继续运行。

为此，我们提供了内建的 panic 函数，它会产生一个运行时错误并终止程序。
该函数接受一个任意类型的实参(一般为字符串)，并在程序终止时打印。
它还能 表明发生了意料之外的事情，比如从无限循环中退出了。
// 用牛顿法计算立方根的一个玩具实现。
func CubeRoot(x float64) float64 {
	z := x / 3 // 任意初始值
	for i := 0; i < 1e6; i++ {
		prevz := z
		z -= (z*z*z - x) / (3 * z * z)
		if veryClose(z, prevz) {
			return z
		}
	}
	// 一百万次迭代并未收敛，事情出错了。
	panic(fmt.Sprintf("CubeRoot(%g) did not converge", x))
}

这仅仅是个示例，实际的库函数应避免 panic。
若问题可以被屏蔽或解决，最好就是让程序 继续运行而不是终止整个程序。
一个可能的反例就是初始化: 若某个库真的不能让自己工作，且有足够理由产生 Panic，那就由它去吧。
var user = os.Getenv("USER")
func init() {
    if user == "" {
        panic("no value for $USER")
    }
}
*/

// CubeRoot 用牛顿法计算立方根的一个玩具实现。
func CubeRoot(x float64) float64 {
	z := x / 3 // 任意初始值
	for i := 0; i < 1e6; i++ {
		prevz := z
		z -= (z*z*z - x) / (3 * z * z)
		if veryClose(z, prevz) {
			return z
		}
	}
	// 一百万次迭代并未收敛，事情出错了。
	panic(fmt.Sprintf("CubeRoot(%g) did not converge", x))
}

func veryClose(z float64, prevz float64) bool {
	return false
}

var user = os.Getenv("USER")

func init() {
	if user == "" {
		panic("no value for $USER")
	}
}
