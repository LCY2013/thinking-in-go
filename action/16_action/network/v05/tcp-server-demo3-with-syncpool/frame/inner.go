package frame

import (
	"encoding/binary"
	"errors"
	"io"
)

var ErrShortWrite = errors.New("short write")
var ErrShortRead = errors.New("short read")

type innerFrameCodec struct {
}

func (inner *innerFrameCodec) Encode(writer io.Writer, payload FramePayload) error {
	var f = payload
	var totalLen int32 = int32(len(payload)) + 4

	// binary.Read 或 Write 会根据参数的宽度，读取或写入对应的字节个数的字节，这里 totalLen 使用 int32，那么 Read 或 Write 只会操作数据流中的 4 个字节；
	err := binary.Write(writer, binary.BigEndian, &totalLen)
	if err != nil {
		return err
	}

	// write the frame payload to outbound stream
	n, err := writer.Write([]byte(f))
	if err != nil {
		return err
	}

	if n != len(payload) {
		return ErrShortWrite
	}

	return nil
}

func (inner *innerFrameCodec) Decode(reader io.Reader) (FramePayload, error) {
	var totalLen int32

	// binary.Read 或 Write 会根据参数的宽度，读取或写入对应的字节个数的字节，这里 totalLen 使用 int32，那么 Read 或 Write 只会操作数据流中的 4 个字节；
	err := binary.Read(reader, binary.BigEndian, &totalLen)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, totalLen-4)
	n, err := io.ReadFull(reader, buf)
	if n != int(totalLen-4) {
		return nil, ErrShortRead
	}
	return FramePayload(buf), nil
}

func NewInnerFrameCodec() StreamFrameCodec {
	return &innerFrameCodec{}
}
