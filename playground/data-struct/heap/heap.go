package main

import (
	"fmt"
	"strconv"
)

type heapType int

const (
	MaxHeap heapType = iota // MaxHeap is a heap where the root node is the largest node in the tree.
	MinHeap                 // MinHeap is a heap where the root node is the smallest node in the tree.
)

const (
	UnlimitedSize = -1
)

// Heap is a binary tree with the following properties:
// 1. The value of each node is smaller/greater than or equal to the value of its parent, with the root being the largest node in the tree.
// 2. The binary tree is complete.
// 3. All nodes are filled from left to right.
// 4. The height of the tree is the smallest possible.
// 5. The root node is at index 0.
type Heap[T any] struct {
	data         []T
	size         int
	maxSize      int
	hType        heapType
	defaultValue T
}

type HeapOption[T any] func(*Heap[T])

func WithMaxSize[T any](maxSize int) HeapOption[T] {
	return func(h *Heap[T]) {
		h.maxSize = maxSize
	}
}

// NewTopMaxK returns a new max top heap.
func NewTopMaxK[T any](size int) *Heap[T] {
	return &Heap[T]{
		data:    make([]T, 0, 10),
		size:    0,
		hType:   MinHeap,
		maxSize: size,
	}
}

// NewTopMinK returns a new min top heap.
func NewTopMinK[T any](size int) *Heap[T] {
	return &Heap[T]{
		data:    make([]T, 0, 10),
		size:    0,
		hType:   MaxHeap,
		maxSize: size,
	}
}

// NewHeap default returns a new max top heap.
func NewHeap[T any](options ...HeapOption[T]) *Heap[T] {
	heap := &Heap[T]{
		data:    make([]T, 0, 10),
		size:    0,
		hType:   MaxHeap,
		maxSize: UnlimitedSize,
	}

	for _, option := range options {
		option(heap)
	}

	return heap
}

// NewMaxTopHeap returns a new max top heap.
func NewMaxTopHeap[T any](options ...HeapOption[T]) *Heap[T] {
	heap := &Heap[T]{
		data:    make([]T, 0, 10),
		size:    0,
		hType:   MaxHeap,
		maxSize: UnlimitedSize,
	}

	for _, option := range options {
		option(heap)
	}

	return heap
}

// NewMinTopHeap returns a new min top heap.
func NewMinTopHeap[T any](options ...HeapOption[T]) *Heap[T] {
	heap := &Heap[T]{
		data:    make([]T, 0, 10),
		size:    0,
		hType:   MinHeap,
		maxSize: UnlimitedSize,
	}

	for _, option := range options {
		option(heap)
	}

	return heap
}

// IsEmpty returns true if the heap is empty.
func (h *Heap[T]) IsEmpty() bool {
	return h.size == 0
}

// Size returns the number of elements in the heap.
func (h *Heap[T]) Size() int {
	return h.size
}

// right returns the index of the right child of the node at index idx.
func (h *Heap[T]) right(idx int) int {
	return 2*idx + 2
}

// left returns the index of the left child of the node at index idx.
func (h *Heap[T]) left(idx int) int {
	return 2*idx + 1
}

// parent returns the index of the parent of the node at index idx.
func (h *Heap[T]) parent(idx int) int {
	return (idx - 1) / 2
}

// Swap swaps the values at the given indices.
func (h *Heap[T]) Swap(idxA, idxB int) {
	h.data[idxA], h.data[idxB] = h.data[idxB], h.data[idxA]
}

// Peek returns the value of the root node of the heap.
func (h *Heap[T]) Peek() T {
	if h.size == 0 {
		return h.defaultValue
	}
	return h.data[0]
}

