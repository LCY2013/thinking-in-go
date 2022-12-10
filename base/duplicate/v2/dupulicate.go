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
	"strings"
)

/*
dup程序从 标准输入得到一些文件名，然后用os.Open函数来打开每一个文件获取内容。
*/
func main() {
	// 创建统计重复行的字典
	counts := make(map[string]int)
	countsFile := make(map[string]map[string]int)

	// 读取命令行参数信息
	files := os.Args[1:]

	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, filePath := range files {
			file, err := os.Open(filePath)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "v2 : %v\n", err)
				continue
			}
			countLines(file, counts)

			file.Close()

			file, err = os.Open(filePath)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "v2 : %v\n", err)
				continue
			}

			countFileLines(file, countsFile)

			file.Close()
		}
	}

	// 输出统计结果
	for content, num := range counts {
		fmt.Printf("%d\t%s\n", num, content)
	}

	fmt.Println("-----------------------")

	for content, fileMap := range countsFile {
		var fileResult string
		for fileName, num := range fileMap {
			fileResult += fmt.Sprintf("%s\t%d\t", fileName, num)
		}
		fmt.Printf("%s\n", content)
		fmt.Printf("\t%s\n", fileResult)
	}
}

/*
统计文件行里面的重复行
*/
func countLines(f *os.File, counts map[string]int) {
	// 将文件作为一个输入源
	scanner := bufio.NewScanner(f)

	// 统计行数
	for scanner.Scan() {
		content := scanner.Text()
		counts[content]++
	}
}

func countFileLines(f *os.File, counts map[string]map[string]int) {
	fileSplit := strings.Split(f.Name(), "/")
	fileName := fileSplit[len(fileSplit)-1]

	// 将文件作为一个输入源
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		content := scanner.Text()
		if counts[content] == nil {
			counts[content] = make(map[string]int)
		}
		counts[content][fileName]++
	}
}
