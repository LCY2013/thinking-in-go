## cgo

```go
package main

/*
int sum(int a, int b) {
  return a+b;
}
*/

import "C"

func main() {
    println(C.sum(1, 1))
}
```

```go
package main

import "fmt"

func main() {

    defer fmt.Println("defer1")
    defer fmt.Println("defer2")

    fmt.Println("do something")
}
```

```go
package main

import (
    "fmt"
    "time"
)

func dosomething() {
    panic("panic1")

    fmt.Println("do something2")

}

func main() {

    dosomething()
    fmt.Println("do something1")

    time.Sleep(1 * time.Second)
}
```

## panic基本使用

```go
package main

import (
    "fmt"
    "time"
)

func dosomething() {
    panic("panic1")

    fmt.Println("do something2")

}

func main() {

    go dosomething()
    fmt.Println("do something1")

    time.Sleep(1 * time.Second)
}
```

## panic + defer

```go
package main

import (
    "fmt"
    "time"
)

func dosomething() {

    defer fmt.Println("defer2")

    panic("panic1")

    fmt.Println("do something2")

}

func main() {
    defer fmt.Println("defer1")

    dosomething()
    fmt.Println("do something1")

    time.Sleep(1 * time.Second)
}
```

- panic在退出协程之前会执行所有已注册的defer

## panic + defer + recover

```go
package main

import (
	"fmt"
	"time"
)

func dosomething() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()
	panic("panic1")

	fmt.Println("do something2")

}

func main() {
	defer fmt.Println("defer1")

	dosomething()
	fmt.Println("do something1")

	time.Sleep(1 * time.Second)
}
```

## 对象到反射对象

```go
func main() {
	s := "moody"

	stype := reflect.TypeOf(s)
	fmt.Println("TypeOf s:", stype)

	svalue := reflect.ValueOf(s)
	fmt.Println("ValueOf s:", svalue)
}
```


## 反射对象到对象

```go
func main() {
    s := "moody"

    stype := reflect.TypeOf(s)
    fmt.Println("TypeOf s:", stype)

    svalue := reflect.ValueOf(s)
    fmt.Println("ValueOf s:", svalue)

    s2 := svalue.Interface().(string)

    fmt.Println("s2:", s2)
}
```

## 发射调用方法

```go
package main

import (
	"fmt"
	"reflect"
)

func MyAdd(a, b int) int { return a + b }

func CallAdd(f func(a int, b int) int) {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return
	}
	argv := make([]reflect.Value, 2)
	argv[0] = reflect.ValueOf(1)
	argv[1] = reflect.ValueOf(1)

	result := v.Call(argv)

	fmt.Println(result[0].Int())
}

func main() {

	CallAdd(MyAdd)
}
```
 
