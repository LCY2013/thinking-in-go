package main

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"strconv"
)

func main() {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoFile("./lua-config/config.lua"); err != nil {
		panic(err)
	}

	// 从lua中获取配置
	host := L.GetGlobal("host").String()
	port, _ := strconv.Atoi(L.GetGlobal("port").String())

	fmt.Printf("host=%s, port=%d", host, port)
}
