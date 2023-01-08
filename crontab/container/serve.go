package container

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

var ginEngine map[string]*gin.Engine

type Serve struct {
	ServeName    string `json:"serveName"`
	ServePort    int    `json:"servePort"`
	ReadTimeout  int    `json:"readTimeout"`
	WriteTimeout int    `json:"writeTimeout"`
}

// BuildMultipleGinServe 构建多个服务
func BuildMultipleGinServe(serves []Serve) []*Server {
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
		//engine := gin.Default()
		engine := gin.New()

		servers[idx] = NewHandlerServer(serve.ServeName,
			fmt.Sprintf(":%d", serve.ServePort),
			engine,
			WithHandleServerReadTimeout(serve.ReadTimeout),
			WithHandleServerWriteTimeout(serve.WriteTimeout))

		ginEngine[serve.ServeName] = engine

		engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			// 你的自定义格式
			/*return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)*/
			log.WithFields(log.Fields{
				"ClientIP":     param.ClientIP,
				"Timestamp":    param.TimeStamp.Format(time.RFC3339Nano),
				"Method":       param.Method,
				"Path":         param.Path,
				"Proto":        param.Request.Proto,
				"StatusCode":   param.StatusCode,
				"Latency":      fmt.Sprintf("%dns", param.Latency.Nanoseconds()),
				"UserAgent":    param.Request.UserAgent(),
				"ErrorMessage": param.ErrorMessage,
			}).Log(log.InfoLevel)
			return ""
		}))

		engine.Use(gin.Recovery())
	}

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		//log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
		log.WithFields(log.Fields{
			"httpMethod":  httpMethod,
			"uri":         absolutePath,
			"handlerName": handlerName,
			"rt":          nuHandlers,
		}).Log(log.InfoLevel)
	}

	return servers
}

// GinEngineByServeName 通过serve name 获取engine
func GinEngineByServeName(name string) (*gin.Engine, bool) {
	engine, ok := ginEngine[name]
	return engine, ok
}