// TopK push a new value to the heap and pop the top value.
// if the heap size is less than maxSize, push the value to the heap and return the top value.
// if the heap size is equal to maxSize, push the value to the heap and return the top value.
func (h *Heap[T]) TopK(v T) T {
	// unlimited size
	if h.maxSize == UnlimitedSize {
		h.Push(v)
		return h.Peek()
	}

	// less than maxSize
	if h.size < h.maxSize {
		h.Push(v)
		return h.Peek()
	}

	top := h.Peek()

	var anyValue any
	anyValue = h.defaultValue

	switch anyValue.(type) {
	case int, int32, int64:
		pv, _ := strconv.ParseInt(fmt.Sprintf("%v", top), 10, 64)
		vv, _ := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
		if vv < pv && h.hType == MaxHeap ||
			vv > pv && h.hType == MinHeap {
			h.Pop()
			h.Push(v)
		}
	case float32, float64:
		pv, _ := strconv.ParseFloat(fmt.Sprintf("%v", top), 64)
		vv, _ := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
		if vv < pv && h.hType == MaxHeap ||
			vv > pv && h.hType == MinHeap {
			h.Pop()
			h.Push(v)
		}
	default:
		tv, ok := any(top).(interface{ Compare(t any) bool })
		if !ok {
			panic("top value not support compare")
		}
		vv, ok := any(v).(interface{ Compare(t any) bool })
		if !ok {
			panic("top value not support compare")
			return top
		}
		if vv != nil && tv != nil && tv.Compare(vv) {
			h.Pop()
			h.Push(v)
		}
	}

	return top
}

// Push adds a new value to the heap.
func (h *Heap[T]) Push(v T) {
	if h.maxSize != UnlimitedSize && h.size >= h.maxSize {
		return
	}
	if h.size == len(h.data) {
		h.data = append(h.data, v)
	} else {
		h.data[h.size] = v
	}

	h.size++
	h.shuffleUp(h.size - 1)
}

// shuffleUp moves the node at index idx up the heap until it is in the correct position.
func (h *Heap[T]) shuffleUp(idx int) {
	if idx == 0 {
		return
	}
	for {
		// find parent index
		p := h.parent(idx)
		if p == idx {
			break
		}
		var anyValue any
		anyValue = h.defaultValue

		switch anyValue.(type) {
		case int, int32, int64:
			var pv, cv int64
			if p < h.size {
				pv, _ = strconv.ParseInt(fmt.Sprintf("%v", h.data[p]), 10, 64)
			}
			if idx < h.size {
				cv, _ = strconv.ParseInt(fmt.Sprintf("%v", h.data[idx]), 10, 64)
			}
			// max top heap
			// if the parent is smaller than the child, swap them
			if cv > pv && h.hType == MaxHeap {
				h.Swap(p, idx)
				idx = p
				continue
			}
			// min top heap
			// if the parent is greater than the child, swap them
			if cv < pv && h.hType == MinHeap {
				h.Swap(p, idx)
				idx = p
				continue
			}
			return
		case float64, float32:
			var pv, cv int64
			if p < h.size {
				pv, _ = strconv.ParseInt(fmt.Sprintf("%v", h.data[p]), 10, 64)
			}
			if idx < h.size {
				cv, _ = strconv.ParseInt(fmt.Sprintf("%v", h.data[idx]), 10, 64)
			}
			// max top heap
			// if the parent is smaller than the child, swap them
			if cv > pv && h.hType == MaxHeap {
				h.Swap(p, idx)
				idx = p
				continue
			}
			// min top heap
			// if the parent is greater than the child, swap them
			if cv < pv && h.hType == MinHeap {
				h.Swap(p, idx)
				idx = p
				continue
			}
			return
		default:
			cV, ok := any(h.data[idx]).(interface{ Compare(t any) bool })
			if !ok {
				return
			}
			// custom compare
			if cV.Compare(h.data[p]) {
				h.Swap(p, idx)
				idx = p
			}
			return
		}
	}
}

// Pop removes the root node from the heap and returns its value.
func (h *Heap[T]) Pop() T {
	if h.size == 0 {
		return h.defaultValue
	}
	ret := h.data[0]
	h.Swap(0, h.size-1)
	h.data[h.size-1] = h.defaultValue
	h.size--
	h.shuffleDown(0)
	return ret
}

