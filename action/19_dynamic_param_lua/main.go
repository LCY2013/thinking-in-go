package main

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"os"
)

func main() {
	L := lua.NewState()
	defer L.Close()

	// 将参数传递到Lua解释器
	L.SetGlobal("argc", lua.LNumber(len(os.Args)))
	table := L.NewTable()
	for i, arg := range os.Args {
		L.SetTable(table, lua.LNumber(i), lua.LString(arg))
	}
	L.SetGlobal("argv", table)

	// 执行Lua脚本
	if err := L.DoFile("script.lua"); err != nil {
		panic(err)
	}

	// 获取结果并输出
	result := L.GetGlobal("result")
	fmt.Println(result.String())
}
