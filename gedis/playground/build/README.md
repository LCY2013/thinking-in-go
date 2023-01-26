## go的特点

|            | 一次编码 | 一次编译 | 不需要运行环境 | 没有虚拟化损失 | 不需要自行处理 GC | 面向对象 | 非常易用的并发能力 |
| ---        | ---      | ---     | ---            | ---            | ---               | ---      | ---               |
| C          | X        | √        | X             | √              | X                 | X        | X                 |
| C++        | X        | √        | X             | √              | X                 | √        | X                 |
| Java       | √        | X        | √             | X              | √                 | √        | X                 |
| JavaScript | √        | O        | √             | X              | √                 | √        | X                 |
| Go         | √        | X        | √             | √              | √                 | √        | √                 |

## 查看从代码到SSA中间码的整个过程

```
$env:GOSSAFUNC="main"
```

```
export GOSSAFUNC=main
```

```
go build
```

## 查看 Plan9 汇编代码

```
go build -gcflags -S main.go
```

## 使用 Modules

```
go get github.com/xxx/xxx 
```

```
go get github.com/xxx/xxx@0.0.1
```

## 使用goproxy.cn作为代理

```
go env -w GOPROXY=https://goproxy.cn,direct
```

## go.mod 文件追加

```
replace github.com/xxx/xxx => xxx/xxx
```

## go vender 缓存到本地

```
go mod vendor
```

```
go build -mod vendor
```

## 创建 Go Module

```
go mod init github.com/xxx/xxx
```
