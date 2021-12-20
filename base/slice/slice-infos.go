package main

import (
	"fmt"
	"unsafe"
)

func main() {
	arrLength()
	defaultArrValue()
	displayDeclaration()
	setValueIndex()
	errDefineArr()
	multidimensionalArray()
	declaredSlice()
	autoScala()
	appendTrap()
}

/*
append 操作的这种自动扩容行为，有些时候会给我们开发者带来一些困惑，
比如基于一个已有数组建立的切片，一旦追加的数据操作触碰到切片的容量上限（实质上也是数组容量的上界)，
切片就会和原数组解除“绑定”，后续对切片的任何修改都不会反映到原数组中了。

在 append 25 之后，切片的元素已经触碰到了底层数组 u 的边界了。
然后我们再 append 26 之后，append 发现底层数组已经无法满足 append 的要求，
于是新创建了一个底层数组（数组长度为 cap(s) 的 2 倍，即 8），并将 slice 的元素拷贝到新数组中了。
在这之后，我们即便再修改切片的第一个元素值，原数组 u 的元素也不会发生改变了，
因为这个时候切片 s 与数组 u 已经解除了“绑定关系”，s 已经不再是数组 u 的“描述符”了。
这种因切片的自动扩容而导致的“绑定关系”解除，有时候会成为一个小陷阱，要注意这一点。
*/
func appendTrap() {
	u := [...]int{11, 12, 13, 14, 15}
	fmt.Println("array:", u) // [11, 12, 13, 14, 15]
	s := u[1:3]
	fmt.Printf("slice(len=%d, cap=%d): %v\n", len(s), cap(s), s) // [12, 13]
	s = append(s, 24)
	fmt.Println("after append 24, array:", u)
	fmt.Printf("after append 24, slice(len=%d, cap=%d): %v\n", len(s), cap(s), s)
	s = append(s, 25)
	fmt.Println("after append 25, array:", u)
	fmt.Printf("after append 25, slice(len=%d, cap=%d): %v\n", len(s), cap(s), s)
	s = append(s, 26)
	fmt.Println("after append 26, array:", u)
	fmt.Printf("after append 26, slice(len=%d, cap=%d): %v\n", len(s), cap(s), s)
	s[0] = 22
	fmt.Println("after reassign 1st elem of slice, array:", u)
	fmt.Printf("after reassign 1st elem of slice, slice(len=%d, cap=%d): %v\n", len(s), cap(s), s)
}

/*
“动态扩容”指的就是，当通过 append 操作向切片追加数据的时候，如果这时切片的 len 值和 cap 值是相等的，也就是说切片底层数组已经没有空闲空间再来存储追加的值了，Go 运行时就会对这个切片做扩容操作，来保证切片始终能存储下追加的新值。

切片变量 nums 之所以可以存储下新追加的值，就是因为 Go 对其进行了动态扩容，也就是重新分配了其底层数组，从一个长度为 6 的数组变成了一个长为 12 的数组。

*/
func autoScala() {
	var s []int
	s = append(s, 11)
	fmt.Println(len(s), cap(s)) //1 1
	s = append(s, 12)
	fmt.Println(len(s), cap(s)) //2 2
	s = append(s, 13)
	fmt.Println(len(s), cap(s)) //3 4
	s = append(s, 14)
	fmt.Println(len(s), cap(s)) //4 4
	s = append(s, 15)
	fmt.Println(len(s), cap(s)) //5 8
}

/*
type slice struct {
    array unsafe.Pointer
    len   int
    cap   int
}

*/
func declaredSlice() {
	var nums = []int{1, 2, 3, 4, 5, 6}
	fmt.Println(nums)
	fmt.Println(len(nums))
	nums = append(nums, 7)
	fmt.Println(len(nums))
	//方法一：通过 make 函数来创建切片，并指定底层数组的长度。我们直接看下面这行代码：
	//sl := make([]byte, 6, 10) // 其中10为cap值，即底层数组长度，6为切片的初始长度

	//如果没有在 make 中指定 cap 参数，那么底层数组长度 cap 就等于 len，比如：
	//sl := make([]byte, 6) // cap = len = 6

	// 方法二：采用 array[low : high : max]语法基于一个已存在的数组创建切片。这种方式被称为数组的切片化，比如下面代码：
	arr := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sl := arr[3:7:9]
	sl[0] += 10
	// 我们看到，基于数组创建的切片，它的起始元素从 low 所标识的下标值开始，切片的长度（len）是 high - low，它的容量是 max - low。
	// 而且，由于切片 sl 的底层数组就是数组 arr，对切片 sl 中元素的修改将直接影响数组 arr 变量。
	// 比如，如果我们将切片的第一个元素加 10，那么数组 arr 的第四个元素将变为 14：
	fmt.Println("arr[3] =", arr[3]) // 14
	// 切片好比打开了一个访问与修改数组的“窗口”，通过这个窗口，我们可以直接操作底层数组中的部分元素。这有些类似于我们操作文件之前打开的“文件描述符”（Windows 上称为句柄），通过文件描述符我们可以对底层的真实文件进行相关操作。可以说，切片之于数组就像是文件描述符之于文件。

	// 方法三：基于切片创建切片。
	//切片变量 nums 在进行一次 append 操作后切片容量变为 12 的问题。这里要清楚一个概念：切片与数组最大的不同，就在于其长度的不定长，这种不定长需要 Go 运行时提供支持，这种支持就是切片的“动态扩容”。
}

func multidimensionalArray() {
	var mArr [2][3][4]int
	fmt.Println(mArr)
}

func errDefineArr() {
	var arr = [6]int{11, 12, 13, 14, 15, 16}
	fmt.Println(arr[0], arr[5]) // 11 16
	//fmt.Println(arr[-1])        // 错误：下标值不能为负数
	//fmt.Println(arr[8])         // 错误：小标值超出了arr的长度范围
}

/*
但如果我们要对一个长度较大的稀疏数组进行显式初始化，这样逐一赋值就太麻烦了，还有什么更好的方法吗？
我们可以通过使用下标赋值的方式对它进行初始化，比如下面代码。
*/
func setValueIndex() {
	var arr = [...]int{
		0:  0,
		99: 39, // 将第100个元素(下标值为99)的值赋值为39，其余元素值均为0
	}
	fmt.Printf("%T\n", arr) // [100]int
}

/*
如果要显式地对数组初始化，我们需要在右值中显式放置数组类型，并通过大括号的方式给各个元素赋值（如下面代码中的 arr）。
当然，我们也可以忽略掉右值初始化表达式中数组类型的长度，用“…”替代，Go 编译器会根据数组元素的个数，自动计算出数组长度（如下面代码中的 arr3）：
*/
func displayDeclaration() {
	var arr = [6]int{
		11, 12, 13, 14, 15, 16,
	} // [11 12 13 14 15 16]
	fmt.Println(arr)
	var arr2 = [...]int{
		21, 22, 23,
	} // [21 22 23]
	fmt.Printf("%T\n", arr2) // [3]int
}

/*
默认值
*/
func defaultArrValue() {
	var arr [6]int // [0 0 0 0 0 0]
	fmt.Println(arr)
}

/*
数组大小就是所有元素的大小之和，这里数组元素的类型为 int。
在 64 位平台上，int 类型的大小为 8，数组 arr 一共有 6 个元素，因此它的总大小为 6x8=48 个字节。
*/
func arrLength() {
	var arr = [6]int{1, 2, 3, 4, 5, 6}
	fmt.Println("数组长度：", len(arr))           // 6
	fmt.Println("数组大小：", unsafe.Sizeof(arr)) // 48
}
