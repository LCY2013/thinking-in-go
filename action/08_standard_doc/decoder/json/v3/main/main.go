package main

import (
	"encoding/json"
	"fmt"
	"fufeng.org/standard/decoder/json/v3"
	"log"
)

func main() {
	// 将字符串反序列化成变量
	var c map[string]interface{}
	err := json.Unmarshal([]byte(v2.JSON), &c)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	fmt.Println("Name:", c["name"])
	fmt.Println("Title:", c["title"])
	fmt.Println("Contact")
	fmt.Println("Home:", c["contact"].(map[string]interface{})["home"])
	fmt.Println("Cell:", c["contact"].(map[string]interface{})["cell"])
}
