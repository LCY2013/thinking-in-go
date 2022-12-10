package main

import (
	"log"
	"os"
	"time"

	v1 "fufeng.org/concurrent_mode/runner/v1"
)

// 使用通道来监视程序运行的时间，以及程序在运行时间过长时如何终止程序

// timeout 规定了必须在多少秒内完成
const timeout = 3 * time.Second

// main 程序的主入口
func main() {
	log.Println("Starting work.")

	// 为本次执行分配超时时间
	r := v1.New(timeout)

	// 加入要执行的任务
	r.Add(createTask(), createTask(), createTask())

	// 执行任务并处理结果
	if err := r.Start(); err != nil {
		switch err {
		case v1.ErrInterrupt:
			log.Println("Terminating due to interrupt.")
			os.Exit(2)
		case v1.ErrTimeout:
			log.Println("Terminating due to timeout.")
			os.Exit(1)
		}
	}

	log.Println("Process ended.")
}

// createTask 创建一个根据id休眠指定秒数的任务示例
func createTask() func(int) {
	return func(id int) {
		log.Printf("Processor - Task #%d.", id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}
