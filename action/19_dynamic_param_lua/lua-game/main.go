package main

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

func main() {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoFile("./lua-game/game.lua"); err != nil {
		panic(err)
	}

	// 调用Lua中的函数
	L.CallByParam(lua.P{
		Fn:      L.GetGlobal("main"),
		NRet:    1,
		Protect: true,
	})

	// 获取返回值
	ret := L.Get(-1)
	if b, ok := ret.(lua.LBool); ok {
		fmt.Println("Game Over:", b)
	}
}
