```go
func do3() {
	fmt.Println("dododo")
}
func do2() {
	do3()
}
func do1() {
	do2()
}
```

```go
func s() {

    i := 0
    for true {
        i++
    }
}
```

## 无限开启协程

```go
func main() {
		go func(i int) {
			fmt.Println(i)
			time.Sleep(time.Second)
		}(i)
}

```

## 利用 channel 的缓存区

```go
func main() {
	ch := make(chan struct{}, 30000)
	for i := 0; i < math.MaxInt32; i++ {
		ch <- struct{}{}
		go func(i int) {
			log.Println(i)
			time.Sleep(time.Second)
			<-ch
		}(i)
	}
	time.Sleep(time.Hour)
}

```

## 协程池

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go get github.com/Jeffail/tunny
```

```go
func main() {
    pool := tunny.NewFunc(3000, func(i interface{}) interface{} {
        log.Println(i)
        time.Sleep(time.Second)
        return nil
    })
    defer pool.Close()

    for i := 0; i < 1000000; i++ {
        go pool.Process(i)
    }
    time.Sleep(time.Second * 4)
}
```

