package main

import (
	"errors"
	"fmt"
	"fufeng.org/advanced/errors/named_type/v2"
)

func main() {
	if v2.ErrType == errors.New("EOF") {
		fmt.Println("Error : ", v2.ErrType)
	}
}
