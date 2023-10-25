package main

import (
	"fmt"
	data_struct "pg/data-struct"
)

func main() {
	preorder := []int{3, 9, 2, 1, 7}
	inorder := []int{9, 3, 1, 2, 7}
	fmt.Printf("前序遍历 = ")
	printArray(preorder)
	fmt.Printf("中序遍历 = ")
	printArray(inorder)

	root := buildTree(preorder, inorder)
	fmt.Printf("构建的二叉树为：\n")
	root.PrintTree()
}

// printArray 打印数组
func printArray(preorder []int) {
	for _, v := range preorder {
		fmt.Printf("%d ", v)
	}
	fmt.Println()
}

// buildTree 给定一个二叉树的前序遍历 preorder 和中序遍历 inorder ，请从中构建二叉树，返回二叉树的根节点
func buildTree(preorder []int, inorder []int) *data_struct.TreeNode {
	// 构建中序遍历的 map，获得左右子树索引信息
	inorderMap := make(map[int]int)
	for idx, v := range inorder {
		inorderMap[v] = idx
	}

	return builderTreeHelper(preorder, inorderMap, 0, 0, len(preorder)-1)
}

func builderTreeHelper(preorder []int, inorderMap map[int]int, preIdx, left, right int) *data_struct.TreeNode {
	if left > right {
		return nil
	}
	// 构建跟节点
	root := &data_struct.TreeNode{
		Val: preorder[preIdx],
	}

	// 构建左子树
	root.Left = builderTreeHelper(preorder, inorderMap, preIdx+1, left, inorderMap[root.Val]-1)
	// 构建右子树
	root.Right = builderTreeHelper(preorder, inorderMap, preIdx+1+inorderMap[root.Val]-left, inorderMap[root.Val]+1, right)
	return root
}
