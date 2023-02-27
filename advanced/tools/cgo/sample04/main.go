package main

//void SayHello(const char* s);
import "C"

/*
使用自己的 C 函数

先自定义一个叫 SayHello 的 C 函数来实现打印，然后从 Go 语言环境中调用这个 SayHello 函数。

可以将 SayHello 函数放到当前目录下的一个 C 语言源文件中（后缀名必须是 .c）。
因为是编写在独立的 C 文件中，为了允许外部引用，所以需要去掉函数的 static 修饰符。

注意，如果之前运行的命令是 go run hello.go 或 go build hello.go 的话，
此处须使用 go run "your/package" 或 go build "your/package" 才可以。
若本就在包路径下的话，也可以直接运行 go run . 或 go build。

既然 SayHello 函数已经放到独立的 C 文件中了，我们自然可以将对应的 C 文件编译打包为静态库或动态库文件供使用。
如果是以静态库或动态库方式引用 SayHello 函数的话，需要将对应的 C 源文件移出当前目录（CGO 构建程序会自动构建当前目录下的 C 源文件，从而导致 C 函数名冲突）。
*/
func main() {
	C.SayHello(C.CString("hello, world!\n"))
}
