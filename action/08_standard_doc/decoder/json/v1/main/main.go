package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	v12 "fufeng.org/standard/decoder/json/v1"
)

func main() {
	uri := "http://ajax.googleapis.com/ajax/services/search/web?v=1.0&rsz=8&q=golang"
	// 向谷歌发起搜索
	resp, err := http.Get(uri)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	defer resp.Body.Close()

	// 将 json 响应码解码成结构体类型
	var gr v12.GgResponse
	err = json.NewDecoder(resp.Body).Decode(&gr)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	fmt.Println(gr)
}
