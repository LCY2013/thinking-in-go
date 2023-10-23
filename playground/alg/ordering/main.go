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
	length := 100000000
	orderArr := make([]int, length)
	for i := 0; i < length; i++ {
		orderArr[i] = rand.Intn(length)
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
			quickSort(orderArr, 0, len(orderArr)-1)
			//quickSortTailCall(orderArr, 0, len(orderArr)-1)
		})

	//fmt.Println(orderArr)
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
