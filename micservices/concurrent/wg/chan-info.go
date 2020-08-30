/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2019 Ramostear.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-08-12
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
	"strings"
	"sync"
	"time"
)

func input(ch chan string) {
	defer wg.Done()
	defer close(ch)
	var input string
	fmt.Println("Enter 'EOF' to exit : ")
	for {
		_, err := fmt.Scanf("%s", &input)
		if err != nil {
			fmt.Println("read input err : ", err)
		}
		if input == "EOF" || strings.EqualFold(input, "EOF") {
			fmt.Println("bye")
			break
		}
		ch <- input
	}
}

func output(ch chan string) {
	defer wg.Done()
	for value := range ch {
		fmt.Println("Your input : ", value)
	}
}

// go switch 只会匹配其中一个，不需要break跳出匹配
func switchInfo() {
	nowTime := time.Now()
	switch nowTime.Weekday() {
	case time.Saturday:
		fmt.Println("take a rest")
	case time.Sunday:
		fmt.Println("take a rest")
	default:
		fmt.Println("you need to work")
	}
	switch {
	case nowTime.Weekday() >= time.Monday && nowTime.Weekday() <= time.Friday:
		fmt.Println("you need to work")
	default:
		fmt.Println("take a rest")
	}
}

var wg sync.WaitGroup

func main() {
	ch := make(chan string)
	// 等待组设置两个等待
	wg.Add(2)
	//wg.Add(1)
	go input(ch)
	go output(ch)
	wg.Wait()
	fmt.Println("main out")
	//runtime.GOMAXPROCS()
}
