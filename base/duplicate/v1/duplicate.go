/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-09-29
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
	"bufio"
	"fmt"
	"os"
)

/*
dup输出标准输入流中的出现多次的行，在行内容前是出现次数的计数。这个程序将引入if 表达式，map内置数据结构和bufio的package。
*/
func main() {
	// 创建一个字典用于统计行和该行出现的次数
	counts := make(map[string]int)
	// 创建一个bufio读取系统的标准输入
	scanner := bufio.NewScanner(os.Stdin)

	// 开始循环输入数据
	for scanner.Scan() {
		content := scanner.Text()
		if content == "" {
			break
		}
		counts[content]++
	}

	// 打印输入的内容
	for line, count := range counts {
		if count > 1 {
			fmt.Printf("%d\t%s\n", count, line)
		}
	}
}

/*
注意:
对map进行range循环时，其迭代顺序是不确定的，从实践来看， 很可能每次运行都会有不一样的结果
(译注:这是Go语言的设计者有意为之的，因为其底层实现不保证 插入顺序和遍历顺序一致，也希望程序员不要依赖遍历时的顺序，所以干脆直接在遍历的时候做了随机 化处理，醉了。
补充:好像说随机序可以防止某种类型的攻击)，来避免程序员在业务中依赖遍历时的顺序。
*/
