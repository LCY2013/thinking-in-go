package main

import "fmt"

// 现在数据库有一张表，用来存储一个多叉树，id为主键，pid 表示父节点的 id，已知 "-1" 表示根节点，现在要求打印出从根节点到每个子节点的路径（可以是无序的）。
//
// | id      | pid    |
// |---------|--------|
// | "A"     | "-1"   |
// | "A-1"   | "A"    |
// | "A-2"   | "A"    |
// | "A-3"   | "A"    |
// | "A-2-1" | "A-2"  |
// | "A-2-2" | "A-2"  |
// | "A-2-3" | "A-2"  |
//
// Input: [
//   {
//       "id": "A",
//       "pid": "-1"
//   },
//   {
//       "id": "A-1",
//       "pid": "A"
//   },
//   {
//       "id": "A-2",
//       "pid": "A"
//   },
//   {
//       "id": "A-3",
//       "pid": "A"
//   },
//   {
//       "id": "A-2-1",
//       "pid": "A-2"
//   },
//   {
//       "id": "A-2-2",
//       "pid": "A-2"
//   },
//   {
//       "id": "A-2-3",
//       "pid": "A-2"
//   }
// ]
// Output: [
//   "/A",
//   "/A/A-1",
//   "/A/A-2",
//   "/A/A-3",
//   "/A/A-2/A-2-1",
//   "/A/A-2/A-2-2",
//   "/A/A-2/A-2-3",
// ]

type Node struct {
	ID       string
	PID      string
	Children []*Node
}

func main() {
	nodes := []*Node{
		{
			ID:  "A",
			PID: "-1",
		},
		{
			ID:  "A-1",
			PID: "A",
		},
		{
			ID:  "A-2",
			PID: "A",
		},
		{
			ID:  "A-3",
			PID: "A",
		},
		{
			ID:  "A-2-1",
			PID: "A-2",
		},
		{
			ID:  "A-2-2",
			PID: "A-2",
		},
		{
			ID:  "A-2-3",
			PID: "A-2",
		},
	}

	// 时间复杂度 O(n)
	root := buildTree(nodes)
	prefix := "/"
	fmt.Println(prefix + root.ID)
	printTree(root, prefix)
}

func printTree(root *Node, prefix string) {
	for _, child := range root.Children {
		fmt.Println("/" + child.ID)
	}
	for _, child := range root.Children {
		printTree(child, prefix+root.ID+prefix)
	}
}

func buildTree(nodes []*Node) *Node {
	var root *Node

	nodesMap := make(map[string]*Node)
	for _, node := range nodes {
		nodesMap[node.ID] = node
	}

	for _, node := range nodes {
		if node.PID == "-1" {
			root = node
			nodesMap[node.ID] = node
		} else {
			pNode := nodesMap[node.PID]
			pNode.Children = append(pNode.Children, node)
		}
	}

	return root
}
