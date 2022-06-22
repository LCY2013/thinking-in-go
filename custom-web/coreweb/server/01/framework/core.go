package framework

import (
	"log"
	"net/http"
)

/*
自定义web框架核心实现文件
*/

// Core 框架核心结构
type Core struct {
	router map[string]ControllerHandler
}

// NewCore 初始化框架核心结构
func NewCore() *Core {
	return &Core{router: map[string]ControllerHandler{}}
}

// Get 设置Get请求路由信息
func (c *Core) Get(url string, handler ControllerHandler) {
	c.router[url] = handler
}

// ServeHTTP 框架核心结构实现Handler接口
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("core.serveHTTP")
	ctx := NewContext(request, response)

	// 路由选择器
	handler := c.router["foo"]

	log.Println("core.router")
	err := handler(ctx)
	if err != nil {
		return
	}
}
