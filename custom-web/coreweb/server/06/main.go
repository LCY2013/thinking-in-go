package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lcy2013/custom-web/coreweb/server/06/framework"
	"github.com/lcy2013/custom-web/coreweb/server/06/framework/middleware"
)

func main() {
	core := framework.NewCore()
	// core.Use(
	// 	middleware.Test1(),
	// 	middleware.Test2())
	core.Use(middleware.Recovery())
	core.Use(middleware.Cost())
	// core.Use(middleware.Timeout(1 * time.Second))

	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}

	go func() {
		server.ListenAndServe()
		//err := server.ListenAndServe()
		/*if err != nil {
			log.Fatal(err)
		}*/
	}()

	// 当前的 Goroutine 等待信号量
	quit := make(chan os.Signal, 2)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGSTOP,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
		syscall.SIGABRT, syscall.SIGSYS, syscall.SIGTERM)
	// 这里会阻塞当前 Goroutine 等待信号
	<-quit

	go func() {
		select {
		case <-quit:
			log.Println("Shutting down now")
			os.Exit(1)
		}
	}()

	// 调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 调用Server.Shutdown graceful结束
	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
