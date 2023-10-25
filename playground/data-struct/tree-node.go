package data_struct

import "fmt"

func (t *TreeNode) PrintTree() {
	printTreeHelper(t, nil, false)
	fmt.Println()
}

type Trunk struct {
	Prev *Trunk
	Str  string
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

func showTrunks(trunk *Trunk) {
	if trunk == nil {
		return
	}
	showTrunks(trunk.Prev)
	fmt.Printf("%s", trunk.Str)
}

// TreeNode 平衡二叉树节点
type TreeNode struct {
	Val    int
	Left   *TreeNode
	Right  *TreeNode
	Height int
}
