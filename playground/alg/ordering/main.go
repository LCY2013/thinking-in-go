package main

import "fmt"

func main() {
	orderArr := []int{1, 3, 2, 5, 4, 8, 9, 10000, 100000, 1, 3, 10, 50, 80, 100, 22, 42, 12}
	// 选择排序
	//selectionSort(orderArr)
	// 冒泡排序
	//bubbleSort(orderArr)
	// 插入排序
	//insertionSort(orderArr)
	// 希尔排序
	shellSort(orderArr)
	// 快速排序
	quickSort(orderArr, 0, len(orderArr)-1)

	fmt.Println(orderArr)
}

// quickSort 快速排序
func quickSort(arr []int, start, end int) {

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
		for ; j > 0 && arr[j] > base; j-- {
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
