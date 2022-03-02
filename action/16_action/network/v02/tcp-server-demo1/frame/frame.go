package frame

import "io"

type FramePayload []byte

type StreamFrameCodec interface {
	Encode(writer io.Writer, payload FramePayload) error // data -> frame，并写入io.Writer
	Decode(reader io.Reader) (FramePayload, error)       // 从io.Reader中提取frame payload，并返回给上层
}
