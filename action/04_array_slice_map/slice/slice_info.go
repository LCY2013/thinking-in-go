package slice

import "fmt"

func FuncSlice() {
	// 创建一个整数型切片
	// 其长度为3个元素，容量为5个元素
	_ = make([]int, 3, 5)

	// 创建一个整型切片，使其长度大于容量 - 错误示例
	// # command-line-arguments
	//./slice_info.go:9:10: len larger than cap in make([]int)
	//_ = make([]int, 5, 3)

	// 创建字符串切片，其长度和容量都是5个元素
	_ = []string{"Red", "Blue", "Green", "Yellow", "Pink"}

	// 创建一个整数型切片，其长度和容量都是3个元素
	_ = []int{1, 2, 3}

	// 创建字符串切片，第100个节点初始化一个空字符串
	_ = []string{99: ""}

	// 创建有三个元素的整型数组
	_ = [3]int{1, 2, 3}

	// 创建有三个元素的整型切片
	_ = []int{1, 2, 3}

	// 创建一个 nil 整型切片
	var _ []int
	// 使用make创建空的整型切片
	_ = make([]int, 0)
	// 使用字面量创建空的整型切片
	_ = []int{}

	// 创建一个整型切片，其容量和长度都是5
	slice := []int{1, 2, 3, 4, 5}
	// 改变 slice 切片的某个位置的值
	slice[2] = 30
	// 创建一个新的切片，该切片长度为2，从 slice 下标 1 开始
	newSlice := slice[1:3]
	_ = newSlice

	// 创建5个字符元素的切片，长度和容量都是5
	source := []string{"Apple", "Orange", "Plum", "Banana", "Grape"}

	// 将第三个元素切片，并限制容量，其长度为1个元素，容量为2个元素
	newSource := source[2:3:4]
	fmt.Println(newSource)

	// slice[i:j:k] 或者 [2:3:4]
	// 长度：j - i  3 - 2
	// 容量：k - i  4 - 2

	// 向 source 中新增元素
	source = append(source, "fufeng")
	fmt.Println(source)

	// 定义两个切片，并让切片相加
	s1 := []int{1, 3}
	s2 := []int{2, 4}
	s3 := append(s1, s2...)
	fmt.Println(s3)

	// 创建一个整型切片，其长度和容量都是4
	s4 := []int{100, 200, 300, 400}
	// 迭代每一个元素并显示其值
	for index, value := range s4 {
		fmt.Printf("Index: %d , Value: %d\n", index, value)
	}

	// range 创建的是切片的副本，如下显示每个元素以及其指针地址
	for index, value := range s4 {
		fmt.Printf("Value: %d, Value-Addr: %X, Element-Addr: %X\n",
			value, &value, &s4[index])
	}

	// 使用空白标识符来忽略索引值
	for _, value := range s4 {
		fmt.Printf("Value: %d\n", value)
	}

	// 传统 for 循环
	for index := 0; index < len(s4); index++ {
		fmt.Printf("Index: %d, Value: %d\n", index, s4[index])
	}

	// 多维切片
	sliceGroup := [][]int{{10}, {20, 30}}
	// 为第一个切片添加一个元素 20
	sliceGroup[0] = append(sliceGroup[0], 20)
	fmt.Println(sliceGroup)
}

// FuncToFuncSlice 函数间传递切片
func FuncToFuncSlice() {
	// 分配包含100万个整数型的切片
	slice := make([]int, 1e6)
	// 将slice传递给函数 foo
	foo(slice)
}

// foo 函数接受一个 整数型的slice，并且返回该slice
func foo(slice []int) []int {
	return slice
}
