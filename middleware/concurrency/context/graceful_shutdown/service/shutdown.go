////go:build graceful

package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

type Option func(*App)

// ShutdownCallback 采用 context.Context 来控制超时，而不是用 time.After 是因为
// - 超时本质上是使用这个回调的人控制的
// - 我们还希望用户知道，他的回调必须要在一定时间内处理完毕，而且他必须显式处理超时错误
type ShutdownCallback func(ctx context.Context)

// WithShutdownCallbacks 通过此处设置shutdown回调
func WithShutdownCallbacks(callbacks ...ShutdownCallback) Option {
	return func(app *App) {
		app.callbacks = append(app.callbacks, callbacks...)
	}
}

// App 定义应用信息
type App struct {
	// 一个应用启动多个服务
	servers []*Server

	// 优雅退出整个超时时间，默认30秒
	shutdownTimeout time.Duration

	// 优雅退出时候等待处理已有请求时间，默认10秒钟
	waitTime time.Duration
	// 自定义回调超时时间，默认三秒钟
	callbackTimeout time.Duration

	// 关闭app时回调
	callbacks []ShutdownCallback
}

func NewApp(servers []*Server, opts ...Option) *App {
	res := &App{
		waitTime:        10 * time.Second,
		callbackTimeout: 3 * time.Second,
		shutdownTimeout: 30 * time.Second,
		servers:         servers,
	}
	for _, opt := range opts {
		opt(res)
	}

	return res
}

// StartAndServe 启动整个app
func (app *App) StartAndServe() {
	for _, server := range app.servers {
		srv := server
		go func() {
			if err := srv.Start(); err != nil {
				if err == http.ErrServerClosed {
					log.Printf("服务器[%s]已关闭", srv.name)
				} else {
					log.Printf("服务器[%s]异常退出", srv.name)
				}
			}
		}()
	}

	// 从这里开始开始启动监听系统信号
	// ch := make(...) 首先创建一个接收系统信号的 channel ch
	// 定义要监听的目标信号 signals []os.Signal
	// 调用 signal
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, signals...)
	<-ch

	go func() {
		select {
		case <-ch:
			log.Printf("强制退出")
			os.Exit(1)
		case <-time.After(app.shutdownTimeout):
			log.Printf("超时强制退出")
			os.Exit(1)
		}
	}()

	app.shutdown()
}

// shutdown 关闭应用
func (app *App) shutdown() {
	log.Println("开始关闭应用，停止接受新的请求")
	for _, server := range app.servers {
		// 思考：这里为什么可以不用并发控制，即不用锁，也不用原子操作
		server.rejectReq()
	}

	log.Println("等待正在执行请求完结")
	// 这里可以改造为实时统计正在处理的请求数量，为0 则下一步
	time.Sleep(app.waitTime)

	log.Println("开始关闭服务器")
	var wg sync.WaitGroup
	wg.Add(len(app.servers))

	for _, srv := range app.servers {
		srvCp := srv
		go func() {
			if err := srvCp.stop(); err != nil {
				log.Printf("关闭服务[%s]失败", srvCp.name)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	log.Println("开始执行自定义回调")
	// 执行回调
	wg.Add(len(app.callbacks))

	for _, callback := range app.callbacks {
		cb := callback
		go func() {
			timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), app.callbackTimeout)
			cb(timeoutCtx)
			cancelFunc()
			wg.Done()
		}()
	}

	wg.Wait()

	// 释放资源
	log.Println("开始释放资源")
	app.close()
}

func (app *App) close() {
	// 在这里释放掉一些可能的资源
	time.Sleep(time.Second)
	log.Println("应用已关闭")
}

// Server 定义
type Server struct {
	srv  *http.Server
	name string
	mux  *serveMux
}

// serverMux 既可以看做是装饰器模式，也可以看做委托模式
type serveMux struct {
	reject atomic.Bool
	*http.ServeMux
}

func (s *serveMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.reject.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("服务已关闭"))
		return
	}
	s.ServeMux.ServeHTTP(w, r)
}

func NewServer(name, addr string) *Server {
	mux := &serveMux{ServeMux: http.NewServeMux()}
	return &Server{
		name: name,
		mux:  mux,
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) rejectReq() {
	s.mux.reject.Store(true)
}

func (s *Server) stop() error {
	log.Printf("服务器[%s]关闭中...", s.name)
	return s.srv.Shutdown(context.Background())
}

func (s *Server) stopWithCtx(ctx context.Context) error {
	log.Printf("服务器[%s]关闭中...", s.name)
	return s.srv.Shutdown(ctx)
}
