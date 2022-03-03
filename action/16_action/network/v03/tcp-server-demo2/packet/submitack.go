package packet

import "bytes"

type SubmitAck struct {
	ID     string
	Result uint8
}

func (ack *SubmitAck) Decode(pktBody []byte) error {
	ack.ID = string(pktBody[0:8])
	ack.Result = uint8(pktBody[8])
	return nil
}

func (ack *SubmitAck) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(ack.ID[:8]), []byte{ack.Result}}, nil), nil
}
