package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

// 利用代码生成自动注入 Trace 函数
// 要实现向目标代码中的函数 / 方法自动注入 Trace 函数，首先要做的就是将 Trace 函数相关的代码打包到一个 module 中以方便其他 module 导入。
// 下面就先来看看将 Trace 函数放入一个独立的 module 中的步骤。
// 将 Trace 函数放入一个独立的 module 中
// 创建一个名为 instrument_trace 的目录，进入这个目录后，通过 go mod init 命令创建一个名为 github.com/lcy2013/instrument_trace 的 module：
// $mkdir instrument_trace
// $cd instrument_trace
// $go mod init github.com/lcy2013/instrument_trace
// 将最新版的 trace.go 放入到该目录下，将包名改为 trace，并仅保留 Trace 函数、Trace 使用的函数以及包级变量，其他函数一律删除掉。

// 不过，由于 http2curGoroutineID 不是一个导出函数，我们无法直接使用。可以把它复制出来改造一下：
var goroutineSpace = []byte("goroutine ")

// curGoroutineID 改造了两个地方。
// 一个地方是通过直接创建一个 byte 切片赋值给 b，替代原 http2curGoroutineID 函数中从一个 pool 池获取 byte 切片的方式，
// 另外一个是使用 strconv.ParseUint 替代了原先的 http2parseUintBytes。
// 改造后，就可以直接使用 curGoroutineID 函数来获取 Goroutine 的 ID 信息了。
func curGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	// Parse the 4707 out of "goroutine 4707 ["
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	}
	return n
}

var mu sync.Mutex
var m = make(map[uint64]int)

// Trace 使用了一个 map 类型变量 m 来保存每个 Goroutine 当前的缩进信息：m 的 key 为 Goroutine 的 ID，值为缩进的层次。
// 然后，考虑到 Trace 函数可能在并发环境中运行，根据 map 不支持并发写的注意事项，增加了一个 sync.Mutex 实例 mu 用于同步对 m 的写操作。
// 对于一个 Goroutine 来说，每次刚进入一个函数调用，就在输出入口跟踪信息之前，将缩进层次加一，并输出入口跟踪信息，加一后的缩进层次值也保存到 map 中。
// 然后，在函数退出前，取出当前缩进层次值并输出出口跟踪信息，之后再将缩进层次减一后保存到 map 中。
// 除了增加缩进层次信息外，在这一版的 Trace 函数实现中，也把输出出入口跟踪信息的操作提取到了一个独立的函数 printTrace 中，这个函数会根据传入的 Goroutine ID、函数名、箭头类型与缩进层次值，按预定的格式拼接跟踪信息并输出。
func Trace() func() {
	// 通过 runtime.Caller 函数获得当前 Goroutine 的函数调用栈上的信息
	// runtime.Caller 的参数标识的是要获取的是哪一个栈帧的信息。
	// 当参数为 0 时，返回的是 Caller 函数的调用者的函数信息，在这里就是 Trace 函数。
	// 但我们需要的是 Trace 函数的调用者的信息，于是我们传入 1。
	//
	// Caller 函数有四个返回值：
	// 第一个返回值代表的是程序计数（pc）；
	// 第二个和第三个参数代表对应函数所在的源文件名以及所在行数，这里我们暂时不需要；
	// 最后一个参数代表是否能成功获取这些信息，如果获取失败，抛出 panic。
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}
	// 通过 runtime.FuncForPC 函数和程序计数器（PC）得到被跟踪函数的函数名称。
	// runtime.FuncForPC 返回的名称中不仅仅包含函数名，还包含了被跟踪函数所在的包名。
	fn := runtime.FuncForPC(pc)
	// 获取上一个栈针的函数名称
	name := fn.Name()
	// 获取当前goroutine的信息
	gid := curGoroutineID()

	mu.Lock()
	indents := m[gid]
	m[gid] = indents + 1
	mu.Unlock()
	printTrace(gid, name, "->", indents+1)
	return func() {
		mu.Lock()
		indent := m[gid]
		m[gid] = indent - 1
		mu.Unlock()
		printTrace(gid, name, "<-", indent)
	}
}

// printTrace 格式化缩进输出
func printTrace(id uint64, name, arrow string, indent int) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents = fmt.Sprintf("%s%s", indents, "	")
	}
	fmt.Printf("g[%05d]: %s%s%s\n", id, indents, arrow, name)
}

// 将这个程序由单 Goroutine 改为多 Goroutine 并发的，这样才能验证支持多 Goroutine 的新版 Trace 函数是否好用：

func A1() {
	defer Trace()()
	B1()
}

func B1() {
	defer Trace()()
	C1()
}

func C1() {
	defer Trace()()
	D()
}

func D() {
	defer Trace()()
}

func A2() {
	defer Trace()()
	B2()
}

func B2() {
	defer Trace()()
	C2()
}

func C2() {
	defer Trace()()
	D()
}

/*
结果：
g[00001]: enter: main.A2
g[00001]: enter: main.B2
g[00001]: enter: main.C2
g[00001]: enter: main.D
g[00001]: exit: main.D
g[00001]: exit: main.C2
g[00001]: exit: main.B2
g[00001]: exit: main.A2
g[00006]: enter: main.A1
g[00006]: enter: main.B1
g[00006]: enter: main.C1
g[00006]: enter: main.D
g[00006]: exit: main.D
g[00006]: exit: main.C1
g[00006]: exit: main.B1
g[00006]: exit: main.A1

新示例程序输出了带有 Goroutine ID 的出入口跟踪信息，通过 Goroutine ID 可以快速确认某一行输出是属于哪个 Goroutine 的。

问题？
函数调用链跟踪已经支持了多 Goroutine，并且可以输出有层次感的跟踪信息了，
但对于 Trace 特性的使用者而言，他们依然需要手工在自己的函数中添加对 Trace 函数的调用。
那么是否可以将 Trace 特性自动注入特定项目下的各个源码文件中呢？
*/
func main() {
	//defer Trace()()
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		A1()
	}()
	A2()
	wg.Wait()
}
