package main

import "fmt"

func main() {
	// 邻接矩阵无向图示例
	//adjMat()
	// 邻接链表无向图示例
	adjList()
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
	// 减一条边
	fmt.Printf("\t减掉边 3 -> 1\n")
	g.RemoveEdge(3, 1)
	g.Print()
}
