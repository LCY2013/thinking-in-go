package array

import "fmt"

// array 声明一个包含5个元素的数组
var array [5]int

func arrayFunc() {
	// 声明一个包含5个元素的整型数组，用具体值初始化每个元素
	_ = [5]int{10, 20, 30, 40, 50}

	// ============================

	// 声明一个包含2个元素的整型数组,由具体值初始化数组大小
	_ = [...]int{1, 2}

	// ============================

	// 声明一个包含5个元素的整型数组,由具体值初始化索引值为1 和 2 的元素
	_ = [5]int{1: 30, 2: 40}

	// ============================

	// 修改数组元素值
	array[1] = 10

	// ============================

	// 声明一个包含5个元素的指向整型的数组
	// 用整型指针初始化索引为0和1的数组元素
	intPoint := [5]*int{0: new(int), 1: new(int)}
	// 为索引 0 ， 1 负值
	*intPoint[0] = 10
	*intPoint[1] = 20

	// ============================

	// 声明一个包含5元素的字符串数组
	var arrayStr01 [5]string

	// 声明第二个包含5元素的字符串数组，用颜色初始化数组
	arrayStr02 := [5]string{"Red", "Blue", "Green", "Yellow", "Pink"}

	// 将 arrayStr02 负值给 arrayStr01
	arrayStr01 = arrayStr02

	fmt.Println(arrayStr01 == arrayStr02)

	// 错误示例
	// 声明一个包含4个元素的字符串
	var arrayStr03 [4]string

	// 将 arrayStr02 负值给 arrayStr03 , compile error
	// # command-line-arguments
	// ./array_info.go:44:13: cannot use arrayStr02 (type [5]string) as type [4]string in assignment
	// arrayStr03 = arrayStr02

	_ = arrayStr03

	// ===============================
	// 声明一个包含3个元素的指针字符串
	var arrayStr04 [3]*string

	// 声明第二个包含3个元素的指针字符串
	arrayStr05 := [3]*string{new(string), new(string), new(string)}

	// 使用颜色给每个指针字符串负值
	*arrayStr05[0] = "Red"
	*arrayStr05[1] = "Blue"
	*arrayStr05[2] = "Green"

	// 将 arrayStr05 负值给 arrayStr04
	arrayStr04 = arrayStr05

	fmt.Println(arrayStr04 == arrayStr05)
}

// mutArrayFunc 多维数组
func mutArrayFunc() {
	// 声明一个二维数组，两个纬度为 4 个元素和 两个元素
	var _ [4][2]int

	// 使用数组字面量来声明并初始化上面的二维数组
	_ = [4][2]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}}

	// 初始化外层数组中索引为1个和3的元素
	_ = [4][2]int{1: {3, 4}, 3: {7, 8}}

	// 声明并初始化内层和外层元素
	_ = [4][2]int{1: {1: 10}, 3: {0: 20}}
}

// acceptArrayFunc 访问多维数组
func acceptArrayFunc() {
	// 声明一个 2x2 的二位数组
	var array [2][2]int

	// 设置每个元素的整型值
	array[0][0] = 10
	array[0][1] = 20
	array[1][0] = 30
	array[1][1] = 40

	// 声明另一个2x2的二位数组
	var arrayOther [2][2]int
	arrayOther = array

	fmt.Println(array == arrayOther)

	// 将索引为1 的纬度信息负值给同类的新型数组
	var arrayInt [2]int = array[1]

	// 将外层索引为1，内层索引为0的整型值负值给新值
	newInt := array[1][1]

	fmt.Println(arrayInt, newInt)
}

// go 语言都是值拷贝，大数据最好使用指针类型，下面是错误示例
// funcToFuncArray 函数间传递数组
// 函数到函数之间传递大对象数组
func funcToFuncArray() {
	// 声明一个需要8 MB 的数组
	var array [1e6]int

	// 将数组传递给 foo函数
	foo(array)
}

// foo 函数接受一个 100 万个整数型的数组
func foo(array [1e6]int) {

}

// 大数据使用，正确示例如下
func funcToFuncArrayRight() {
	// 声明一个需要 8MB 的数组
	var array [1e6]int

	// 将数组传递给 bar函数
	bar(&array)
}

// bar 函数接受一个 100 万个整数型的指针数组
func bar(array *[1e6]int) {

}
