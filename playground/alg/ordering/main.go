package main

import "fmt"

func main() {
	orderArr := []int{1, 3, 2, 5, 4, 8, 9, 1, 3, 10, 50, 22, 42, 12}
	// 选择排序
	//selectionSort(orderArr)
	// 冒泡排序
	//bubbleSort(orderArr)
	// 插入排序
	//insertionSort(orderArr)

	fmt.Println(orderArr)
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

/*
for i := 1; i < len(arr); i++ {
		// 找到基准点
		base := arr[i]
		// 找到排序区域的最后一个位置
		j := i - 1
		for ; j > 0 && arr[j] > base; j-- {
			arr[j+1] = arr[j]
		}
		// 将base插入正确位置
		arr[j+1] = base
	}
*/
