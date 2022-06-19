## gomock

gomock 是Go 编程语言的模拟框架。它与 Go 的内置testing包很好地集成，但也可以在其他环境中使用。

### 安装

> go get github.com/golang/mock/mockgen

### 运行 mockgen

mockgen has two modes of operation: source and reflect.

#### Source mode

源模式从源文件生成模拟接口。它通过使用 -source 标志启用。在此模式下可能有用的其他标志是 -imports 和 -aux_files。

例子：

> mockgen -source=foo.go [other options]

#### reflect mode

反射模式通过构建一个使用反射来理解接口的程序来生成模拟接口。它通过传递两个非标志参数来启用：一个导入路径和一个逗号分隔的符号列表。

您可以使用 ”.” 引用当前路径的包。

例子：

```shell
mockgen database/sql/driver Conn,Driver

# Convenient for `go:generate`.
mockgen . Conn,Driver
```

#### Flags

该mockgen命令用于在给定包含要模拟的接口的 Go 源文件的情况下为模拟类生成源代码。它支持以下标志：

- -source：包含要模拟的接口的文件。

- -destination：将生成的源代码写入其中的文件。如果您不设置此项，代码将打印到标准输出。

- -package：用于生成的模拟类源代码的包。如果不设置，包名将mock_与输入文件的包连接。

- -imports：应在生成的源代码中使用的显式导入列表，指定为表单元素的逗号分隔列表 foo=bar/baz，其中bar/baz是要导入foo的包，是生成的源代码中用于包的标识符。

- -aux_files：应查阅以解决例如在不同文件中定义的嵌入式接口的附加文件列表。这被指定为表单元素的逗号分隔列表 foo=bar/baz.go，其中bar/baz.go是源文件，foo是 -source 文件使用的该文件的包名称。

- -build_flags:（仅限反射模式）标志逐字传递给go build.

- -mock_names：生成的模拟的自定义名称列表。这被指定为表单元素的逗号分隔列表 Repository=MockSensorRepository,Endpoint=MockSensorEndpoint，其中 Repository是接口名称和MockSensorRepository所需的模拟名称（模拟工厂方法和模拟记录器将以模拟命名）。如果其中一个接口没有指定自定义名称，则将使用默认命名约定。

- -self_package：生成代码的完整包导入路径。此标志的目的是通过尝试包含自己的包来防止生成代码中的导入循环。如果将 mock 的包设置为其输入之一（通常是主包）并且输出是 stdio，则可能会发生这种情况，因此 mockgen 无法检测到最终输出包。设置此标志将告诉 mockgen 要排除哪个导入。

- -copyright_file: 用于将版权标头添加到生成的源代码的版权文件。

- -debug_parser：仅打印解析器结果。

- -exec_only：（反射模式）如果设置，则执行此反射程序。

- -prog_only：（反射模式）只生成反射程序；将其写入标准输出并退出。

- -write_package_comment：如果为真，则写入包文档注释（godoc）。（默认为真）

有关使用的示例mockgen，请参见sample/目录。在简单的情况下，您只需要-source标志。

如果您在使用反射模式和供应商依赖项时遇到此错误，您可以选择三种解决方法：
1. 使用源模式。
2. 包括一个空的 import import _ "github.com/golang/mock/mockgen/model"。
3. 添加--build_flags=--mod=mod到您的 mockgen 命令。

#### gomock 常用方法
func InOrder(calls ...*Call)  InOrder声明给定调用的调用顺序

```text
type Call struct {
   t TestReporter // for triggering test failures on invalid call setup

   receiver   interface{}  // the receiver of the method call
   method     string       // the name of the method
   methodType reflect.Type // the type of the method
   args       []Matcher    // the args
   origin     string       // file and line number of call setup

   preReqs []*Call // prerequisite calls

   // Expectations
   minCalls, maxCalls int

   numCalls int // actual number made

   // actions are called when this Call is called. Each action gets the args and
   // can set the return values by returning a non-nil slice. Actions run in the
   // order they are created.
   actions []func([]interface{}) []interface{}
}
```
Call表示对mock对象的一个期望调用 

- func (c *Call) After(preReq *Call) *Call After声明调用在preReq完成后执行 

- func (c *Call) AnyTimes() *Call 允许调用0次或多次

- func (c *Call) Do(f interface{}) *Call 声明在匹配时要运行的操作

- func (c *Call) MaxTimes(n int) *Call 设置最大的调用次数为n次

- func (c *Call) MinTimes(n int) *Call 设置最小的调用次数为n次

- func (c *Call) Return(rets ...interface{}) *Call Return声明模拟函数调用返回的值

- func (c *Call) SetArg(n int, value interface{}) *Call SetArg声明使用指针设置第n个参数的值

- func (c *Call) Times(n int) *Call 设置调用的次数为n次

- func NewController(t TestReporter) *Controller 获取控制对象

- func WithContext(ctx context.Context, t TestReporter) (*Controller, context.Context)WithContext返回一个控制器和上下文，如果发生任何致命错误时会取消。

- func (ctrl *Controller) Call(receiver interface{}, method string, args ...interface{}) []interface{} Mock对象调用，不应由用户代码调用。

- func (ctrl *Controller) Finish() 检查所有预计调用的方法是否被调用，每个控制器都应该调用。本函数只应该被调用一次。

- func (ctrl *Controller) RecordCall(receiver interface{}, method string, args ...interface{}) *Call 被mock对象调用，不应由用户代码调用。

- func (ctrl *Controller) RecordCallWithMethodType(receiver interface{}, method string, methodType reflect.Type, args ...interface{}) *Call 被mock对象调用，不应由用户代码调用。

- func Any() Matcher 匹配任意值

- func AssignableToTypeOf(x interface{}) Matcher AssignableToTypeOf是一个匹配器，用于匹配赋值给模拟调用函数的参数和函数的参数类型是否匹配。

- func Eq(x interface{}) Matcher 通过反射匹配到指定的类型值，而不需要手动设置

- func Nil() Matcher  返回nil

- func Not(x interface{}) Matcher  不递归给定子匹配器的结果


