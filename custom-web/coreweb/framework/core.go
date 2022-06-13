package framework

import "net/http"

/*
自定义web框架核心实现文件
*/

// Core 框架核心结构
type Core struct {
}

// NewCore 初始化框架核心结构
func NewCore() *Core {
	return &Core{}
}

// ServeHTTP 框架核心结构实现Handler接口
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// TODO
}
