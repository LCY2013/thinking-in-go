package main

import "fmt"

// AdjMat GraphAdjMat 邻接矩阵的方式构建无向图信息
type AdjMat struct {
	// 顶点列表
	vertexList []int
	// 邻接矩阵，元素代表“边”，索引代表“顶点索引”
	edge [][]int
}

// NewGraphAdjMat 邻接矩阵的方式构建无向图信息
func NewGraphAdjMat() *AdjMat {
	return &AdjMat{
		vertexList: make([]int, 0),
		edge:       make([][]int, 0),
	}
}

// Size 获取顶点数量
func (g *AdjMat) Size() int {
	return len(g.vertexList)
}

// VertexIndex 获取到对应顶点的索引信息
func (g *AdjMat) VertexIndex(v int) (idx int, exist bool) {
	for vIdx, ver := range g.vertexList {
		if ver == v {
			idx = vIdx
			exist = true
			break
		}
	}
	return
}

// AddVertex 添加顶点
func (g *AdjMat) AddVertex(v int) {
	// 添加顶点信息
	g.vertexList = append(g.vertexList, v)
	// 初始化边信息
	for idx := range g.edge {
		g.edge[idx] = append(g.edge[idx], 0)
	}
	g.edge = append(g.edge, make([]int, g.Size()))
}

// RemoveVertex 移除顶点
func (g *AdjMat) RemoveVertex(v int) {
	// 找到某个顶点索引
	vIdx, exist := g.VertexIndex(v)

	if !exist {
		return
	}

	// 移除顶点
	g.vertexList = append(g.vertexList[:vIdx], g.vertexList[vIdx+1:]...)
	// 移除边信息
	g.edge = append(g.edge[:vIdx], g.edge[vIdx+1:]...)
	for idx := range g.edge {
		g.edge[idx] = append(g.edge[idx][:vIdx], g.edge[idx][vIdx+1:]...)
	}
}

// AddEdge 添加边
func (g *AdjMat) AddEdge(iv, jv int) {
	i, exist := g.VertexIndex(iv)
	if !exist {
		return
	}
	j, exist := g.VertexIndex(jv)
	if !exist {
		return
	}

	if j == i || i >= g.Size() || j >= g.Size() {
		return
	}
	g.edge[i][j] = 1
	g.edge[j][i] = 1
}

// RemoveEdge 移除边
func (g *AdjMat) RemoveEdge(iv, jv int) {
	i, exist := g.VertexIndex(iv)
	if !exist {
		return
	}
	j, exist := g.VertexIndex(jv)
	if !exist {
		return
	}

	if j == i || i >= g.Size() || j >= g.Size() {
		return
	}
	g.edge[i][j] = 0
	g.edge[j][i] = 0
}

// Print 打印邻接矩阵
func (g *AdjMat) Print() {
	fmt.Printf("\t顶点列表 = %v\n", g.vertexList)
	fmt.Printf("\t邻接矩阵 = \n")
	for i := range g.edge {
		fmt.Printf("\t\t\t%v\n", g.edge[i])
	}
}

func main() {
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
