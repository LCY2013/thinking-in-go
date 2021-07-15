package search

import (
	"encoding/json"
	"os"
)

// 用于读取json数据文件

const dataFile = "data/data.json"

// Feed 包含我们需要处理的数据源的信息
type Feed struct {
	Name string `json:"site"`
	URI  string `json:"link"`
	Type string `json:"type"`
}

// RetrieveFeeds 读取并反序列化源数据文件
func RetrieveFeeds() ([]*Feed, error) {
	// 打开文件
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	// 当函数返回时，关闭文件
	defer file.Close()

	// 将文件解码到一个切片里
	// 这个切片的每一项都指向一个 Feed 类型的指针
	var feeds []*Feed
	err = json.NewDecoder(file).Decode(&feeds)

	// 这个函数不需要检测错误，调用者会做这件事
	return feeds, err
}
