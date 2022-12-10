package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/LCY2013/thinking-in-go/middleware/concurrency/context/graceful_shutdown/service"
)

// main 注意要从命令行启动，否则不同的 IDE 可能会吞掉关闭信号
// go build -tags=answer
func main() {
	businessServer := service.NewServer("business", "localhost:8080")
	businessServer.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello"))
	}))

	adminServer := service.NewServer("admin", "localhost:8081")
	adminServer.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("admin"))
	}))

	app := service.NewApp([]*service.Server{businessServer, adminServer}, service.WithShutdownCallbacks(StoreCacheToDBCallback))
	app.StartAndServe()
}

func StoreCacheToDBCallback(ctx context.Context) {
	done := make(chan struct{}, 1)
	go func() {
		// 业务逻辑，比如说这里我们模拟的是将本地缓存刷新到数据库里面
		// 这里简单的睡一段时间来模拟
		log.Printf("刷新缓存中……")
		time.Sleep(time.Millisecond * 500)
		done <- struct{}{}
	}()

	select {
	case <-done:
		log.Printf("缓存被刷新到了 DB")
	case <-ctx.Done():
		log.Printf("缓存刷新超时")
	}
}
