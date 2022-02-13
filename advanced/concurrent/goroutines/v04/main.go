package main

import (
	"fmt"
	"time"
)

func search(term string) (string, error) {
	time.Sleep(time.Millisecond * 200)
	return "some value", nil
}

func process(term string) error {
	record, err := search(term)
	if err != nil {
		return err
	}
	fmt.Printf("Received: %s\n", record)
	return nil
}

func main() {

}
