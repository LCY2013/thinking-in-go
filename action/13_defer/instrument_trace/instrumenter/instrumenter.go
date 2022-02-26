package instrumenter

type Instrumenter interface {
	// Instrument 接受一个 Go 源文件路径，返回注入了 Trace 函数的新源文件内容以及一个 error 类型值，作为错误状态标识。
	// 默认提供了一种自动注入 Trace 函数的实现，那就是 ast.instrumenter，它注入 Trace 的实现原理是这样的：
	// 抽象语法树（abstract syntax tree，AST）是源代码的抽象语法结构的树状表现形式，树上的每个节点都表示源代码中的一种结构。
	// 因为 Go 语言是开源编程语言，所以它的抽象语法树的操作包也和语言一起开放给了 Go 开发人员，我们可以基于 Go 标准库以及Go 实验工具库提供的 ast 相关包，快速地构建基于 AST 的应用，这里的 ast.instrumenter 就是一个应用 AST 的典型例子。
	// 一旦通过 ast 相关包解析 Go 源码得到相应的抽象语法树后，便可以操作这棵语法树，并按我们的逻辑在语法树中注入 Trace 函数，最后再将修改后的抽象语法树转换为 Go 源码，就完成了整个自动注入的工作了。
	Instrument(string) ([]byte, error)
}
