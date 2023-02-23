//go:build amd64 || arm64

package tools

import (
	"reflect"
	"sort"
	"unsafe"
)

/*
切片类型强制转换

为了安全，当两个切片类型 []T 和 []Y 的底层原始切片类型不同时，Go 语言是无法直接转换类型的。
不过安全都是有一定代价的，有时候这种转换是有它的价值的——可以简化编码或者是提升代码的性能。
比如在 64 位系统上，需要对一个 []float64 切片进行高速排序，我们可以将它强制转为 []int 整数切片，
然后以整数的方式进行排序（因为 float64 遵循 IEEE754 浮点数标准特性，当浮点数有序时对应的整数也必然是有序的）。
*/

// 下面的代码通过两种方法将 []float64 类型的切片转换为 []int 类型的切片：

func SortFloat64FastV1(sortFloat []float64) {
	// 强制类型转换
	var sortInt []int = ((*[1 << 20]int)(unsafe.Pointer(&sortFloat[0])))[:len(sortFloat):cap(sortFloat)]

	// 以 int 方式给 float64 排序
	sort.Ints(sortInt)
}

func SortFloat64FastV2(sortFloat []float64) {
	// 通过 reflect.SliceHeader 更新切片头部信息实现转换
	var sortInt []int
	floatHdr := (*reflect.SliceHeader)(unsafe.Pointer(&sortFloat))
	intHdr := (*reflect.SliceHeader)(unsafe.Pointer(&sortInt))
	*intHdr = *floatHdr

	// 以 int 方式给 float64 排序
	sort.Ints(sortInt)
}

/*
第一种强制转换是先将切片数据的开始地址转换为一个较大的数组的指针，然后对数组指针对应的数组重新做切片操作。
中间需要 unsafe.Pointer 来连接两个不同类型的指针传递
需要注意的是，Go语言实现中非0大小数组的长度不得超过 2GB，
因此需要针对数组元素的类型大小计算数组的最大长度范围
（[]uint8 最大 2GB，[]uint16 最大 1GB，以此类推，但是 []struct{} 数组的长度可以超过 2GB）。

第二种转换操作是分别取到两个不同类型的切片头信息指针，任何类型的切片头部信息底层都是对应 reflect.SliceHeader 结构，
然后通过更新结构体方式来更新切片信息，从而实现 a 对应的 []float64 切片到 c 对应的 []int 类型切片的转换。

通过基准测试，我们可以发现用 sort.Ints 对转换后的 []int 排序的性能要比用 sort.Float64s 排序的性能好一点。
不过需要注意的是，这个方法可行的前提是要保证 []float64 中没有 NaN 和 Inf 等非规范的浮点数（
因为浮点数中 NaN 不可排序，正 0 和负 0 相等，但是整数中没有这类情形）。
*/
