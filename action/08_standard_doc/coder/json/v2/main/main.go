package main

import (
	"encoding/json"
	"fmt"
	v2 "fufeng.org/standard/coder/json/v2"
	"log"
)

func main() {
	// 将字符串反序列化成变量
	var c v2.Contact
	err := json.Unmarshal([]byte(v2.JSON), &c)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	fmt.Println(c)
}
