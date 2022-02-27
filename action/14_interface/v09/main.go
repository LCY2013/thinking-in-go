package main

/*
Go 接口的应用模式或惯例

先要了解一个前置原则，那就是在实际真正需要的时候才对程序进行抽象。
再通俗一些来讲，就是不要为了抽象而抽象。
接口本质上是一种抽象，它的功能是解耦，所以这条原则也在告诉我们：不要为了使用接口而使用接口。

举一个简单的例子，如果我们要给一个计算器添加一个整数加法的功能特性，本来一个函数就可以实现：
func Add(a int64, b int64) int64 {
  return a+b
}

但如果你非要引入一个接口，结果代码可能就变成了这样：
type Adder interface {
    Add(int64, int64) int64
}
func Add(adder Adder, a int64, b int64) int64 {
  return adder.Add(a, b)
}

一切皆组合

编程语言的语法元素间和语言特性也存在着正交的情况，并且通过将这些正交的特性组合起来，我们可以实现更为高级的特性。
在语言设计层面，Go 语言就为广大 Gopher 提供了诸多正交的语法元素供后续组合使用，包括：
1、Go 语言无类型体系（Type Hierarchy），没有父子类的概念，类型定义是正交独立的；
2、方法和类型是正交的，每种类型都可以拥有自己的方法集合，方法本质上只是一个将 receiver 参数作为第一个参数的函数而已；
3、接口与它的实现者之间无“显式关联”，也就说接口与 Go 语言其他部分也是正交的。

#### go 组合的方式
1、垂直组合
垂直组合更多用在将多个类型通过“类型嵌入（Type Embedding）”的方式实现新类型的定义。

传统面向对象编程语言（比如：C++）大多是通过继承的方式建构出自己的类型体系的，但 Go 语言并没有类型体系的概念。
Go 语言通过类型的组合而不是继承让单一类型承载更多的功能，由于这种方式与硬件配置升级的垂直扩展很类似，所以这里叫它垂直组合。

又因为不是继承，那么通过垂直组合定义的新类型与被嵌入的类型之间就没有所谓“父子关系”的概念了，也没有向上、向下转型（Type Casting），被嵌入的类型也不知道将其嵌入的外部类型的存在。
调用方法时，方法的匹配取决于方法名字，而不是类型。

这样的垂直组合更多应用在新类型的定义方面。通过这种垂直组合，可以达到方法实现的复用、接口定义重用等目的。

在实现层面，Go 语言通过类型嵌入（Type Embedding）实现垂直组合，组合方式主要有以下这么几种。

第一种：通过嵌入接口构建接口
通过在接口定义中嵌入其他接口类型，实现接口行为聚合，组成大接口。这种方式在标准库中非常常见，也是 Go 接口类型定义的惯例。
比如这个 ReadWriter 接口类型就采用了这种类型嵌入方式：
// $GOROOT/src/io/io.go
type ReadWriter interface {
    Reader
    Writer
}

第二种：通过嵌入接口构建结构体类型
type MyReader struct {
  io.Reader // underlying reader
  N int64   // max bytes remaining
}
在结构体中嵌入接口，可以用于快速构建满足某一个接口的结构体类型，来满足某单元测试的需要，之后我们只需要实现少数需要的接口方法就可以了。
尤其是将这样的结构体类型变量传递赋值给大接口的时候，就更能体现嵌入接口类型的优势了。

第三种：通过嵌入结构体类型构建新结构体类型
在结构体中嵌入接口类型名和在结构体中嵌入其他结构体，都是“委派模式（delegate）”的一种应用。
对新结构体类型的方法调用，可能会被“委派”给该结构体内部嵌入的结构体的实例，通过这种方式构建的新结构体类型就“继承”了被嵌入的结构体的方法的实现。

括嵌入接口类型在内的各种垂直组合更多用于类型定义层面，本质上它是一种类型组合，也是一种类型之间的耦合方式。

2、水平组合
假设现在有一个任务，要编写一个函数，实现将一段数据写入磁盘的功能。通常可以很容易地写出下面的函数：

func Save(f *os.File, data []byte) error

这个函数使用一个 *os.File 来表示数据写入的目的地，这个函数实现后可以工作得很好。

但这里依旧存在一些问题来看一下：

首先，这个函数很难测试。os.File 是一个封装了磁盘文件描述符（又称句柄）的结构体，只有通过打开或创建真实磁盘文件才能获得这个结构体的实例，这就意味着，如果我们要对 Save 这个函数进行单元测试，就必须使用真实的磁盘文件。测试过程中，通过 Save 函数写入文件后，我们还需要再次操作文件、读取刚刚写入的内容来判断写入内容是否正确，并且每次测试结束前都要对创建的临时文件进行清理，避免给后续的测试带去影响。

其次，Save 函数违背了接口分离原则。
根据业界广泛推崇的 Robert Martin（Bob 大叔）的接口分离原则（ISP 原则，Interface Segregation Principle），也就是客户端不应该被迫依赖他们不使用的方法，我们会发现 os.File 不仅包含 Save 函数需要的与写数据相关的 Write 方法，还包含了其他与保存数据到文件操作不相关的方法。比如，你也可以看下 *os.File 包含的这些方法：
func (f *File) Chdir() error
func (f *File) Chmod(mode FileMode) error
func (f *File) Chown(uid, gid int) error
... ...

最后，Save 函数对 os.File 的强依赖让它失去了扩展性。像 Save 这样的功能函数，它日后很大可能会增加向网络存储写入数据的功能需求。但如果到那时我们再来改变 Save 函数的函数签名（参数列表 + 返回值）的话，将影响到 Save 函数的所有调用者。

综合考虑这几种原因，发现 Save 函数所在的“器官”与 os.File 所在的“器官”之间采用了一种硬连接的方式，而以 os.File 这样的结构体作为“关节”让它连接的两个“器官”丧失了相互运动的自由度，让它与它连接的两个“器官”构成的联结体变得“僵直”。

那么，我们应该如何更换“关节”来改善 Save 的设计呢？我们来试试接口。新版的 Save 函数原型如下：
func Save(w io.Writer, data []byte) error

我们用 io.Writer 接口类型替换掉了 *os.File。这样一来，新版 Save 的设计就符合了接口分离原则，因为 io.Writer 仅包含一个 Write 方法，而且这个方法恰恰是 Save 唯一需要的方法。

这里我们以 io.Writer 接口类型表示数据写入的目的地，既可以支持向磁盘写入，也可以支持向网络存储写入，并支持任何实现了 Write 方法的写入行为，这让 Save 函数的扩展性得到了质的提升。

一直强调的，接口本质是契约，具有天然的降低耦合的作用。基于这点，我们对 Save 函数的测试也将变得十分容易，比如下面示例代码：
main_test.go

接口应用的几种模式
通过接口进行水平组合的基本模式就是：使用接受接口类型参数的函数或方法。

1、基本模式
接受接口类型参数的函数或方法是水平组合的基本语法，形式是这样的：
func YourFuncName(param YourInterfaceType)

函数 / 方法参数中的接口类型作为“关节（连接点）”，支持将位于多个包中的多个类型与 YourFuncName 函数连接到一起，共同实现某一新特性。

同时，接口类型和它的实现者之间隐式的关系却在不经意间满足了：依赖抽象（DIP）、里氏替换原则（LSP）、接口隔离（ISP）等代码设计原则，这在其他语言中是需要很“刻意”地设计谋划的，但对 Go 接口来看，这一切却是自然而然的。

2、创建模式
Go 社区流传一个经验法则：“接受接口，返回结构体（Accept interfaces, return structs）”，这其实就是一种把接口作为“关节”的应用模式。我这里把它叫做创建模式，是因为这个经验法则多用于创建某一结构体类型的实例。

下面是 Go 标准库中，运用创建模式创建结构体实例的代码摘录：
// $GOROOT/src/sync/cond.go
type Cond struct {
    ... ...
    L Locker
}
func NewCond(l Locker) *Cond {
    return &Cond{L: l}
}
// $GOROOT/src/log/log.go
type Logger struct {
    mu     sync.Mutex
    prefix string
    flag   int
    out    io.Writer
    buf    []byte
}
func New(out io.Writer, prefix string, flag int) *Logger {
    return &Logger{out: out, prefix: prefix, flag: flag}
}
// $GOROOT/src/log/log.go
type Writer struct {
    err error
    buf []byte
    n   int
    wr  io.Writer
}
func NewWriterSize(w io.Writer, size int) *Writer {
    // Is it already a Writer?
    b, ok := w.(*Writer)
    if ok && len(b.buf) >= size {
        return b
    }
    if size <= 0 {
        size = defaultBufSize
    }
    return &Writer{
        buf: make([]byte, size),
        wr:  w,
    }
}

创建模式在 sync、log、bufio 包中都有应用。以上面 log 包的 New 函数为例，这个函数用于实例化一个 log.Logger 实例，它接受一个 io.Writer 接口类型的参数，返回 *log.Logger。从 New 的实现上来看，传入的 out 参数被作为初值赋值给了 log.Logger 结构体字段 out。

创建模式通过接口，在 NewXXX 函数所在包与接口的实现者所在包之间建立了一个连接。大多数包含接口类型字段的结构体的实例化，都可以使用创建模式实现。这个模式比较容易理解，我们就不再深入了。

3、包装器模式
在基本模式的基础上，当返回值的类型与参数类型相同时，我们能得到下面形式的函数原型：
func YourWrapperFunc(param YourInterfaceType) YourInterfaceType

通过这个函数，我们可以实现对输入参数的类型的包装，并在不改变被包装类型（输入参数类型）的定义的情况下，返回具备新功能特性的、实现相同接口类型的新类型。
这种接口应用模式我们叫它包装器模式，也叫装饰器模式。包装器多用于对输入数据的过滤、变换等操作。

下面就是 Go 标准库中一个典型的包装器模式的应用：
// $GOROOT/src/io/io.go
func LimitReader(r Reader, n int64) Reader { return &LimitedReader{r, n} }
type LimitedReader struct {
    R Reader // underlying reader
    N int64  // max bytes remaining
}
func (l *LimitedReader) Read(p []byte) (n int, err error) {
    // ... ...
}

通过 LimitReader 函数的包装后，我们得到了一个具有新功能特性的 io.Reader 接口的实现类型，也就是 LimitedReader。这个新类型在 Reader 的语义基础上实现了对读取字节个数的限制。

接下来我们再具体看 LimitReader 的一个使用示例：
func main() {
    r := strings.NewReader("hello, gopher!\n")
    lr := io.LimitReader(r, 4)
    if _, err := io.Copy(os.Stdout, lr); err != nil {
        log.Fatal(err)
    }
}

由于包装器模式下的包装函数（如上面的 LimitReader）的返回值类型与参数类型相同，因此我们可以将多个接受同一接口类型参数的包装函数组合成一条链来调用，形式是这样的：
YourWrapperFunc1(YourWrapperFunc2(YourWrapperFunc3(...)))

在上面示例的基础上自定义一个包装函数：CapReader，通过这个函数的包装，我们能得到一个可以将输入的数据转换为大写的 Reader 接口实现：
func CapReader(r io.Reader) io.Reader {
    return &capitalizedReader{r: r}
}
type capitalizedReader struct {
    r io.Reader
}
func (r *capitalizedReader) Read(p []byte) (int, error) {
    n, err := r.r.Read(p)
    if err != nil {
        return 0, err
    }
    q := bytes.ToUpper(p)
    for i, v := range q {
        p[i] = v
    }
    return n, err
}
func main() {
    r := strings.NewReader("hello, gopher!\n")
    r1 := CapReader(io.LimitReader(r, 4))
    if _, err := io.Copy(os.Stdout, r1); err != nil {
        log.Fatal(err)
    }
}

这里，我们将 CapReader 和 io.LimitReader 串在了一起形成一条调用链，这条调用链的功能变为：截取输入数据的前四个字节并将其转换为大写字母。

4、适配器模式
适配器模式不是基本模式的直接衍生模式，但这种模式是后面中间件模式的前提。

适配器模式的核心是适配器函数类型（Adapter Function Type）。
适配器函数类型是一个辅助水平组合实现的“工具”类型。
这里我要再强调一下，它是一个类型。
它可以将一个满足特定函数签名的普通函数，显式转换成自身类型的实例，转换后的实例同时也是某个接口类型的实现者。

最典型的适配器函数类型莫过于 http.HandlerFunc了。这里看一个应用 http.HandlerFunc 的例子：
func greetings(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome!")
}
func main() {
    http.ListenAndServe(":8080", http.HandlerFunc(greetings))
}

这个例子通过 http.HandlerFunc 这个适配器函数类型，将普通函数 greetings 快速转化为满足 http.Handler 接口的类型。
而 http.HandleFunc 这个适配器函数类型的定义是这样的：
// $GOROOT/src/net/http/server.go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
type HandlerFunc func(ResponseWriter, *Request)
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}

经过 HandlerFunc 的适配转化后，我们就可以将它的实例用作实参，传递给接收 http.Handler 接口的 http.ListenAndServe 函数，从而实现基于接口的组合。

5、中间件（Middleware）

中间件（Middleware）这个词的含义可大可小。
在 Go Web 编程中，“中间件”常常指的是一个实现了 http.Handler 接口的 http.HandlerFunc 类型实例。
实质上，这里的中间件就是包装模式和适配器模式结合的产物。

func validateAuth(s string) error {
    if s != "123456" {
        return fmt.Errorf("%s", "bad auth token")
    }
    return nil
}
func greetings(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome!")
}
func logHandler(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        t := time.Now()
        log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t)
        h.ServeHTTP(w, r)
    })
}
func authHandler(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        err := validateAuth(r.Header.Get("auth"))
        if err != nil {
            http.Error(w, "bad auth param", http.StatusUnauthorized)
            return
        }
        h.ServeHTTP(w, r)
    })
}
func main() {
    http.ListenAndServe(":8080", logHandler(authHandler(http.HandlerFunc(greetings))))
}

所谓中间件（如：logHandler、authHandler）本质就是一个包装函数（支持链式调用），但它的内部利用了适配器函数类型（http.HandlerFunc），将一个普通函数（比如例子中的几个匿名函数）转型为实现了 http.Handler 的类型的实例。

运行这个示例，并用 curl 工具命令对其进行测试，可以得到下面结果：
$curl http://localhost:8080
bad auth param
$curl -H "auth:123456" localhost:8080/
Welcome!

总结：
尽量避免使用空接口作为函数参数类型

Go 语言之父 Rob Pike 曾说过：空接口不提供任何信息（The empty interface says nothing）。
我们应该怎么理解这句话的深层含义呢？
在 Go 语言中，一方面你不用像 Java 那样显式声明某个类型实现了某个接口，但另一方面，你又必须声明这个接口，这又与接口在 Java 等静态类型语言中的工作方式更加一致。

这种不需要类型显式声明实现了某个接口的方式，可以让种类繁多的类型与接口匹配，包括那些存量的、并非由你编写的代码以及你无法编辑的代码（比如：标准库）。Go 的这种处理方式兼顾安全性和灵活性，其中，这个安全性就是由 Go 编译器来保证的，而为编译器提供输入信息的恰恰是接口类型的定义。

比如我们看下面的接口：
// $GOROOT/src/io/io.go
type Reader interface {
  Read(p []byte) (n int, err error)
}

建议广大 Gopher 尽可能地抽象出带有一定行为契约的接口，并将它作为函数参数类型，尽量不要使用可以“逃过”编译器类型安全检查的空接口类型（interface{}）。

在这方面，Go 标准库已经为我们作出了”表率“。
全面搜索标准库后，你可以发现以interface{}为参数类型的方法和函数少之甚少。
不过，也还有，使用interface{}作为参数类型的函数或方法主要有两类：
1、容器算法类，比如：container 下的 heap、list 和 ring 包、sort 包、sync.Map 等；

2、格式化 / 日志类，比如：fmt 包、log 包等。

这些使用interface{}作为参数类型的函数或方法都有一个共同特点，就是它们面对的都是未知类型的数据，所以在这里使用具有“泛型”能力的interface{}类型。
我们也可以理解为是在 Go 语言尚未支持泛型的这个阶段的权宜之计。
等 Go 泛型落地后，很多场合下 interface{}就可以被泛型替代了。
*/
