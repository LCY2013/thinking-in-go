## 基本类型大小

```go
package main

import (
    "fmt"
    "unsafe"
)

func main(){
    i := 1234
    j:= int32(1)
    f:=float32(3.141)
    bytes := [5]byte{'h','e','l','l','o'}
    primes := [4]int{2,3,5,7}
    p:= &primes

    r:= rune(666)


    fmt.Println(unsafe.Sizeof(i))
    fmt.Println(unsafe.Sizeof(j))
    fmt.Println(unsafe.Sizeof(f))
    fmt.Println(unsafe.Sizeof(bytes))
    fmt.Println(unsafe.Sizeof(primes))
    fmt.Println(unsafe.Sizeof(p))
    fmt.Println(unsafe.Sizeof(r))
    fmt.Println(unsafe.Sizeof("你好"))

}
```

## 字符串转StringHeader

```go
func main() {

    str := "你a" 

    for _, s := range str {
        fmt.Printf("unicode: %c %d %T\n", s, s, s)
    }

    for i := 0; i < len(str); i++ {
        fmt.Printf("ascii: %c %d %T\n", str[i], str[i], str[i])
    }

}
```

## 字符串的两种遍历

```go
    str := "快乐 everyday"

    str1 := str[3:5]

        for _, s := range str1{
        fmt.Printf("unicode: %c %d %T\n", s,s)
    }

    for i:=0;i<len(str1) ;i++  {
        fmt.Printf("ascii: %c %d %T\n", str[i], str[i])
```

## 切片的三种创建方式

```go
arr[0:3] or slice[0:3]
slice := []int{1, 2, 3}
slice := make([]int, 10)
```

## Go隐式接口特点

```go
package main

import "fmt"

type taxi struct {
}

func (t taxi) Drive() {
    fmt.Println("Drive taxi")
}

func (t taxi) MakeMoney() {
    fmt.Println("Make Money")
}

type Car interface {
    Drive()
}

type MoneyMaker interface {
    MakeMoney()
}

func main() {

}

```

## 结构体和指针实现接口

```go
type Car interface {
    Drive()
}

type truck struct {
}

func (t *truck) Drive() {

}

func main() {
    var a Car = truck{}
    fmt.Println(reflect.TypeOf(a))
}

```

## 变量对齐

```go
type Args struct {
    num1 int32
    num2 int32
}

type Flag struct {
    num1 int16
    num2 int32
}

func main() {
    fmt.Println(unsafe.Sizeof(Args{}))
    fmt.Println(unsafe.Sizeof(Flag{}))
}
```

```go

unsafe.Alignof(Args{}) // 8
unsafe.Alignof(Flag{}) // 4
```

```go
func main() {
    fmt.Printf("bool size:%d align: %d\n", unsafe.Sizeof(bool(true)), unsafe.Alignof(bool(true)))
    fmt.Printf("byte size:%d align: %d\n", unsafe.Sizeof(byte(0)), unsafe.Alignof(byte(0)))
    fmt.Printf("int8 size:%d align: %d\n", unsafe.Sizeof(int8(0)), unsafe.Alignof(int8(0)))
    fmt.Printf("int16 size:%d align: %d\n", unsafe.Sizeof(int16(0)), unsafe.Alignof(int16(0)))
    fmt.Printf("int32 size:%d align: %d\n", unsafe.Sizeof(int32(0)), unsafe.Alignof(int32(0)))
    fmt.Printf("int64 size:%d align: %d\n", unsafe.Sizeof(int64(0)), unsafe.Alignof(int64(0)))
}   
```

- Go中每一个变量都有自己的对齐系数
- struct的对齐系数取决于当中的最大成员

## 字长对齐

- 对于特定系统，也有系统对齐系数，一般为系统字长
- 变量要尽量放置在一个系统字长里

```go
    var a bool
    var b int16
    var c int

    fmt.Printf("%p\n", &a)
    fmt.Printf("%p\n", &b)
    fmt.Printf("%p\n", &c)
```

