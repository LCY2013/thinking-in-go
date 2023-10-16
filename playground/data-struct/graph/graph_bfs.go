package main

import "fmt"

// Bfs 广度优先搜索
// 广度优先遍历的序列是否唯一？
// 不唯一。广度优先遍历只要求按“由近及远”的顺序遍历，而多个相同距离的顶点的遍历顺序是允许被任意打乱的。
func Bfs(list *AdjList) {
	// 初始化一个队列
	var queue []int
	// 初始化一个map，用于记录已经访问过的顶点
	visited := make(map[int]struct{})

	for _, v := range list.vertexList {
		if _, ok := visited[v]; ok {
			continue
		}
		queue = append(queue, v)
		visited[v] = struct{}{}
		for len(queue) > 0 {
			// 出队列
			out := queue[0]
			queue = queue[1:]
			fmt.Printf("%d\t", out)
			// 获取顶点的所有邻接点
			for _, jv := range list.vertexMap[out] {
				if _, ok := visited[jv]; ok {
					continue
				}
				queue = append(queue, jv)
				visited[jv] = struct{}{}
			}
		}
	}

}

// bfs 广度优先搜索
func bfs() {
	// 初始化一个无向图
	graphAdjList := NewGraphAdjList()
	// 初始化顶点
	for i := 0; i < 10; i++ {
		graphAdjList.AddVertex(i)
	}
	// 初始化边
	graphAdjList.AddEdge(0, 1)
	graphAdjList.AddEdge(0, 3)
	graphAdjList.AddEdge(1, 2)
	graphAdjList.AddEdge(1, 4)
	graphAdjList.AddEdge(2, 5)
	graphAdjList.AddEdge(3, 4)
	graphAdjList.AddEdge(3, 6)
	graphAdjList.AddEdge(4, 5)
	graphAdjList.AddEdge(4, 7)
	graphAdjList.AddEdge(5, 8)
	graphAdjList.AddEdge(6, 7)
	graphAdjList.AddEdge(7, 8)
	fmt.Printf("\n初始化后，图为:\n")
	graphAdjList.Print()

	// 广度优先搜索
	fmt.Printf("\n广度优先搜索结果:\n")
	Bfs(graphAdjList)
}