// shuffleDown moves the node at index idx down the heap until it is in the correct position.
func (h *Heap[T]) shuffleDown(idx int) {
	// loop until the node at idx is a leaf
	for {
		lIdx := h.left(idx)
		rIdx := h.right(idx)
		curIdx := idx

		anyValue := any(h.defaultValue)
		switch anyValue.(type) {
		case int, int32, int64:
			if h.hType == MaxHeap {
				// max top heap
				// if the parent is smaller than the child, swap them
				var lv, rv int64
				if lIdx < h.size {
					lv, _ = strconv.ParseInt(fmt.Sprintf("%v", h.data[lIdx]), 10, 64)
				}
				if rIdx < h.size {
					rv, _ = strconv.ParseInt(fmt.Sprintf("%v", h.data[rIdx]), 10, 64)
				}
				vv, _ := strconv.ParseInt(fmt.Sprintf("%v", h.data[curIdx]), 10, 64)
				//if lIdx < h.size && any(h.data[lIdx]).(int) > any(h.data[idx]).(int) {
				if lIdx < h.size && lv > vv {
					curIdx = lIdx
				}
				//if rIdx < h.size && any(h.data[rIdx]).(int) > any(h.data[idx]).(int) {
				vv, _ = strconv.ParseInt(fmt.Sprintf("%v", h.data[curIdx]), 10, 64)
				if rIdx < h.size && rv > vv {
					curIdx = rIdx
				}

				if curIdx == idx {
					return
				}

				h.Swap(idx, curIdx)
				idx = curIdx
			}
			if h.hType == MinHeap {
				// max top heap
				// if the parent is smaller than the child, swap them
				var lv, rv int64
				if lIdx < h.size {
					lv, _ = strconv.ParseInt(fmt.Sprintf("%v", h.data[lIdx]), 10, 64)
				}
				if rIdx < h.size {
					rv, _ = strconv.ParseInt(fmt.Sprintf("%v", h.data[rIdx]), 10, 64)
				}
				vv, _ := strconv.ParseInt(fmt.Sprintf("%v", h.data[curIdx]), 10, 64)
				//if lIdx < h.size && any(h.data[lIdx]).(int) < any(h.data[idx]).(int) {
				if lIdx < h.size && lv < vv {
					curIdx = lIdx
				}

				vv, _ = strconv.ParseInt(fmt.Sprintf("%v", h.data[curIdx]), 10, 64)
				//if rIdx < h.size && any(h.data[rIdx]).(int) < any(h.data[idx]).(int) {
				if rIdx < h.size && rv < vv {
					curIdx = rIdx
				}

				if curIdx == idx {
					return
				}

				h.Swap(idx, curIdx)
				idx = curIdx
			}
		case float32, float64:
			if h.hType == MaxHeap {
				// max top heap
				// if the parent is smaller than the child, swap them
				var lv, rv float64
				if lIdx < h.size {
					lv, _ = strconv.ParseFloat(fmt.Sprintf("%v", h.data[lIdx]), 64)
				}
				if rIdx < h.size {
					rv, _ = strconv.ParseFloat(fmt.Sprintf("%v", h.data[rIdx]), 64)
				}

				vv, _ := strconv.ParseFloat(fmt.Sprintf("%v", h.data[curIdx]), 64)
				//if lIdx < h.size && any(h.data[lIdx]).(float64) > any(h.data[idx]).(float64) {
				if lIdx < h.size && lv > vv {
					curIdx = lIdx
				}

				vv, _ = strconv.ParseFloat(fmt.Sprintf("%v", h.data[curIdx]), 64)
				//if rIdx < h.size && any(h.data[rIdx]).(float64) > any(h.data[idx]).(float64) {
				if rIdx < h.size && rv > vv {
					curIdx = rIdx
				}

				if curIdx == idx {
					return
				}

				h.Swap(idx, curIdx)
				idx = curIdx
			}
			if h.hType == MinHeap {
				// max top heap
				// if the parent is smaller than the child, swap them
				var lv, rv float64
				if lIdx < h.size {
					lv, _ = strconv.ParseFloat(fmt.Sprintf("%v", h.data[lIdx]), 64)
				}
				if rIdx < h.size {
					rv, _ = strconv.ParseFloat(fmt.Sprintf("%v", h.data[rIdx]), 64)
				}
				vv, _ := strconv.ParseFloat(fmt.Sprintf("%v", h.data[curIdx]), 64)
				//if lIdx < h.size && any(h.data[lIdx]).(float64) < any(h.data[idx]).(float64) {
				if lIdx < h.size && lv < vv {
					curIdx = lIdx
				}

				vv, _ = strconv.ParseFloat(fmt.Sprintf("%v", h.data[curIdx]), 64)
				//if rIdx < h.size && any(h.data[rIdx]).(float64) < any(h.data[idx]).(float64) {
				if rIdx < h.size && rv < vv {
					curIdx = rIdx
				}

				if curIdx == idx {
					return
				}

				h.Swap(idx, curIdx)
				idx = curIdx
			}
		default:
			var lv, rv interface{ Compare(t any) bool }
			if lIdx < h.size {
				lv = any(h.data[lIdx]).(interface{ Compare(t any) bool })
			}
			if rIdx < h.size {
				rv = any(h.data[rIdx]).(interface{ Compare(t any) bool })
			}
			pv := any(h.data[curIdx]).(interface{ Compare(t any) bool })
			if lv == nil && rv == nil {
				return
			}
			// max top heap
			// if the parent is smaller than the child, swap them
			if lv != nil && lv.Compare(pv) {
				curIdx = lIdx
			}

			pv = any(h.data[curIdx]).(interface{ Compare(t any) bool })
			if rv != nil && rv.Compare(pv) {
				curIdx = rIdx
			}

			if curIdx == idx {
				return
			}

			h.Swap(idx, curIdx)
			idx = curIdx
		}
	}
}

