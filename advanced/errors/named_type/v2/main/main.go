package main

import (
	"errors"
	"fmt"
)

func main() {
	if v2.ErrType == errors.New("EOF") {
		fmt.Println("Error : ", v2.ErrType)
	}
}
