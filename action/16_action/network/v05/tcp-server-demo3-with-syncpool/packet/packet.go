package packet

import (
	"bytes"
	"fmt"
	"sync"
)

// Packet协议定义
/*
### packet header
1 byte: commandID

### submit packet

8字节 ID 字符串
任意字节 payload

### submit ack packet

8字节 ID 字符串
1字节 result
*/

const (
	CommandConn   = iota + 0x01 // 0x01
	CommandSubmit               // 0x02
)

const (
	CommandConnAck   = iota + 0x80 // 0x80
	CommandSubmitAck               // 0x81
)

type Packet interface {
	Decode([]byte) error     // []byte -> struct
	Encode() ([]byte, error) // struct -> []byte
}

var SubmitPool = sync.Pool{
	New: func() interface{} {
		return &Submit{}
	},
}

func Decode(packet []byte) (Packet, error) {
	commandID := packet[0]
	pktBody := packet[1:]

	switch commandID {
	case CommandConn:
		return nil, nil
	case CommandConnAck:
		return nil, nil
	case CommandSubmit:
		// 从Pool池获取submit对象
		submit := SubmitPool.Get().(*Submit)
		err := submit.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return submit, nil
	case CommandSubmitAck:
		submitAck := SubmitAck{}
		err := submitAck.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &submitAck, nil
	default:
		return nil, fmt.Errorf("decode unknown commandID [%d]", commandID)
	}
}

func Encode(p Packet) ([]byte, error) {
	var commandID uint8
	var pktBody []byte
	var err error

	switch p.(type) {
	case *Submit:
		commandID = CommandSubmit
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *SubmitAck:
		commandID = CommandSubmitAck
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("encode unknown commandID [%d]", commandID)
	}
	return bytes.Join([][]byte{[]byte{commandID}, pktBody}, nil), nil
}
