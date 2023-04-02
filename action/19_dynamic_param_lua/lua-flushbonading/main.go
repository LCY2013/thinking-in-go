package main

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

func main() {
	L := lua.NewState()
	defer L.Close()

	// 在Lua中定义函数
	if err := L.DoString(`
        function add(a, b)
            return a + b
        end
    `); err != nil {
		panic(err)
	}

	// 调用Lua中的函数
	L.CallByParam(lua.P{
		Fn:      L.GetGlobal("add"),
		NRet:    1,
		Protect: true,
	}, lua.LNumber(10), lua.LNumber(20))

	// 获取返回值
	ret := L.Get(-1)
	fmt.Println("add result:", ret)
}
