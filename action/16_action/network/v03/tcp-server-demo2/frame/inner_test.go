package frame

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestEncode(t *testing.T) {
	codec := NewInnerFrameCodec()
	buf := make([]byte, 0, 128)
	rw := bytes.NewBuffer(buf)

	err := codec.Encode(rw, []byte("hello"))
	if err != nil {
		t.Errorf("want nil, actual %s", err.Error())
	}

	// 验证Encode正确性
	var totalLen int32
	err = binary.Read(rw, binary.BigEndian, &totalLen)

	if err != nil {
		t.Errorf("want nil, actual %s", err.Error())
	}

	if totalLen != 9 {
		t.Errorf("want 9, actual %d", totalLen)
	}

	left := rw.Bytes()
	if string(left) != "hello" {
		t.Errorf("want hello, actual %s", string(left))
	}
}

func TestDecode(t *testing.T) {
	codec := NewInnerFrameCodec()
	data := []byte{0x0, 0x0, 0x0, 0x9, 'h', 'e', 'l', 'l', 'o'}
	rw := bytes.NewBuffer(data)
	payload, err := codec.Decode(rw)
	if err != nil {
		t.Errorf("want nil, actual %s", string(payload))
	}

	if string(payload) != "hello" {
		t.Errorf("want hello, actual %s", string(payload))
	}
}