// PrintArray print array
func (h *Heap[T]) PrintArray() {
	fmt.Printf("[")
	if h != nil && h.size != 0 {
		for i := 0; i < h.size-1; i++ {
			fmt.Printf("%d, ", h.data[i])
		}
		fmt.Printf("%d", h.data[h.size-1])
	}
	fmt.Printf("]")
}

// TreeNode 平衡二叉树节点
type TreeNode struct {
	Val   any
	Left  *TreeNode
	Right *TreeNode
}

type Trunk struct {
	Prev *Trunk
	Str  string
}

func showTrunks(trunk *Trunk) {
	if trunk == nil {
		return
	}
	showTrunks(trunk.Prev)
	fmt.Printf("%s", trunk.Str)
}

/* Help to print a binary tree, hide more details */
func printTreeHelper(node *TreeNode, prev *Trunk, isLeft bool) {
	if node == nil {
		return
	}

	prevStr := "    "
	trunk := &Trunk{
		Prev: prev,
		Str:  prevStr,
	}
	printTreeHelper(node.Right, trunk, true)
	if prev == nil {
		trunk.Str = "———"
	} else if isLeft {
		trunk.Str = "/———"
		prevStr = "   |"
	} else {
		trunk.Str = "\\———"
		prev.Str = prevStr
	}
	showTrunks(trunk)
	fmt.Printf("%d\n", node.Val)

	if prev != nil {
		prev.Str = prevStr
	}
	trunk.Str = "   |"

	printTreeHelper(node.Left, trunk, false)
}

func printTree(node *TreeNode) {
	printTreeHelper(node, nil, false)
}

func (h *Heap[T]) printHeap() {
	fmt.Printf("堆的数组表示：")
	h.PrintArray()
	fmt.Printf("堆的树状表示：\n")
	root := h.arrToTree()
	printTree(root)
}

