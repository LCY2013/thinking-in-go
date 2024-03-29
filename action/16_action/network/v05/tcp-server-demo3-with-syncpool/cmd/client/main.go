package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/lcy2013/tcp-server-demo3-with-syncpool/frame"
	"github.com/lcy2013/tcp-server-demo3-with-syncpool/packet"
	"github.com/lucasepe/codename"
)

// main 客户端启动了 5 个 Goroutine，模拟 5 个并发连接。
// startClient 函数是每个连接的主处理函数。
func main() {
	var wg sync.WaitGroup
	var num = 100
	wg.Add(num)

	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done()
			startClient(i)
		}(i + 1)
	}

	wg.Wait()
}

func startClient(i int) {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Printf("dial error: [%+v]\n", err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("conn error: [%+v]\n", err)
		}
	}(conn)
	fmt.Printf("[client %d]: dial ok\n", i)

	// 生成payload
	rng, err := codename.DefaultRNG()
	if err != nil {
		panic(err)
	}

	codec := frame.NewInnerFrameCodec()
	var counter int
	go func() {
		for {
			// handle ack
			// read from the connection
			errReadDead := conn.SetReadDeadline(time.Now().Add(time.Second * 5))
			if errReadDead != nil {
				fmt.Printf("conn setReadDeadLine: [%+v]\n", errReadDead)
				return
			}
			ackFramePayLoad, errReadDead := codec.Decode(conn)
			if errReadDead != nil {
				if e, ok := errReadDead.(net.Error); ok {
					if e.Timeout() {
						continue
					}
				}
				panic(errReadDead)
			}
			p, errReadDead := packet.Decode(ackFramePayLoad)
			if errReadDead != nil {
				fmt.Printf("packet decode error: [%+v]\n", errReadDead)
				return
			}
			submitAck, ok := p.(*packet.SubmitAck)
			if !ok {
				panic("no submitAck")
			}
			fmt.Printf("[client %d]: the result of submit ack[%s] is %d\n", i, submitAck.ID, submitAck.Result)
		}
	}()

	for {
		// send submit
		counter++
		id := fmt.Sprintf("%08d", counter) // 8 byte string
		payload := codename.Generate(rng, 4)
		s := &packet.Submit{
			ID:      id,
			Payload: []byte(payload),
		}
		framePayload, err := packet.Encode(s)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[client %d]: send submit id = %s, payload=%s, frame length = %d\n",
			i, s.ID, s.Payload, len(framePayload)+4)
		err = codec.Encode(conn, framePayload)
		if err != nil {
			panic(err)
		}
	}
}
