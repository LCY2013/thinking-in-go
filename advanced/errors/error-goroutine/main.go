package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

/*
在程序启动的时候，如果有强依赖的服务出现故障时 panic 退出
在程序启动的时候，如果发现有配置明显不符合要求， 可以 panic 退出（防御编程）
其他情况下只要不是不可恢复的程序错误，都不应该直接 panic 应该返回 error
在程序入口处，例如 gin 中间件需要使用 recover 预防 panic 程序退出

在程序中我们应该避免使用野生的 goroutine
1、如果是在请求中需要执行异步任务，应该使用异步 worker ，消息通知的方式进行处理，避免请求量大时大量 goroutine 创建
2、如果需要使用 goroutine 时，应该使用同一的 Go 函数进行创建，这个函数中会进行 recover ，避免因为野生 goroutine panic 导致主进程退出
*/
func main() {
	//url := "http://qw-scrm.privatecloud-xinjushang-dev:4000/internal/qw-scrm-svc"
	//url = strings.ReplaceAll(url, "/internal", "")
	//fmt.Println(url)
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-08-31 10:08:55", time.Local)
	fmt.Println(t.Unix())
	fmt.Println(t.Unix() - time.Now().Unix())
}

// serveDebugV1 开启 go 的debug获取火焰图等信息
func serveDebugV1() {
	_ = http.ListenAndServe("127.0.0.1:8001", http.DefaultServeMux)
}

// serveDebugV2 开启 go 的debug获取火焰图等信息
func serveDebugV2() {
	if err := http.ListenAndServe("127.0.0.1:8001", http.DefaultServeMux); err != nil {
		// 该函数最后会调用 os.Exit 导致 go的defer指令不会被执行
		log.Fatal(err)
	}
}

// GO 全局声明的异步 - 执行异步
func GO(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v\n", err)
			}
		}()

		f()
	}()
}

/*
error：
1、我们在应用程序中使用 github.com/pkg/errors 处理应用错误，注意在公共库当中，我们一般不使用这个

2、error 应该是函数的最后一个返回值，当 error 不为 nil 时，函数的其他返回值是不可用的状态，不应该对其他返回值做任何期待
	func f() (io.Reader, *S1, error) 在这里，我们不知道 io.Reader 中是否有数据，可能有，也有可能有一部分

3、错误处理的时候应该先判断错误， if err != nil 出现错误及时返回，使代码是一条流畅的直线，避免过多的嵌套.

4、在应用程序中出现错误时，使用 errors.New 或者 errors.Errorf 返回错误

5、如果是调用应用程序的其他函数出现错误，请直接返回，如果需要携带信息，请使用 errors.WithMessage

6、如果是调用其他库（标准库、企业公共库、开源第三方库等）获取到错误时，请使用 errors.Wrap 添加堆栈信息
	切记，不要每个地方都是用 errors.Wrap 只需要在错误第一次出现时进行 errors.Wrap 即可
	根据场景进行判断是否需要将其他库的原始错误吞掉，例如可以把 repository 层的数据库相关错误吞掉，返回业务错误码，避免后续我们分割微服务或者更换 ORM 库时需要去修改上层代码
	注意我们在基础库，被大量引入的第三方库编写时一般不使用 errors.Wrap 避免堆栈信息重复

7、禁止每个出错的地方都打日志，只需要在进程的最开始的地方使用 %+v 进行统一打印，例如 http/rpc 服务的中间件

8、错误判断使用 errors.Is 进行比较

9、错误类型判断，使用 errors.As 进行赋值

10、如何判定错误的信息是否足够，想一想当你的代码出现问题需要排查的时候你的错误信息是否可以帮助你快速的定位问题，例如我们在请求中一般会输出参数信息，用于辅助判断错误

11、对于业务错误，推荐在一个统一的地方创建一个错误字典，错误字典里面应该包含错误的 code，并且在日志中作为独立字段打印，方便做业务告警的判断，错误必须有清晰的错误文档

12、不需要返回，被忽略的错误必须输出日志信息

13、同一个地方不停的报错，最好不要不停输出错误日志，这样可能会导致被大量的错误日志信息淹没，无法排查问题，比较好的做法是打印一次错误详情，然后打印出错误出现的次数

14、对同一个类型的错误，采用相同的模式，例如参数错误，不要有的返回 404 有的返回 200

15、处理错误的时候，需要处理已分配的资源，使用 defer 进行清理，例如文件句柄

*/

