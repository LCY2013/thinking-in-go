package main

import (
	"fmt"
	"fufeng.org/concurrent_mode/pool"
	"io"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// 展示如何利用 pool 包来共享数据库连接池
const (
	maxGoroutines   = 25 // 要使用goroutine的数量
	pooledResources = 2  // 池中的资源的数量
)

// 模拟共享的资源
type dbConnection struct {
	ID int32
}

// Close 实现了 io.Closer 接口，以便 dbConnection 可以被池管理
// Close 用来管理任意资源的释放工作
func (dbConn *dbConnection) Close() error {
	log.Println("Close: Connection", dbConn.ID)
	return nil
}

// idCounter 用来给每一个 dbConnection 分配一个唯一的 id
var idCounter int32

// createConnection 是一个工厂函数，当需要一个新的连接时，资源池会调用这个函数
func createConnection() (io.Closer, error) {
	// 通过 atomic 同步新增 idCounter 的值
	id := atomic.AddInt32(&idCounter, 1)
	log.Println("Create: New Connection", id)

	// 创建一个新的连接
	return &dbConnection{id}, nil
}

// main 程序的主入口
func main() {
	// 创建一个等待组
	var wg sync.WaitGroup
	// 等待 goroutine 所有项目完成
	wg.Add(maxGoroutines)

	// 创建用来管理连接的池
	p, err := pool.New(createConnection, pooledResources)
	if err != nil {
		log.Println(err)
	}

	// 使用池中的连接来完成查询
	for query := 0; query < maxGoroutines; query++ {
		// 每个 goroutine 需要复制一份要查询值的副本
		// 不然所有的查询会共享一个查询变量
		go func(q int) {
			performQueries(q, p)
			wg.Done()
		}(query)
	}

	// 等待所有 goroutine 的结束
	wg.Wait()

	// 关闭池
	fmt.Println("Shutdown Program.")
	p.Close()
}

// performQueries 用来测试连接的资源池
func performQueries(query int, p *pool.Pool) {
	// 从池中获取一个连接信息
	conn, err := p.Acquires()
	if err != nil {
		log.Println(err)
		return
	}

	// 将该连接放回到池中
	defer p.Release(conn)

	// 用等待模拟数据库查询时间
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	log.Printf("QID[%d] CID[%d] \n", query, conn.(*dbConnection).ID)
}
