package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

/*
用SingleFlight合并重复请求
Package singleflight provides a duplicate function call suppression mechanism.

具体到Go程序运行的层面来说，SingleFlight的作用是在处理多个goroutine同时调用同一个函数的时候，只让一个goroutine去实际调用这个函数，等到这个goroutine返回结果的时候，再把结果返回给其他几个同时调用了相同函数的goroutine，这样可以减少并发调用的数量。在实际应用中也是，它能够在一个服务中减少对下游的并发重复请求。还有一个比较常见的使用场景是用来防止缓存击穿。

singleflight.Group类型提供了三个方法：
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool)

func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result

func (g *Group) Forget(key string)

Do方法，接受一个字符串Key和一个待调用的函数，会返回调用函数的结果和错误。使用Do方法的时候，它会根据提供的Key判断是否去真正调用fn函数。同一个 key，在同一时间只有第一次调用Do方法时才会去执行fn函数，其他并发的请求会等待调用的执行结果。
DoChan方法：类似Do方法，只不过是一个异步调用。它会返回一个通道，等fn函数执行完，产生了结果以后，就能从这个 chan 中接收这个结果。
Forget方法：在SingleFlight中删除一个Key。这样一来，之后这个Key的Do方法调用会执行fn函数，而不是等待前一个未完成的fn 函数的结果。

《查询DNS记录》

Go语言的net标准库里使用的lookupGroup结构，就是Go扩展库提供的原语 singleflight.Group

type Resolver struct {
  ......
 // 源码地址 https://github.com/golang/go/blob/master/src/net/lookup.go#L151
 // lookupGroup merges LookupIPAddr calls together for lookups for the same
 // host. The lookupGroup key is the LookupIPAddr.host argument.
 // The return values are ([]IPAddr, error).
 lookupGroup singleflight.Group
}

它的作用是将对相同域名的DNS记录查询合并成一个查询，下面是net库提供的DNS记录查询方法LookupIp使用lookupGroup这个SingleFlight进行合并查询的相关源码，它使用的是异步查询的方法DoChan。

func LookupIP(host string) ([]IP, error) {
 addrs, err := DefaultResolver.LookupIPAddr(context.Background(), host)
  ......
}

func (r *Resolver) lookupIPAddr(ctx context.Context, network, host string) ([]IPAddr, error) {
  ......
  // 使用SingleFlight的DoChan合并多个查询请求
 ch, called := r.getLookupGroup().DoChan(lookupKey, func() (interface{}, error) {
  defer dnsWaitGroup.Done()
  return testHookLookupIP(lookupGroupCtx, resolverFunc, network, host)
 })
 if !called {
  dnsWaitGroup.Done()
 }

 select {
 case <-ctx.Done():
  ......
 case r := <-ch:
  lookupGroupCancel()
  if trace != nil && trace.DNSDone != nil {
   addrs, _ := r.Val.([]IPAddr)
   trace.DNSDone(ipAddrsEface(addrs), r.Shared, r.Err)
  }
  return lookupIPReturn(r.Val, r.Err, r.Shared)
 }
}

上面的源码做了很多删减，只留了SingleFlight合并查询的部分，如果有兴趣可以去GitHub上看一下完整的源码，访问链接https://github.com/golang/go/blob/master/src/net/lookup.go#L261 ，可直接定位到这部分的源码。

《防止缓存击穿》
在项目里使用缓存时，一个常见的用法是查询一个数据先去查询缓存，如果没有就去数据库里查到数据并缓存到Redis里。那么缓存击穿问题是指，高并发的系统中，大量的请求同时查询一个缓存Key 时，如果这个 Key 正好过期失效，就会导致大量的请求都打到数据库上，这就是缓存击穿。用 SingleFlight 来解决缓存击穿问题再合适不过，这个时候只要这些对同一个 Key 的并发请求的其中一个到数据库中查询就可以了，这些并发的请求可以共享同一个结果。

下面是一个模拟用SingleFlight并发原语合并查询Redis缓存的程序，你可以自己动手测试一下，开10个goroutine去查询一个固定的Key，观察一下返回结果就会发现最终只执行了一次Redis查询。

看一下singleflight.Group的实现原理，通过它的源码也是能学到不少用Go语言编程的技巧的。singleflight.Group由一个互斥锁sync.Mutex和一个映射表组成，每一个 singleflight.call结构体都保存了当前调用对应的信息：
type Group struct {
 mu sync.Mutex
 m  map[string]*call
}

type call struct {
 wg sync.WaitGroup

 val interface{}
 err error

 dups  int
 chans []chan<- Result
}



*/

// client 模拟 redis client
type client struct {
	// 忽略其他配置
	requestGroup singleflight.Group
}

// Get 查询
func (c *client) Get(key string) (interface{}, error) {
	fmt.Println("Querying Database")
	time.Sleep(time.Second)
	v := fmt.Sprintf("Content of key: %s", key)
	return v, nil
}

// SingleFlight 查询
func (c *client) SingleFlight(key string) (interface{}, error) {
	v, err, _ := c.requestGroup.Do(key, func() (interface{}, error) {
		return c.Get(key)
	})
	return v, err
}

// SingleFlightTime 查询超时
func (c *client) SingleFlightTime(key string, timeout time.Duration) (interface{}, error) {
	vChan := c.requestGroup.DoChan(key, func() (interface{}, error) {
		return c.Get(key)
	})
	select {
	case v := <-vChan:
		return v.Val, v.Err
	case <-time.After(timeout):
		return nil, errors.New("timeout")
	}
}

func main() {
	//client := new(client)
	client := &client{}
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			v, err := client.SingleFlight("key")
			if err != nil {
				return
			}
			fmt.Println(v)
		}()
	}
	wg.Wait()

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			v, err := client.SingleFlightTime("key for timeout", time.Second*2)
			if err != nil {
				fmt.Printf("%+v\n", err)
				return
			}
			fmt.Println(v)
		}()
	}
	wg.Wait()
}
