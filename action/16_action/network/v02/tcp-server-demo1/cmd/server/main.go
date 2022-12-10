package main

import (
	"fmt"
	"net"

	"github.com/lcy2013/tcp-server-demo1/frame"
	"github.com/lcy2013/tcp-server-demo1/packet"
)

func main() {
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
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("handleConn error: [%+v]\n", err)
		}
	}()
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("conn error: [%+v]\n", err)
		}
	}(conn)

	codec := frame.NewInnerFrameCodec()
	for {
		//decode the frame to get the payload
		framePayload, err := codec.Decode(conn)
		if err != nil {
			fmt.Printf("handleConn : frame decode error: [%+v]\n", err)
			return
		}

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
