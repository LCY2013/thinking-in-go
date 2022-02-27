package main

import (
	"errors"
	"fmt"
)

//nil error 值 != nil

type MyError struct {
	error
}

var ErrBad = MyError{
	error: errors.New("bad things happened"),
}

func bad() bool {
	return false
}

func returnsError() error {
	var p *MyError = nil
	if bad() {
		p = &ErrBad
	}
	return p
}

func returnsErrorV1() error {
	if bad() {
		return &ErrBad
	}
	return nil
}

func returnsErrorV2() error {
	var p error = nil
	if bad() {
		p = ErrBad
	}
	return p
}

/*
如果你是一个初学者，我猜你的的思路大概是这样的：p 为 nil，returnsError 返回 p，那么 main 函数中的 err 就等于 nil，于是程序输出 ok 后退出。
示例程序并未如我们前面预期的那样输出 ok。
程序显然是进入了错误处理分支，输出了 err 的值。
那这里就有一个问题了：明明 returnsError 函数返回的 p 值为 nil，为什么却满足了if err != nil的条件进入错误处理分支呢？

接口类型变量的内部表示
接口类型“动静兼备”的特性也决定了它的变量的内部表示绝不像一个静态类型变量（如 int、float64）那样简单，可以在$GOROOT/src/runtime/runtime2.go中找到接口类型变量在运行时的表示：
// $GOROOT/src/runtime/runtime2.go
type iface struct {
    tab  *itab
    data unsafe.Pointer
}
type eface struct {
    _type *_type
    data  unsafe.Pointer
}
在运行时层面，接口类型变量有两种内部表示：iface和eface，这两种表示分别用于不同的接口类型变量：
1、eface 用于表示没有方法的空接口（empty interface）类型变量，也就是 interface{}类型的变量；
2、iface 用于表示其余拥有方法的接口 interface 类型变量。

这两个结构的共同点是它们都有两个指针字段，并且第二个指针字段的功能相同，都是指向当前赋值给该接口类型变量的动态类型变量的值。

那它们的不同点在哪呢？
就在于 eface 表示的空接口类型并没有方法列表，因此它的第一个指针字段指向一个_type类型结构，这个结构为该接口类型变量的动态类型的信息，它的定义是这样的：
// $GOROOT/src/runtime/type.go
type _type struct {
    size       uintptr
    ptrdata    uintptr // size of memory prefix holding all pointers
    hash       uint32
    tflag      tflag
    align      uint8
    fieldAlign uint8
    kind       uint8
    // function for comparing objects of this type
    // (ptr to object A, ptr to object B) -> ==?
    equal func(unsafe.Pointer, unsafe.Pointer) bool
    // gcdata stores the GC type data for the garbage collector.
    // If the KindGCProg bit is set in kind, gcdata is a GC program.
    // Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
    gcdata    *byte
    str       nameOff
    ptrToThis typeOff
}

iface 除了要存储动态类型信息之外，还要存储接口本身的信息（接口的类型信息、方法列表信息等）以及动态类型所实现的方法的信息，因此 iface 的第一个字段指向一个itab类型结构。itab 结构的定义如下：
// $GOROOT/src/runtime/runtime2.go
type itab struct {
    inter *interfacetype
    _type *_type
    hash  uint32 // copy of _type.hash. Used for type switches.
    _     [4]byte
    fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

可以看到，itab 结构中的第一个字段inter指向的 interfacetype 结构，存储着这个接口类型自身的信息。你看一下下面这段代码表示的 interfacetype 类型定义， 这个 interfacetype 结构由类型信息（typ）、包路径名（pkgpath）和接口方法集合切片（mhdr）组成。
// $GOROOT/src/runtime/type.go
type interfacetype struct {
    typ     _type
    pkgpath name
    mhdr    []imethod
}
itab 结构中的字段_type则存储着这个接口类型变量的动态类型的信息，字段fun则是动态类型已实现的接口方法的调用地址数组。

*/

func main() {
	//err := returnsError()
	err := returnsErrorV1()
	if err != nil {
		fmt.Printf("error occur: %+v\n", err)
		return
	}
	fmt.Println("ok")
}
