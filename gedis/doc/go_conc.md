## Atomic 机制

- CPU 级别支持的原子操作
- X86平台：给内存加锁，再操作
- Arm平台：先操作，如果操作失败，再重试

```go
func add(p *int32) {
	*p++

	//*p = *p + 1

	//atomic.AddInt32(p, 1)

}

func main() {
	c := int32(0)
	for i := 0; i < 1000; i++ {
		go add(&c)
	}
	time.Sleep(5 * time.Second)
	fmt.Println(c)
}
```

## 读写锁的使用

```go
type Person struct {
	mu     sync.RWMutex
	salary int
	level  int
}

func promote(p *Person) {
	p.mu.Lock()
	p.salary++
	fmt.Println(p.salary)
	p.level++
	fmt.Println(p.level)
	p.mu.Unlock()
}

func printPerson(p *Person) {
	defer p.mu.RUnlock()
	p.mu.RLock()
	fmt.Println(p.salary)
	fmt.Println(p.level)
}

func main() {

	p := Person{level: 1, salary: 10000}

	go promote(&p)
	go promote(&p)
	go promote(&p)

	time.Sleep(time.Second)

}
```

## waitgroup的使用

```go
type Person struct {
	mu     sync.RWMutex
	salary int
	level  int
}

func promote(p *Person, wg *sync.WaitGroup) {
	p.mu.Lock()
	p.salary++
	fmt.Println(p.salary)
	p.level++
	fmt.Println(p.level)
	p.mu.Unlock()
	wg.Done()
}

func printPerson(p *Person, wg *sync.WaitGroup) {
	defer p.mu.RUnlock()
	p.mu.RLock()
	fmt.Println(p.salary)
	fmt.Println(p.level)
}

func main() {

	p := Person{level: 1, salary: 10000}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go promote(&p, &wg)
	go promote(&p, &wg)
	go promote(&p, &wg)
	wg.Wait()
	//time.Sleep(time.Second)

}
```

## 锁拷贝问题

- 锁拷贝可能导致锁的死锁问题
- 使用 vet 工具可以检测锁拷贝问题
- vet 还能检测可能的 bug 或者可疑的构造

```go
func main() {

	m := sync.Mutex{}
	n := m
	fmt.Println(n)
}
```

```sh
go vet 
```

## RACE 竞争检测

```go
var J int

func do() {
	J++
}

func main() {

	for i := 0; i < 200; i++ {
		go do()
	}
}
```

## dead lock 检测

```go
package main

import (
	sync "github.com/sasha-s/go-deadlock"
	"time"
)

var J int

var M = sync.Mutex{}

func do() {
	M.Lock()
	J++
	time.Sleep(100000000000)
	M.Unlock()
}

func main() {
	sync.Opts.DeadlockTimeout = time.Millisecond * 100
	for i := 0; i < 200; i++ {
		go do()
	}
	time.Sleep(10000000000000000) 
}
```



