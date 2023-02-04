```GO
package main

func sum(a, b int) int {
    sum := 0
    sum = a + b
    return sum
}

func main() {
    a := 3
    b := 5
    print(
        sum(a, b)
        )
}
```

## 指针逃逸

- 函数返回了对象的指针

```go
package main

import "fmt"

type Demo struct {
    name string
}

func createDemo(name string) *Demo {
    d := new(Demo) // 局部变量 d 逃逸到堆
    d.name = name
    return d
}

func main() {
    demo := createDemo("demo")
    fmt.Println(demo)
}

```

## GC分析工具

```
go tool pprof
```

```
go tool trace
```

```
go build -gcflags=”-m”
```

```
GODEBUG=”gctrace=1”
```


