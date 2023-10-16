package main

import "fmt"

// Dfs 深度优先搜索
// 1. 从初始访问节点出发，初始访问节点可能有多个邻接节点，深度优先搜索的策略就是首先访问第一个邻接节点，然后再以这个被访问的邻接节点作为初始节点，访问它的第一个邻接节点，可以这样理解：每次都在访问完当前节点后首先访问当前节点的第一个邻接节点。
// 2. 我们可以看到，这样的访问策略是优先往纵向挖掘深入，而不是对一个节点的所有邻接节点进行横向访问。
// 3. 显然，深度优先搜索是一个递归的过程。
// 4. 深度优先搜索会有一个问题：可能会陷入“死胡同”，导致无法继续往下搜索。解决这个问题的方法是：每次访问到一个节点的时候，不是立即访问该节点的所有邻接节点，而是先访问该节点的一个邻接节点，而且是选择最先访问到该节点的邻接节点。
// 5. 深度优先搜索算法的运行时间不一定就是 O(E)，而是 O(E+V)，其中，E 表示边的个数，V 表示顶点的个数。这个时间复杂度怎么算出来的呢？我们可以这样理解：每条边和每个顶点都会被访问一次，因此时间复杂度就是 O(E+V)。
// 6. 深度优先搜索算法的空间消耗主要取决于递归调用的栈空间，栈的大小和递归调用的深度成正比，所以空间复杂度是 O(V)。
// 7. 深度优先搜索算法的应用场景有哪些呢？我们可以通过深度优先搜索算法来解决这样一个问题：给你一张地图，找出两个地点之间的路线。这个问题怎么用深度优先搜索算法来解决呢？我们可以把地图抽象成一个无向图，每个交叉路口看成一个顶点，每条路看成一条边，然后，我们从一个地点出发，假设是从顶点 A 出发，然后深度优先搜索地遍历地图，直到找到终点或者遍历整张地图。
// 8. 深度优先搜索算法有一个缺点，它找到的路径不一定是最短路径。为什么呢？我们来看下面这个例子。我们从图中的顶点 1 出发，深度优先搜索的结果是 1->2->4->5，但是从顶点 1 到顶点 5 的最短路径是 1->3->5。
func Dfs(list *AdjList, visited map[int]struct{}, vertex int) {
	if len(list.vertexMap) == 0 {
		return
	}
	// 前序遍历
	if _, ok := visited[vertex]; ok {
		return
	}
	visited[vertex] = struct{}{}

	fmt.Printf("%d\t", vertex)
	for _, jv := range list.vertexMap[vertex] {
		Dfs(list, visited, jv)
	}
}

// dfs 深度优先搜索
func dfs() {
	// 初始化一个无向图
	graphAdjList := NewGraphAdjList()
	// 初始化顶点
	for i := 0; i < 7; i++ {
		graphAdjList.AddVertex(i)
	}
	// 初始化边
	graphAdjList.AddEdge(0, 1)
	graphAdjList.AddEdge(0, 3)
	graphAdjList.AddEdge(1, 2)
	graphAdjList.AddEdge(2, 5)
	graphAdjList.AddEdge(5, 4)
	graphAdjList.AddEdge(5, 6)
	fmt.Printf("\n初始化后，图为:\n")
	graphAdjList.Print()

	// 广度优先搜索
	fmt.Printf("\n深度度优先搜索结果:\n")
	Dfs(graphAdjList, map[int]struct{}{}, graphAdjList.vertexList[0])
}
