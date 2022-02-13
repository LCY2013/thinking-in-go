package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type recordType struct {
	record string
	err    error
}

func search(term string) (string, error) {
	time.Sleep(time.Millisecond * 200)
	return "some value", nil
}

func process(term string, ch chan recordType) error {
	// 超时context
	timeout, _ := context.WithTimeout(context.Background(), time.Millisecond*205)
	select {
	case <-timeout.Done():
		return errors.New("search canceled")
	case result := <-ch:
		if result.err != nil {
			return result.err
		}
		fmt.Printf("Received: %s\n", result.record)
		return nil
	}
}

func main() {
	// 模拟跟目录
	rootDir := "/test"
	// 模拟异步调用
	//ch := make(chan recordType)
	ch := make(chan recordType, 1)
	// 这里存在goroutine 泄漏，因为ch 是一个不带缓冲的通道，修复可以设置ch带一个缓冲值
	go func() {
		record, err := search(rootDir)
		ch <- recordType{record, err}
	}()

	// 开始处理
	err := process(rootDir, ch)
	if err != nil {
		fmt.Printf("process error: [%+v]", err)
	}
}
