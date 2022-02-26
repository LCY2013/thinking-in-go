package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

// 增加 Goroutine 标识
// v02 版本的 Trace 函数在面对只有一个 Goroutine 的时候，还是可以支撑的，
// 但当程序中并发运行多个 Goroutine 的时候，多个函数调用链的出入口信息输出就会混杂在一起，无法分辨。
// 继续对 Trace 函数进行改造，让它支持多 Goroutine 函数调用链的跟踪。
// 方案就是在输出的函数出入口信息时，带上一个在程序每次执行时能唯一区分 Goroutine 的 Goroutine ID。
// Goroutine 也没有 ID 信息啊！的确如此，Go 核心团队为了避免Goroutine ID 的滥用，故意没有将 Goroutine ID 暴露给开发者。
// 但在 Go 标准库的 h2_bundle.go 中，却发现了一个获取 Goroutine ID 的标准方法，看下面代码：
/*
// $GOROOT/src/net/http/h2_bundle.go
var http2goroutineSpace = []byte("goroutine ")
func http2curGoroutineID() uint64 {
    bp := http2littleBuf.Get().(*[]byte)
    defer http2littleBuf.Put(bp)
    b := *bp
    b = b[:runtime.Stack(b, false)]
    // Parse the 4707 out of "goroutine 4707 ["
    b = bytes.TrimPrefix(b, http2goroutineSpace)
    i := bytes.IndexByte(b, ' ')
    if i < 0 {
        panic(fmt.Sprintf("No space found in %q", b))
    }
    b = b[:i]
    n, err := http2parseUintBytes(b, 10, 64)
    if err != nil {
        panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
    }
    return n
}
*/

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
	fmt.Printf("g[%05d]: enter: %s\n", gid, name)
	return func() {
		fmt.Printf("g[%05d]: exit: %s\n", gid, name)
	}
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
