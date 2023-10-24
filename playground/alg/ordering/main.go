package main

import (
	"fmt"
	"math/rand"
	"time"
)

func runtime(f func()) {
	now := time.Now().UnixMilli()
	f()
	fmt.Printf("speed time: %d\n", time.Now().UnixMilli()-now)
}

func main() {
	// 亿
	//length := 100000000
	length := 10
	orderArr := make([]int, length)
	for i := 0; i < length; i++ {
		orderArr[i] = rand.Intn(length)
	}
	if len(orderArr) <= 100 {
		fmt.Println(orderArr)
	}
	//orderArr := []int{1, 3, 2, 5, 4, 8, 9, 10000, 100000, 1, 3, 10, 50, 80, 100, 22, 42, 12}
	runtime(
		func() {
			// 选择排序
			//selectionSort(orderArr)
			// 冒泡排序
			//bubbleSort(orderArr)
			// 插入排序
			//insertionSort(orderArr)
			// 希尔排序
			//shellSort(orderArr)
			// 快速排序
			//quickSort(orderArr, 0, len(orderArr)-1)
			//quickSortTailCall(orderArr, 0, len(orderArr)-1)
			// 归并排序
			// mergeSort(orderArr, 0, len(orderArr)-1)
			// 堆排序见heap.go
			// 桶排序
			//bucketSort(orderArr)
			// 计数排序
			countingSort(orderArr)
		})

	if len(orderArr) <= 100 {
		fmt.Println(orderArr)
	}
}

// countingSort 计数排序
func countingSort(arr []int) {
	// 找到最大值和最小值
	max, min := arr[0], arr[0]
	for _, v := range arr {
		if max < v {
			max = v
		}
		if min > v {
			min = v
		}
	}

	// 计算计数数组的大小
	countingArr := make([]int, max-min+1)

	// 计算每个元素的个数
	for _, v := range arr {
		countingArr[v-min]++
	}

	// 计算每个元素的位置
	for i := 1; i < len(countingArr); i++ {
		countingArr[i] += countingArr[i-1]
	}

	// 临时数组
	tmp := make([]int, len(arr))
	for i := len(arr) - 1; i >= 0; i-- {
		idx := countingArr[arr[i]-min] - 1
		tmp[idx] = arr[i]
		countingArr[arr[i]-min]--
	}

	// 将临时数组的值赋值给原数组
	for i := 0; i < len(arr); i++ {
		arr[i] = tmp[i]
	}
}

// bucketSort 桶排序
func bucketSort(arr []int) {
	// 桶的最大值, 桶的最小值
	bucketMax, bucketMin := 0, 0
	// 找到最大值和最小值
	for _, v := range arr {
		if bucketMax < v {
			bucketMax = v
		}
		if bucketMin > v {
			bucketMin = v
		}
	}

	// 桶的数量
	bucketNum := bucketMax - bucketMin + 1

	// 桶的大小
	bucketSize := len(arr) / bucketNum
	// 桶的数量
	if len(arr)%bucketNum > 0 {
		bucketNum++
	}

	// 初始化桶
	buckets := make([][]int, bucketNum)
	for i := 0; i < bucketNum; i++ {
		buckets[i] = make([]int, 0)
	}

	// 将元素放入桶中
	for _, v := range arr {
		idx := (v - bucketMin) / bucketSize
		buckets[idx] = append(buckets[idx], v)
	}

	// 对每个桶进行排序, 利用插入排序对小数据量排序
	for _, bucket := range buckets {
		insertionSort(bucket)
	}

	// 将桶中的元素放回原数组
	idx := 0
	for _, bucket := range buckets {
		for _, v := range bucket {
			arr[idx] = v
			idx++
		}
	}
}

// merge 辅助合并函数，用于合并两个有序数组
// 左区间范围[left, mid]
// 右区间范围[mid+1, right]
func merge(arr []int, left, mid, right int) {
	// 先将左右两个区间的内容保存在临时tmp
	tmp := make([]int, right+1-left)
	for idx := 0; idx <= right-left; idx++ {
		tmp[idx] = arr[left+idx]
	}

	// 计算两个区间的左右下标信息
	leftStart, leftEnd := left-left, mid-left
	rightStart, rightEnd := mid+1-left, right-left

	// 记录左右区间开始合并有序数组
	i, j := leftStart, rightStart
	for idx := left; idx <= right; idx++ {
		// 如果左边的数组已经并入完成就执行
		if i > leftEnd {
			arr[idx] = tmp[j]
			j++
		} else if j > rightEnd {
			// 如果右边的数组已经并入完成
			arr[idx] = tmp[i]
			i++
		} else if i <= leftEnd && tmp[i] <= tmp[j] {
			// 如果左边的小于右边就先左边进入
			arr[idx] = tmp[i]
			i++
		} else {
			// 如果右边的小就右边的先进入
			arr[idx] = tmp[j]
			j++
		}
	}
}

