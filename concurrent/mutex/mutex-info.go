/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-11-09
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
	"sync"
)

/*
go1.15.2/src/sync/mutex.go:31

// A Locker represents an object that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}

Mutex、RWMutex 都实现了该接口

*/
func main() {
	// data race
	// case01()

	// mutex
	// case02()

	// mutex struct
	// case03()

	case04()
}

/*
  多个goroutine同时修改一个变量，看是否会造成结果不一致

  结论:
	造成结果(<100000)和预期(100000)不一致

  原因分析:
	count++ 不是一个原子操作，至少包含几个步骤，比如读取变量count 的当前值，对这个值加 1，把结果再保存到 count 中，因为不是原子操作，就可能有并发的问题。

	汇编层如下：
		// count++操作的汇编代码
		MOVQ    "".count(SB),
		AX    LEAQ    1(AX),
		CX    MOVQ    CX,
		"".count(SB)

   工具集:
	  Go 提供了一个检测并发访问共享资源是否有问题的工具： race detector，可以帮助我们自动发现程序有没有 data race 的问题。

	  Go race detector 是基于 Google 的 C/C++ sanitizers 技术实现的，编译器通过探测所有的内存访问，加入代码能监视对这些内存地址的访问（读还是写）。
	  在代码运行的时候，race detector 就能监控到对共享变量的非同步访问，出现 race 的时候，就会打印出警告信息。

	  使用流程：
		在编译（compile）、测试（test）或者运行（run）Go 代码的时候，加上 race 参数，就有可能发现并发问题。
	    比如在上面的例子中，我们可以加上 race 参数运行，检测一下是不是有并发问题。
		如果你 go run -race mutex-info.go，就会输出警告信息。

	开启了 race 的程序部署在线上，还是比较影响性能的

	运行 go tool compile -race -S mutex-info.go，可以查看计数器例子的代码，重点关注一下 count++ 前后的编译后的代码。


*/
func case01() {
	// 定义一个初始化的统计变量
	var count = 0
	// 定义一个等待组
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(num int) {
			// 等待组释放
			defer wg.Done()
			//fmt.Println("wg ",num)
			// 这里对count进行递增操作
			for k := 0; k < 10000; k++ {
				count++
			}
		}(i)
	}

	// 休眠等待goroutine执行完成
	//time.Sleep(time.Second)

	// 等待goroutine执行完成
	wg.Wait()

	fmt.Println(count)
}

//===============================
/*
	利用mutex解决数据竞态问题
*/
func case02() {
	// 定义一个初始化的统计变量
	var count = 0
	// 定义一个等待组
	var wg sync.WaitGroup
	wg.Add(10)

	// 定义互斥锁
	//mutex := sync.Mutex{}
	// Mutex 的零值是还没有 goroutine 等待的未加锁的状态，所以你不需要额外的初始化，直接声明变量即可
	var mutex sync.Mutex

	for i := 0; i < 10; i++ {
		go func(num int) {
			// 等待组释放
			defer wg.Done()
			//fmt.Println("wg ",num)
			// 这里对count进行递增操作
			for k := 0; k < 10000; k++ {
				mutex.Lock()
				count++
				mutex.Unlock()
			}
		}(i)
	}

	// 休眠等待goroutine执行完成
	//time.Sleep(time.Second)

	// 等待goroutine执行完成
	wg.Wait()

	fmt.Println(count)
}

//===============================
func case03() {
	var counter Counter03

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for k := 0; k < 10000; k++ {
				counter.mutex.Lock()
				counter.count++
				counter.mutex.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Println(counter.count)
}

// 通过定义一个结构体绑定互斥锁
type Counter03 struct {
	mutex sync.Mutex
	count int
}

//===============================
func case04() {
	var counter Counter

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for k := 0; k < 10000; k++ {
				counter.incr()
			}
		}()
	}

	wg.Wait()
	fmt.Println(counter.counts())
}

type Counter struct {
	CounterType int
	Name        string

	mutex sync.Mutex
	count uint64
}

func (counter *Counter) incr() {
	defer counter.mutex.Unlock()
	counter.mutex.Lock()
	counter.count++
}

func (counter *Counter) counts() uint64 {
	defer counter.mutex.Unlock()
	counter.mutex.Lock()
	return counter.count
}

// 归并排序
func mergeSort(r []int) []int {
	length := len(r)
	if length <= 1 {
		return r
	}
	num := length / 2
	left := mergeSort(r[:num])
	right := mergeSort(r[num:])
	return merge(left, right)
}

func merge(left, right []int) (result []int) {
	l, r := 0, 0
	for l < len(left) && r < len(right) {
		if left[l] < right[r] {
			result = append(result, left[l])
			l++
		} else {
			result = append(result, right[r])
			r++
		}
	}
	result = append(result, left[l:]...)
	result = append(result, right[r:]...)
	return
}
