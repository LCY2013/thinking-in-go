package main

import (
	"encoding/json"
	"fmt"
	v22 "fufeng.org/standard/decoder/json/v2"
	"log"
)

func main() {
	// 将字符串反序列化成变量
	var c v22.Contact
	err := json.Unmarshal([]byte(v22.JSON), &c)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	fmt.Println(c)
}
