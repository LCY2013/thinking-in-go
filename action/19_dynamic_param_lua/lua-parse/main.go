package main

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

/*
Lua编写高性能的网络协议解析器或者数据转换器
*/

func main() {
	L := lua.NewState()
	defer L.Close()

	// 注册函数到lua中
	L.SetGlobal("parseData", L.NewFunction(parseData))

	// 在lua中调用parseData函数
	L.DoString(`print(parseData("Hello,World!"))`)
}

// 解析数据的函数
func parseData(L *lua.LState) int {
	data := L.CheckString(1)
	fmt.Println("parseData:", data)
	L.Push(lua.LNumber(len(data)))
	return 1
}
