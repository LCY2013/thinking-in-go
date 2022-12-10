package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lcy2013/instrument_trace/instrumenter"
	"github.com/lcy2013/instrument_trace/instrumenter/ast"
)

var (
	wrote bool
)

func init() {
	flag.BoolVar(&wrote, "w", false, "write result to (source) file instead of stdout")
}

func usage() {
	fmt.Println("instrument [-w] xxx.go")
	flag.PrintDefaults()
}

func main() {
	fmt.Println(os.Args)
	flag.Usage = usage
	// 解析命令行参数
	flag.Parse()

	if len(os.Args) < 2 {
		usage()
		return
	}

	var file string
	if len(os.Args) == 3 {
		file = os.Args[2]
	}

	if len(os.Args) == 2 {
		file = os.Args[1]
	}

	// 对文件扩展名校验
	if filepath.Ext(file) != ".go" {
		usage()
		return
	}

	// 声明instrumenter.Instrumenter接口类型变量
	var ins instrumenter.Instrumenter
	// 创建以ast方式实现Instrumenter接口的ast.instrumenter实例
	ins = ast.New("github.com/lcy2013/instrument_trace", "trace", "Trace")
	// 向go源码所有函数注入Trace函数
	newSrc, err := ins.Instrument(file)
	if err != nil {
		panic(err)
	}

	if newSrc == nil {
		// add nothing to the source file. no change
		fmt.Printf("no trace added for %s\n", file)
		return
	}

	// 将生成的新代码内容输出到stdout上
	if !wrote {
		fmt.Println(newSrc)
		return
	}

	// 将生成的新代码写回go源文件
	if err = ioutil.WriteFile(file, newSrc, 0666); err != nil {
		fmt.Printf("write %s error: %v\n", file, err)
		return
	}

	// 写入文件完成
	fmt.Printf("instument trace for %s ok\n", file)
}
