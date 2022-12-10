package main

import (
	"errors"
	"fmt"

	v1 "fufeng.org/advanced/errors/named_type/v1"
)

func main() {
	// 相等，具体可以参考 官方 示例如下
	if v1.ErrNamedType == v1.New("EOF") {
		fmt.Println("Named Type Error")
	}

	// 不想等
	if v1.ErrStructType == errors.New("EOF") {
		fmt.Println("Struct Type Error")
	}
}
