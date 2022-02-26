package ast

import (
	"bytes"
	"fmt"
	instrumented "github.com/lcy2013/instrument_trace/instrumenter"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
)

func New(traceImport, tracePkg, traceFunc string) instrumented.Instrumenter {
	return &instrumenter{traceImport, tracePkg, traceFunc}
}

type instrumenter struct {
	traceImport string
	tracePkg    string
	traceFunc   string
}

func (a instrumenter) Instrument(filename string) ([]byte, error) {
	fset := token.NewFileSet()
	curAst, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", filename, err)
	}

	// 如果整个源码都不包含函数声明，则无需注入操作，直接返回。
	if !hasFuncDecl(curAst) {
		return nil, nil
	}

	// 在AST上添加包导入语句
	astutil.AddImport(fset, curAst, a.traceImport)

	// 向AST上所有函数注入Trace函数
	a.addDeferTraceIntoFuncDecls(curAst)

	buf := &bytes.Buffer{}
	err = format.Node(buf, fset, curAst)
	if err != nil {
		return nil, fmt.Errorf("error formatting new code: %w", err)
	}
	return buf.Bytes(), nil
}

// addDeferTraceIntoFuncDecls 向AST上所有函数注入Trace函数
func (a instrumenter) addDeferTraceIntoFuncDecls(curAst *ast.File) {
	// 遍历所有声明语句
	for _, decl := range curAst.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if ok {
			// 如果是函数声明，则注入Trace
			// 遍历语法树上所有声明语句，如果是函数声明，就调用 instrumenter 的 addDeferStmt 方法进行注入，如果不是，就直接返回。
			a.addDeferStmt(fd)
		}
	}
}

// addDeferStmt 注入Trace
// addDeferStmt 函数体略长，但逻辑也很清晰，就是先判断函数是否已经注入了 Trace，如果有，则略过；如果没有，就构造一个 Trace 语句节点，并将它插入到 AST 中。
// Instrument 的最后一步就是将注入 Trace 后的 AST 重新转换为 Go 代码，这就是期望得到的带有 Trace 特性的 Go 代码了。
func (a instrumenter) addDeferStmt(fd *ast.FuncDecl) (added bool) {
	stmts := fd.Body.List

	// 判断"defer trace.Trace()()"语句是否已经存在
	for _, stmt := range stmts {
		ds, ok := stmt.(*ast.DeferStmt)
		if !ok {
			// 如果不是defer语句，则继续for循环
			continue
		}
		// 如果是defer语句，则要进一步判断是否是defer trace.Trace()()
		ce, ok := ds.Call.Fun.(*ast.CallExpr)
		if !ok {
			continue
		}
		se, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		x, ok := se.X.(*ast.Ident)
		if !ok {
			continue
		}

		if x.Name == a.tracePkg && se.Sel.Name == a.traceFunc {
			// defer trace.Trace()() 已经存在直接返回
			return false
		}
	}

	// 没有找到"defer trace.Trace()()"，注入一个新的跟踪语句
	// 在AST上构造一个defer trace.Trace()()
	ds := &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: a.tracePkg,
					},
					Sel: &ast.Ident{
						Name: a.traceFunc,
					},
				},
			},
		},
	}

	newList := make([]ast.Stmt, len(stmts)+1)
	copy(newList[1:], stmts)
	// 注入新构造的defer语句
	newList[0] = ds
	fd.Body.List = newList
	return true
}

// hasFuncDecl 查询源码中是否有函数声明
func hasFuncDecl(astFile *ast.File) bool {
	if len(astFile.Decls) == 0 {
		return false
	}

	for _, decl := range astFile.Decls {
		_, ok := decl.(*ast.FuncDecl)
		if ok {
			return true
		}
	}

	return false
}
