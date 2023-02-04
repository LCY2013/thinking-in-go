## Channel 声明方法

- chInt := make(chan int)       // unbuffered channel  非缓冲通道
- chBool := make(chan bool, 0)  // unbuffered channel  非缓冲通道
- chStr := make(chan string, 2) // bufferd channel     缓冲通道

## Channel 基本用法

- ch <- x // channel 接收数据 x
- x <- ch // channel 发送数据并赋值给 x
- <- ch // channel 发送数据，忽略接受者

## 错误用法

```go
func main() {
    ch := make(chan string)

    ch <- "ping"

    fmt.Println(<-ch)
}
```

```go
func main() {
    ch := make(chan string)

    go func() {
        ch <- "ping"
    }()

    fmt.Println(<-ch)
}
```

## 内存与通信

- “不要通过共享内存的方式进行通信”
- “而是应该通过通信的方式共享内存”


```go
func watch(p *int) {
    for true {
        if *p == 1 {
            fmt.Println("hello")
            break
        }
    }
}

func main() {
    i := 0
    go watch(&i)

    time.Sleep(time.Second)

    i = 1

    time.Sleep(time.Second)
}
```

```go
func watch(c chan int) {

    if <-c == 1 {
        fmt.Println("hello")
    }
}

func main() {
    c := make(chan int)
    go watch(c)

    time.Sleep(time.Second)

    c <- 1

    time.Sleep(time.Second)
}
```

