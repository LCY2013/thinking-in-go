package main

import "fmt"

type Err struct {
}

func (e Err) Error() string {
	return ""
}

func E() error {
	var err *Err
	return err
	//return nil
}

func main() {
	if err := E(); err != nil {
		fmt.Println("err", err)
	}
}
