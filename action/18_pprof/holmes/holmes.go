package main

import (
	"mosn.io/holmes"
	"net/http"
	"time"
)

/*
https://github.com/mosn/holmes
*/

func init() {
	http.HandleFunc("/make1gb", make1gbslice)
	go http.ListenAndServe(":10003", nil)
}

// main Dump goroutine when goroutine number spikes
func main() {
	h, _ := holmes.New(
		holmes.WithCollectInterval("5s"), // 指定 2s 为区间监控进程的资源占用率， 线上建议设置大于10s的采样区间。
		holmes.WithDumpPath("/tmp"),
		holmes.WithTextDump(),
		holmes.WithDumpToLogger(true),
		holmes.WithGoroutineDump(10, 25, 2000, 10*1000, time.Minute),
	)
	h.EnableMemDump().
		EnableGoroutineDump().
		EnableGCHeapDump().
		EnableCPUDump().
		EnableThreadDump().
		Start()
	time.Sleep(time.Hour)
}

func make1gbslice(wr http.ResponseWriter, req *http.Request) {
	var a = make([]byte, 1073741824)
	_ = a
}
