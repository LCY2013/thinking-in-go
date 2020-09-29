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
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

/*
ReadFile函数需要一 个文件名参数
*/
func main() {
	// 创建统计文本的字典
	counts := make(map[string]int)

	// 遍历命令行参数信息
	for _, filePath := range os.Args[1:] {
		// 通过ReadFile 读取文件字节内容
		fileByte, err := ioutil.ReadFile(filePath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "v3 : %v\n", err)
			continue
		}
		for _, line := range strings.Split(string(fileByte), "\n") {
			counts[line]++
		}
	}

	// 输出命令行参数信息
	for content, num := range counts {
		fmt.Printf("%d\t%s\n", num, content)
	}
}
