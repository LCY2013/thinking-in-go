package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/lcy2013/tcp-server-demo2/frame"
	"github.com/lcy2013/tcp-server-demo2/metrics"
	"github.com/lcy2013/tcp-server-demo2/packet"
)

func main() {
	// pprof for go
	go func() {
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			return
		}
	}()

	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Printf("listen error: [%+v]\n", err)
		return
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept error: [%+v]\n", err)
			break
		}

		// start a new goroutine to handle the new connection.
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	// 连接建立，ClientConnected加1
	metrics.ClientConnected.Inc()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("handleConn error: [%+v]\n", err)
		}
		// 连接断开，ClientConnected减1
		metrics.ClientConnected.Dec()
	}()
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("conn error: [%+v]\n", err)
		}
	}(conn)

	codec := frame.NewInnerFrameCodec()
	for {
		// read from the connection

		// decode the frame to get the payload
		// the payload is undecoded packet
		framePayload, err := codec.Decode(conn)
		if err != nil {
			fmt.Printf("handleConn : frame decode error: [%+v]\n", err)
			return
		}
		// 收到并解码一个消息请求，ReqRecvTotal消息计数器加1
		metrics.ReqRecvTotal.Add(1)

		// do something with the packet
		ackFramePayload, err := handlePacket(framePayload)
		if err != nil {
			fmt.Printf("handlePacket : handle packet error : [%+v]\n", err)
			return
		}

		// write ack frame to the connection
		err = codec.Encode(conn, ackFramePayload)
		if err != nil {
			fmt.Printf("handleConn: frame encode error: [%+v]\n", err)
			return
		}

		// 返回响应后，RspSendTotal消息计数器减1
		metrics.RespSendTotal.Add(1)
	}
}

func handlePacket(payload frame.FramePayload) (ackFramePayload []byte, err error) {
	var p packet.Packet
	p, err = packet.Decode(payload)
	if err != nil {
		fmt.Printf("handleConn: packet decode error: [%+v]\n", err)
		return
	}
	switch p.(type) {
	case *packet.Submit:
		submit := p.(*packet.Submit)
		fmt.Printf("recv submit: id = %s, payload=%s\n", submit.ID, string(submit.Payload))
		submitAck := &packet.SubmitAck{
			ID:     submit.ID,
			Result: 0,
		}
		ackFramePayload, err := packet.Encode(submitAck)
		if err != nil {
			fmt.Printf("handleConn: packet encode error: [%+v]\n", err)
			return nil, err
		}
		return ackFramePayload, nil
	default:
		return nil, fmt.Errorf("unknown packet type")
	}
}