// case4
func case4() error {
	file, err := os.OpenFile("", os.O_RDWR, 666)
	if err != nil {
		return errors.WithStack(err)
	}

	if file == nil {
		return errors.Errorf("打开文件错误: %s", err)
	}

	if file.Name() != "fufeng.txt" {
		return errors.WithMessage(err, "文件名称不对")
	}

	return nil
}

// case6
func case6() error {
	var request struct{}
	var data []byte
	err := json.Unmarshal(data, &request)
	if err != nil {
		return errors.Wrap(err, "其他附加信息")
	}

	// 其他
	return nil
}

// case8
func case8() error {
	err := case6()
	if errors.Is(err, io.EOF) {
		return io.EOF
	}
	return nil
}

// case9
func case9() error {
	err := case8()
	var errEof error = io.EOF
	if errors.As(err, &errEof) {
		return errEof
	}

	return nil
}

/*
panic or error?

1、在 Go 中 panic 会导致程序直接退出，是一个致命的错误，如果使用 panic recover 进行处理的话，会存在很多问题

性能问题，频繁 panic recover 性能不好

容易导致程序异常退出，只要有一个地方没有处理到就会导致程序进程整个退出

不可控，一旦 panic 就将处理逻辑移交给了外部，我们并不能预设外部包一定会进行处理

2、什么时候使用 panic 呢？

对于真正意外的情况，那些表示不可恢复的程序错误，例如索引越界、不可恢复的环境问题、栈溢出，我们才使用 panic

3、使用 error 处理有哪些好处？

简单。

考虑失败，而不是成功(Plan for failure, not success)。

没有隐藏的控制流。

完全交给你来控制 error。

Error are values。
*/

/*
为什么标准库中 errors.New 会返回一个指针？
翻看标准库的源代码可以发现， errors 库中的 errorString 结构体实现了 error 接口，为什么在 New 一个 error 的时候会返回一个结构体的指针呢？
//
// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
//

我们先来看一个例子，我们同样创建了 errorString 的结构体，我们自定义的和标准库中的唯一不同就是，自建的这个返回的是值，而不是指针。
//
type errorString struct {
	text string
}

func (e errorString) Error() string {
	return e.text
}

// New 创建一个自定义错误
func New(s string) error {
	return errorString{text: s}
}

var errorString1 = New("test a")
var err1 = errors.New("test b")

func main() {
	if errorString1 == New("test a") {
		fmt.Println("err string a") // 会输出
	}

	if err1 == errors.New("test b") {
		fmt.Println("err b") // 不会输出
	}
}
//

在 main 函数的对比中我们就可以发现，我们自定义的 errorString 在对比的时候只要对应的字符串相同就会返回 true，但是标准库的包不会。

这是因为，在对比两个 struct 是否相同的时候，会去对比，这两个 struct 里面的各个字段是否是相同的，如果相同就返回 true，但是对比指针的时候会去判断两个指针的地址是否一致。
*/

/*
error type: 错误定义与判断

Sentinel Error

哨兵错误，就是定义一些包级别的错误变量，然后在调用的时候外部包可以直接对比变量进行判定，在标准库当中大量的使用了这种方式
例如下方 io 库中定义的错误
//
// EOF is the error returned by Read when no more input is available.
// Functions should return EOF only to signal a graceful end of input.
// If the EOF occurs unexpectedly in a structured data stream,
// the appropriate error is either ErrUnexpectedEOF or some other error
// giving more detail.
var EOF = errors.New("EOF")

// ErrUnexpectedEOF means that EOF was encountered in the
// middle of reading a fixed-size block or data structure.
var ErrUnexpectedEOF = errors.New("unexpected EOF")

// ErrNoProgress is returned by some clients of an io.Reader when
// many calls to Read have failed to return any data or error,
// usually the sign of a broken io.Reader implementation.
var ErrNoProgress = errors.New("multiple Read calls return no data or error")
//

我们在外部判定的时候一般使用等值判定或者使用 errors.Is 进行判断
//
if err == io.EOF {
	//...
}

if errors.Is(err, io.EOF){
	//...
}
//
*/

/*
error types

这个就类似我们前面定义的 errorString 一样实现了 error 的接口，然后在外部是否类型断言来判断是否是这种错误类型

//
type MyStruct struct {
	s string
    name string
    path string
}



// 使用的时候
func f() {
    switch err.(type) {
        case *MyStruct:
        // ...
        case others:
        // ...
    }
}
//

这种方式相对于哨兵来说，可以包含更加丰富的信息，但是同样也将错误的类型暴露给了外部，例如标准库中的 os.PathError


*/
