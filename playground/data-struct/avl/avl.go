package main

import (
	"fmt"
)

func max(a int, b int) int {
	if a >= b {
		return a
	}
	return b
}

// AvlTree 平衡二叉搜索树
type AvlTree struct {
	Root *AvlNode
}

// Insert 树插入节点
func (t *AvlTree) Insert(v int) {
	t.Root = t.Root.Insert(v)
}

// Remove 树删除节点
func (t *AvlTree) Remove(v int) {
	t.Root = t.Root.Remove(v)
}

// Search 树查找节点
func (t *AvlTree) Search(v int) *AvlNode {
	return t.Root.Search(v)
}

// AvlNode 平衡二叉树节点
type AvlNode struct {
	Val    int
	Left   *AvlNode
	Right  *AvlNode
	Height int
}

// height avlNode 高度
func (n *AvlNode) height() int {
	if n == nil {
		return -1
	}
	return n.Height
}

// updateHeight 更新avlNode 高度
func (n *AvlNode) updateHeight() {
	n.Height = max(n.Left.height(), n.Right.height()) + 1
}

// balanceFactor 平衡因子
func (n *AvlNode) balanceFactor() int {
	if n == nil {
		return -1
	}
	return n.Left.height() - n.Right.height()
}

// rightRotate 右旋
func (n *AvlNode) rightRotate() *AvlNode {
	var child, grantChild *AvlNode
	if n == nil || n.Left == nil {
		return n
	}
	child = n.Left
	grantChild = child.Right
	child.Right = n
	n.Left = grantChild
	n.updateHeight()
	child.updateHeight()
	return child
}

// leftRotate 左旋
func (n *AvlNode) leftRotate() *AvlNode {
	var child, grantChild *AvlNode
	if n == nil || n.Right == nil {
		return n
	}
	child = n.Right
	grantChild = child.Left
	child.Left = n
	n.Right = grantChild
	n.updateHeight()
	child.updateHeight()
	return child
}

// rotate 旋转操作
func (n *AvlNode) rotate() *AvlNode {
	// 左偏树
	if n.balanceFactor() > 1 {
		// 左左树
		if n.Left.balanceFactor() >= 0 {
			// 直接右旋
			return n.rightRotate()
		} else { // 左右树
			// 先左转再右转
			n.Left = n.Left.leftRotate()
			return n.rightRotate()
		}
	}
	if n.balanceFactor() < -1 { // 右偏树
		// 右右树
		if n.Right.balanceFactor() <= 0 {
			// 直接左转
			return n.leftRotate()
		} else { // 右左树
			// 先右转再左转
			n.Right = n.Right.rightRotate()
			return n.leftRotate()
		}
	}
	return n
}

// Insert 插入节点
func (n *AvlNode) Insert(v int) *AvlNode {
	// 如果当前节点不存在，直接插入
	if n == nil {
		return &AvlNode{
			Val: v,
		}
	}

	// 插入到指定位置
	// 递归：递 为寻找对应节点， 归为链接新建节点和上级节点
	if v < n.Val {
		n.Left = n.Left.Insert(v)
	} else if v > n.Val {
		n.Right = n.Right.Insert(v)
	} else {
		return n
	}
	// 更新高度
	n.updateHeight()
	// 旋转，使得子树平衡
	return n.rotate()
}

// copyTo 将source内容复制给n
func (n *AvlNode) copyTo(source *AvlNode) {
	n.Val = source.Val
	n.Left = source.Left
	n.Right = source.Right
	n.Height = source.Height
}

// Remove 删除节点
func (n *AvlNode) Remove(v int) *AvlNode {
	// 如果当前节点不存在，直接插入
	if n == nil {
		return n
	}
	// 删除到指定位置
	// 递归：递 为寻找对应节点， 归为链接新建节点和上级节点
	if v < n.Val {
		n.Left = n.Left.Remove(v)
	} else if v > n.Val {
		n.Right = n.Right.Remove(v)
	} else {
		// 区分删除节点的情况
		// 1. 叶子节点
		if n.Left == nil && n.Right == nil {
			return nil
		} else if n.Left == nil || n.Right == nil { // 2. 只有一个子节点
			child := n.Left
			if n.Right != nil {
				child = n.Right
			}
			n.copyTo(child)
		} else {
			// 3. 有两个子节点
			// 找到右子树最小节点，或者找到左子树最大节点
			// 这里找到右子树最小节点
			rightMinNode := n.Right
			for rightMinNode.Left != nil {
				rightMinNode = rightMinNode.Left
			}
			// 删除右子树最小节点
			n.Right = n.Right.Remove(rightMinNode.Val)
			// 将右子树最小节点的值赋值给当前节点
			n.Val = rightMinNode.Val
		}
	}

	// 更新节点高度
	n.updateHeight()
	// 旋转，平衡子树
	return n.rotate()
}

// SearchRecursion 递归实现查找节点
func (n *AvlNode) SearchRecursion(v int) *AvlNode {
	if n == nil {
		return nil
	}

	if v < n.Val {
		return n.Left.SearchRecursion(v)
	} else if v > n.Val {
		return n.Right.SearchRecursion(v)
	} else {
		return n
	}
}

// Search 循环实现查找节点
func (n *AvlNode) Search(v int) *AvlNode {
	if n == nil {
		return nil
	}

	cur := n
	for cur != nil {
		if v > cur.Val { // 查询的值大于当前节点的值，去右边查询
			cur = cur.Right
		} else if v < cur.Val { // 查询的值小于当前节点的值，去左边查询
			cur = cur.Left
		} else {
			// 找到或者没找到直接退出循环
			break
		}
	}

	return cur
}

func (t *AvlTree) insertAndPrint(v int) {
	t.Insert(v)
	fmt.Printf("\n插入节点 %d 后，AVL 树为 \n", v)
	t.printTree()
}

func (t *AvlTree) removeAndPrint(v int) {
	t.Remove(v)
	fmt.Printf("\n移除节点 %d 后，AVL 树为 \n", v)
	t.printTree()
}

func (t *AvlTree) printTree() {
	printTreeHelper(t.Root, nil, false)
	fmt.Println()
}

type Trunk struct {
	Prev *Trunk
	Str  string
}

/* Help to print a binary tree, hide more details */
func printTreeHelper(node *AvlNode, prev *Trunk, isLeft bool) {
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

func showTrunks(trunk *Trunk) {
	if trunk == nil {
		return
	}
	showTrunks(trunk.Prev)
	fmt.Printf("%s", trunk.Str)
}

func main() {
	/* 初始化空 AVL 树 */
	tree := &AvlTree{}
	/* 插入节点 */
	// 请关注插入节点后，AVL 树是如何保持平衡的
	tree.insertAndPrint(1)
	tree.insertAndPrint(2)
	tree.insertAndPrint(3)
	tree.insertAndPrint(4)
	tree.insertAndPrint(5)
	tree.insertAndPrint(8)
	tree.insertAndPrint(7)
	tree.insertAndPrint(9)
	tree.insertAndPrint(10)
	tree.insertAndPrint(6)

	/* 插入重复节点 */
	tree.insertAndPrint(7)

	/* 删除节点 */
	// 请关注删除节点后，AVL 树是如何保持平衡的
	tree.removeAndPrint(8) // 删除度为 0 的节点
	tree.removeAndPrint(5) // 删除度为 1 的节点
	tree.removeAndPrint(4) // 删除度为 2 的节点

	/* 查询节点 */
	searchNode := tree.Search(7)
	fmt.Printf("\n查找到的节点对象节点值 = %d \n", searchNode.Val)
}
