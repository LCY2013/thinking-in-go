package main

import "fmt"

/*
defer 函数支持

第一点：明确哪些函数可以作为 deferred 函数

Go 语言中除了自定义函数 / 方法，还有 Go 语言内置的 / 预定义的函数，Go 语言内置函数的完全列表：
Functions:
  append cap close complex copy delete imag len
  make new panic print println real recover

第二点：注意 defer 关键字后面表达式的求值时机
一定要牢记一点：defer 关键字后面的表达式，是在将 deferred 函数注册到 deferred 函数栈的时候进行求值的。

第三点：知晓 defer 带来的性能损耗

*/
func bar() (int, int) {
	return 1, 2
}

// foo 这里展示了不同的支持和不支持defer 的内置函数
// append、cap、len、make、new、imag 等内置函数都是不能直接作为 deferred 函数的，而 close、copy、delete、print、recover 等内置函数则可以直接被 defer 设置为 deferred 函数。
// 那些不能直接作为 deferred 函数的内置函数，我们可以使用一个包裹它的匿名函数来间接满足要求，以 append 为例是这样的
// defer func() {
//   _ = append(sl, 11)
// }()
func foo() {
	var c chan int
	var sl []int
	var m = make(map[string]int, 10)
	m["item1"] = 1
	m["item2"] = 2
	//var a = complex(1.0, -1.4)

	var sl1 []int
	defer bar()
	//defer append(sl, 11)
	//defer cap(sl)
	defer close(c)
	//defer complex(2, -2)
	defer copy(sl1, sl)
	defer delete(m, "item2")
	//defer imag(a)
	//defer len(sl)
	//defer make([]int, 10)
	//defer new(*int)
	defer panic(1)
	defer print("hello, defer\n")
	defer println("hello, defer")
	//defer real(a)
	defer recover()
}

// foo1 defer 关键字后面表达式的求值时机
func foo1() {
	for i := 0; i <= 3; i++ {
		defer fmt.Println(i)
	}
}
func foo2() {
	for i := 0; i <= 3; i++ {
		defer func(n int) {
			fmt.Println(n)
		}(i)
	}
}
func foo3() {
	for i := 0; i <= 3; i++ {
		defer func() {
			fmt.Println(i)
		}()
	}
}

func main() {
	//foo()
	fmt.Println("foo1 result:")
	foo1()
	fmt.Println("\nfoo2 result:")
	foo2()
	fmt.Println("\nfoo3 result:")
	foo3()
}
