package main

import (
	"fmt"
	"strconv"
	"strings"
)

// AdjList GraphAdjList 邻接链表的方式构建无向图信息
type AdjList struct {
	// 顶点列表, 包含边信息
	vertexList []int
	vertexMap  map[int][]int
}

// NewGraphAdjList 邻接链表的方式构建无向图信息
func NewGraphAdjList() *AdjList {
	return &AdjList{
		vertexMap: make(map[int][]int),
	}
}

// Size 获取顶点数量
func (g *AdjList) Size() int {
	return len(g.vertexMap)
}

// AddVertex 添加顶点信息
func (g *AdjList) AddVertex(v int) {
	g.vertexMap[v] = make([]int, 0)
	g.vertexList = append(g.vertexList, v)
}

// RemoveVertex 移除顶点信息
func (g *AdjList) RemoveVertex(v int) {
	delete(g.vertexMap, v)
	for idx, ver := range g.vertexList {
		if ver == v {
			g.vertexList = append(g.vertexList[:idx], g.vertexList[idx+1:]...)
			break
		}
	}
}

// AddEdge 添加边信息
func (g *AdjList) AddEdge(iv, jv int) {
	g.vertexMap[iv] = append(g.vertexMap[iv], jv)
	g.vertexMap[jv] = append(g.vertexMap[jv], iv)
}

// RemoveEdge 移除边信息
func (g *AdjList) RemoveEdge(iv, jv int) {
	for idx, v := range g.vertexMap[iv] {
		if v == jv {
			g.vertexMap[iv] = append(g.vertexMap[iv][:idx], g.vertexMap[iv][idx+1:]...)
			break
		}
	}
	for idx, v := range g.vertexMap[jv] {
		if v == iv {
			g.vertexMap[jv] = append(g.vertexMap[jv][:idx], g.vertexMap[jv][idx+1:]...)
			break
		}
	}
}

// Print 打印邻接表
func (g *AdjList) Print() {
	var builder strings.Builder
	fmt.Printf("邻接表 = \n")
	for _, v := range g.vertexList {
		list := g.vertexMap[v]
		builder.WriteString("\t\t" + strconv.Itoa(v) + ": ")
		for _, l := range list {
			builder.WriteString(strconv.Itoa(l) + " ")
		}
		fmt.Println(builder.String())
		builder.Reset()
	}
}