func (h *Heap[T]) arrToTree() *TreeNode {
	root := &TreeNode{
		Val: h.data[0],
	}
	queue := []*TreeNode{root}
	for i := 0; i < h.size; {
		node := queue[0]
		queue = queue[1:]
		i++
		if i < h.size {
			node.Left = &TreeNode{
				Val: h.data[i],
			}
			queue = append(queue, node.Left)
		}
		i++
		if i < h.size {
			node.Right = &TreeNode{
				Val: h.data[i],
			}
			queue = append(queue, node.Right)
		}
	}
	return root
}

func main() {
	//maxTopHeapCase()
	//minTopHeapCase()
	topKCase()
}

func topKCase() {
	heap := NewTopMaxK[int](3)
	/* 初始化堆 */
	// 初始化大顶堆
	top := heap.TopK(1)
	fmt.Printf("\ntopk堆顶元素 1 入堆后，返回：%d\n", top)
	heap.printHeap()
	top = heap.TopK(2)
	fmt.Printf("\ntopk堆顶元素 2 入堆后，返回：%d\n", top)
	heap.printHeap()
	top = heap.TopK(3)
	fmt.Printf("\ntopk堆顶元素 3 入堆后，返回：%d\n", top)
	heap.printHeap()
	top = heap.TopK(4)
	fmt.Printf("\ntopk堆顶元素 4 入堆后，返回：%d\n", top)
	heap.printHeap()
	top = heap.TopK(5)
	fmt.Printf("\ntopk堆顶元素 5 入堆后，返回：%d\n", top)
	heap.printHeap()
	top = heap.TopK(10)
	fmt.Printf("\ntopk堆顶元素 10 入堆后，返回：%d\n", top)
	heap.printHeap()
	top = heap.TopK(1)
	fmt.Printf("\ntopk堆顶元素 1 入堆后，返回：%d\n", top)
	heap.printHeap()
}

func minTopHeapCase() {
	heap := NewMinTopHeap[int]()
	/* 初始化堆 */
	// 初始化大顶堆
	heapNums := []int{9, 8, 6, 6, 7, 5, 2, 1, 4, 3, 6, 2}
	for _, num := range heapNums {
		heap.Push(num)
	}
	fmt.Printf("输入数组并建堆后\n")
	heap.printHeap()

	for !heap.IsEmpty() {
		pop := heap.Pop()
		fmt.Printf("\n堆顶元素 %d 出堆后\n", pop)
		heap.printHeap()
	}
}

func maxTopHeapCase() {
	heap := NewMaxTopHeap[int]()
	/* 初始化堆 */
	// 初始化大顶堆
	heapNums := []int{9, 8, 6, 6, 7, 5, 2, 1, 4, 3, 6, 2}
	for _, num := range heapNums {
		heap.Push(num)
	}
	fmt.Printf("输入数组并建堆后\n")
	heap.printHeap()

	topHeap := NewMaxTopHeap[int]()
	nums := []int{9, 8, 6, 6, 7, 5, 2, 1, 4, 3, 6, 2}
	topHeap.size = len(nums)
	for i := range nums {
		topHeap.data = append(topHeap.data, nums[i])
	}
	fmt.Printf("输入指定堆后\n")
	topHeap.printHeap()

	/* 获取堆顶元素 */
	fmt.Printf("\n堆顶元素为 %d\n", topHeap.Peek())

	/* 元素入堆 */
	topHeap.Push(7)
	fmt.Printf("\n元素 7 入堆后\n")
	topHeap.printHeap()

	/* 堆顶元素出堆 */
	top := topHeap.Pop()
	fmt.Printf("\n堆顶元素 %d 出堆后\n", top)
	topHeap.printHeap()

	/* 获取堆大小 */
	fmt.Printf("\n堆元素数量为 %d\n", topHeap.Size())

	/* 判断堆是否为空 */
	fmt.Printf("\n堆是否为空 %v\n", topHeap.IsEmpty())

	for !topHeap.IsEmpty() {
		pop := topHeap.Pop()
		fmt.Printf("\n堆顶元素 %d 出堆后\n", pop)
		topHeap.printHeap()
	}
}