// mergeSort 归并排序
func mergeSort(arr []int, left, right int) {
	if left >= right {
		return
	}
	mid := (left + right) / 2
	mergeSort(arr, left, mid)
	mergeSort(arr, mid+1, right)
	merge(arr, left, mid, right)
}

// quickSortTailCall 快排尾递归优化空间
func quickSortTailCall(arr []int, left, right int) {
	for left < right {
		pivot := partitionMedian(arr, left, right)
		if pivot-left < right-pivot {
			quickSortTailCall(arr, left, pivot-1)
			left = pivot + 1
		} else {
			quickSortTailCall(arr, pivot+1, right)
			right = pivot - 1
		}
	}
}

// quickSort 快速排序
func quickSort(arr []int, left, right int) {
	if left >= right {
		return
	}
	// 获取新的基准点
	//pivot := partition(arr, left, right)
	pivot := partitionMedian(arr, left, right)

	// 左
	//quickSort(arr, left, pivot-1)
	partitionMedian(arr, left, pivot-1)
	//quickSort(arr, pivot+1, right)
	partitionMedian(arr, pivot+1, right)
}

// partition 快速排序类-哨兵划分
func partition(arr []int, left, right int) int {
	i, j := left, right
	for i < j {
		for i < j && arr[j] >= arr[left] {
			j--
		}
		for i < j && arr[i] <= arr[left] {
			i++
		}
		arr[i], arr[j] = arr[j], arr[i]
	}
	arr[left], arr[i] = arr[i], arr[left]
	return i
}

func medianThree(arr []int, left, mid, right int) int {
	// 异或规则为 0 ^ 0 = 1 ^ 1 = 0, 0 ^ 1 = 1 ^ 0 = 1
	if (arr[left] < arr[mid]^arr[left]) != (arr[left] < arr[right]) {
		return left
	} else if (arr[mid] < arr[left]) != (arr[mid] < arr[right]) {
		return mid
	}
	return right
}

// partitionMedian 快速排序类-哨兵划分 三数取中值
func partitionMedian(arr []int, left, right int) int {
	i, j := left, right

	mid := medianThree(arr, left, (left+right)/2, right)
	arr[left], arr[mid] = arr[mid], arr[left]

	for i < j {
		for i < j && arr[j] >= arr[left] {
			j--
		}
		for i < j && arr[i] <= arr[left] {
			i++
		}
		arr[i], arr[j] = arr[j], arr[i]
	}
	arr[left], arr[i] = arr[i], arr[left]
	return i
}

// shellSort 希尔排序
func shellSort(arr []int) {
	gap := 1
	for len(arr)/3 > gap {
		gap = 3*gap + 1
	}

	for ; gap > 0; gap /= 3 {
		for i := gap; i < len(arr); i++ {
			base := arr[i]
			j := i
			for ; j >= gap && arr[j-gap] > base; j -= gap {
				arr[j] = arr[j-gap]
			}
			arr[j] = base
		}
	}
}

// insertionSort 插入排序
func insertionSort(arr []int) {
	for i := 1; i < len(arr); i++ {
		// 找到基准点
		base := arr[i]
		// 找到未排序的最后一个元素
		j := i - 1
		for ; j >= 0 && arr[j] > base; j-- {
			arr[j+1] = arr[j]
		}
		// 最后有一个j--找到了第一个小于等于base的值，所以需要+1找到第一个大于base的下标
		arr[j+1] = base
	}
}

// bubbleSort 冒泡排序
func bubbleSort(arr []int) {
	for i := len(arr) - 1; i > 0; i-- {
		for j := 0; j < i; j++ {
			if arr[j] > arr[j+1] {
				tmp := arr[j]
				arr[j] = arr[j+1]
				arr[j+1] = tmp
			}
		}
	}
}

// selectionSort 选择排序
func selectionSort(arr []int) {
	for i := 0; i < len(arr)-1; i++ {
		min := i
		for j := i + 1; j < len(arr); j++ {
			if arr[min] > arr[j] {
				min = j
			}
		}
		tmp := arr[i]
		arr[i] = arr[min]
		arr[min] = tmp
	}
}
