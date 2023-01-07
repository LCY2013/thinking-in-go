package container

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var ginEngine map[string]*gin.Engine

type Serve struct {
	ServeName string `json:"serveName"`
	ServePort int    `json:"servePort"`
}

// BuildMultipleGinServe 构建多个服务
func BuildMultipleGinServe(serves []*Serve) []*Server {
	if serves == nil || len(serves) == 0 {
		return nil
	}

	var (
		serveLen = len(serves)
		servers  = make([]*Server, serveLen)
	)

	if ginEngine == nil {
		ginEngine = make(map[string]*gin.Engine, serveLen)
	}

	for idx, serve := range serves {
		engine := gin.Default()
		servers[idx] = NewHandlerServer(serve.ServeName, fmt.Sprintf(":%d", serve.ServePort), engine)
		ginEngine[serve.ServeName] = engine
	}

	return servers
}

// GinEngineByServeName 通过serve name 获取engine
func GinEngineByServeName(name string) (*gin.Engine, bool) {
	engine, ok := ginEngine[name]
	return engine, ok
}
