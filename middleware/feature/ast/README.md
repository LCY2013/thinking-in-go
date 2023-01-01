## AST 编程

- [GORM Gen 子项目](https://github.com/go-gorm/gen)

- [依赖注入 Google Wire](https://github.com/google/wire)

### ast.Node接口
ast.Node 接口代表了 AST 的节点。任何一个Go语言在AST里面都代表了一个节点，多个语句之间可以组成一个更加复杂的节点。

三个关键子类：
- 表达式 Expr: 例如各种内置类型、方法调用、数组索引等
- 语句 Stmt: 如赋值语句、if语句等
- 声明 Decl: 各种声明

#### GenDecl 
GenDecl 代表通用声明，一般是:
- 类型声明，以type开头的
- 变量声明，以var开头的
- 常量声明，以const开头的
- import声明，以import开头的

#### FuncDecl
FuncDecl 是方法声明:
- Doc 是文档。注意，不是所有的注释都是文档，要符合Go规范的文档
- Recv 接收器
- Name 方法名
- Type 方法签名：里面包含泛型参数、参数类型、返回值

#### StructType
StructType 代表了一个结构体声明
- Fields: 字段。实际上，定义在该结构体上的方法，也放在这里。




