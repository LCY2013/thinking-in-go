package main

import "fmt"

func main() {
	// 邻接矩阵无向图示例
	//adjMat()
	// 邻接链表无向图示例
	//adjList()
	// 广度优先遍历
	//bfs()
	// 深度优先遍历
	dfs()
}

func adjList() {
	// 构造一个无向图
	g := NewGraphAdjList()
	vertices := []int{1, 3, 2, 5, 4}
	for _, i := range vertices {
		g.AddVertex(i)
	}
	g.Print()

	// 加一条边
	fmt.Printf("添加边 1 -> 3\n")
	g.AddEdge(1, 3)
	g.Print()

	// 删除一个顶点后
	fmt.Printf("\t删除顶点 3\n")
	g.RemoveVertex(3)
	g.Print()

	// 减一条边
	fmt.Printf("减掉边 3 -> 1\n")
	g.RemoveEdge(3, 1)
	g.Print()
}

func adjMat() {
	// 构造一个无向图
	g := NewGraphAdjMat()
	vertices := []int{1, 3, 2, 5, 4}
	for _, i := range vertices {
		g.AddVertex(i)
	}
	g.Print()

	// 加一条边
	fmt.Printf("\t添加边 1 -> 3\n")
	g.AddEdge(1, 3)
	g.Print()

	// 删除一个顶点后
	fmt.Printf("\t删除顶点 3\n")
	g.RemoveVertex(3)
	g.Print()

	// 减一条边
	fmt.Printf("\t减掉边 3 -> 1\n")
	g.RemoveEdge(3, 1)
	g.Print()
}
