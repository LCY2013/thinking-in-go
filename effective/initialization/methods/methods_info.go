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

import "fmt"

// ByteSlice
/*
Pointers vs. Values

正如 constants_info.ByteSize 那样，可以为任何已命名的类型(除了指针或接口)定义方法; 接收者可 不必为结构体。

在之前讨论切片时，编写了一个 Append 函数。 也可将其定义为切片的方法。为此，首先要声明一个已命名的类型来绑定该方法，然后使该方法的接收者成为该类型的值。
*/
type ByteSlice []byte

// AppendValue
/*
下面这样仍然需要该方法返回更新后的切片。
*/
func (slice ByteSlice) AppendValue(data []byte) []byte {
	return append(slice, data...)
}

//Append
/*
为了消除这种不便，可通过重新定义该方法， 将一个指向 ByteSlice 的指针作为该方法的接收者， 这样该方法就能重写调用者提供的切片了。
*/
func (p *ByteSlice) Append(data []byte) {
	slice := *p
	// 这里就可以没有return了
	slice = append(slice, data...)
	*p = slice
}

//Waite
/*
其实可以做得更好。若将函数修改为与标准 Write 类似的方法，就像下面这样。

那么类型 *ByteSlice 就满足了标准的 io.Writer 接口，这将非常实用。 例如，我们可以通过 打印将内容写入。
*/
func (p *ByteSlice) Write(data []byte) (n int, err error) {
	slice := *p
	slice = append(slice, data...)
	*p = slice
	return len(*p), nil
}

/*
下面将 ByteSlice 的地址传入，因为只有 *ByteSlice 才满足 io.Writer。
以指针或值为接收者的区别在于: 值方法可通过指针和值调用， 而指针方法只能通过指针来调用。

之所以会有这条规则是因为指针方法可以修改接收者;
通过值调用它们会导致方法接收到该值的副本，因此任何修改都将被丢弃，因此该语言不允许这种错误。
不过有个方便的例外: 若该值是可寻址的，那么该语言就会自动插入取址操作符来对付一般的通过值调用的指针方法。
在下面的例子中，变量 b 是可寻址的，因此只需通过 b.Write 来调用它的 Write 方法，编译器会将它重写为 (&b).Write。

顺便一提，在字节切片上使用 Write 的想法已被 bytes.Buffer 所实现。
*/
func main() {
	var b ByteSlice
	length, err := fmt.Fprintf(&b, "This hour has %d days\n", 7)
	if err != nil {
		return
	}
	fmt.Println(length)

	c := &ByteSlice{}
	*c = c.AppendValue([]byte("hello"))
	fmt.Printf("%T\t%#v\n", c, c)

	d := ByteSlice{}
	d = d.AppendValue([]byte("world"))
	fmt.Printf("%T\t%+v\n", d, d)
}
