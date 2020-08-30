/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2019 Ramostear.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo</a>
 * @date : 2020-08-08
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
	"os"
)

func main() {
	//openFilesOP()
	opFile()
}

/*
操作文件
rwx 124
*/
func opFile() {
	file, err := os.OpenFile("lcy.txt", os.O_RDWR, 6)
	if err != nil {
		fmt.Println("读取文件出错 ", err)
	}
	fmt.Println("读取文件成功")
	defer file.Close()

	writeString, err := file.WriteString("hi")
	if err != nil {
		fmt.Println("写文件出错 ", err)
	}
	fmt.Println("写入的字节数是 ", writeString)
}

/*打开文件操作*/
func openFilesOP() {
	create, err := os.Create("fufeng.txt")
	if err != nil {
		fmt.Printf("crate file fail %s\n", err)
	}
	defer create.Close()

	if open, err := os.Open("fufeng.txt"); err == nil {
		defer open.Close()
		_, err = open.WriteString("hello fufeng")
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("open file fail ", err)
	}

	_ = os.Rename("fufeng.txt", "lcy.txt")

	if openFile, err := os.OpenFile("lcy.txt", os.O_RDWR, 6); err == nil {
		defer openFile.Close()
		_, _ = openFile.WriteString("hello , fufeng ")
	} else {
		fmt.Println(err)
	}
}
