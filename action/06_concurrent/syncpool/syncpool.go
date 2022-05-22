package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
Go 并发相关库 sync 里面有一个有趣的 package Pool，sync.Pool 是个有趣的库，用很少的代码实现了很巧的功能。
第一眼看到 Pool 这个名字，就让人想到池子，元素池化是常用的性能优化的手段（性能优化的几把斧头：并发，预处理，缓存）。
比如，创建一个 100 个元素的池，然后就可以在池子里面直接获取到元素，免去了申请和初始化的流程，大大提高了性能。
释放元素也是直接丢回池子而免去了真正释放元素带来的开销。

但是再仔细一看 sync.Pool 的实现，发现比我预期的还更有趣。
sync.Pool 除了最常见的池化提升性能的思路，最重要的是减少 GC 。
常用于一些对象实例创建昂贵的场景。
注意，Pool 是 Goroutine 并发安全的。

sync.Pool 只是本身的 Pool 数据结构是并发安全的，并不是说 Pool.New 函数一定是线程安全的。
Pool.New 函数可能会被并发调用 ，如果 New 函数里面的实现是非并发安全的，那就会有问题。
关于 createBuffer 函数的实现里，对于 numCalcCreated 的计数加是用原子操作的：atomic.AddInt32(&numCalcsCreated, 1) 。

为什么 sync.Pool 不适合用于像 socket 长连接或数据库连接池?

因为，不能对 sync.Pool 中保存的元素做任何假设，以下事情是都可以发生的：

Pool 池里的元素随时可能释放掉，释放策略完全由 runtime 内部管理；
Get 获取到的元素对象可能是刚创建的，也可能是之前创建好 cache 住的。使用者无法区分；
Pool 池里面的元素个数你无法知道；
所以，只有的你的场景满足以上的假定，才能正确的使用 Pool 。
sync.Pool 本质用途是增加临时对象的重用率，减少 GC 负担。
划重点：临时对象。所以说，像 socket 这种带状态的，长期有效的资源是不适合 Pool 的。

sync.Pool 本质用途是增加临时对象的重用率，减少 GC 负担；
不能对 Pool.Get 出来的对象做预判，有可能是新的（新分配的），有可能是旧的（之前人用过，然后 Put 进去的）；
不能对 Pool 池里的元素个数做假定，你不能够；
sync.Pool 本身的 Get, Put 调用是并发安全的，sync.New 指向的初始化函数会并发调用，里面安不安全只有自己知道；
当用完一个从 Pool 取出的实例时候，一定要记得调用 Put，否则 Pool 无法复用这个实例，通常这个用 defer 完成；
*/

// 统计真正的创建次数
var numCalcCreated int32

// createBuffer 第一个步骤就是创建一个 Pool 实例，关键一点是配置 New 方法，声明 Pool 元素创建的方法。
func createBuffer() interface{} {
	// 这里要注意下，非常重要的一点。这里必须使用原子加，不然有并发问题；
	atomic.AddInt32(&numCalcCreated, 1)
	buffer := make([]byte, 1024)
	return &buffer
}

func main() {
	// 创建实例
	bufferPool := sync.Pool{
		New: createBuffer,
	}

	// 多 goroutine 并发测试
	numWorkers := 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			// 申请一个 buffer 实例
			//13 buffer objects were created.
			buffer := bufferPool.Get()
			// 1048576 buffer objects were created.
			//buffer := createBuffer()
			_ = buffer.(*[]byte)
			// 释放回池
			bufferPool.Put(buffer)
		}()
	}

	wg.Wait()

	fmt.Printf("%d buffer objects were created.\n", numCalcCreated)
}
