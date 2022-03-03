package frame

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type ReturnErrorWriter struct {
	W  io.Writer
	Wn int // 第几次调用Write返回错误
	wc int // 写操作次数计数
}

func (w *ReturnErrorWriter) Write(data []byte) (n int, err error) {
	w.Wn++
	if w.wc >= w.Wn {
		return 0, errors.New("write error")
	}
	return w.W.Write(data)
}

type ReturnErrorReader struct {
	R  io.Reader
	Rn int // 第几次调用Read返回错误
	rc int // 读操作次数计数
}

func (r *ReturnErrorReader) Read(data []byte) (n int, err error) {
	r.Rn++
	if r.rc >= r.Rn {
		return 0, errors.New("read error")
	}
	return r.R.Read(data)
}

func TestEncodeWithWriteFail(t *testing.T) {
	codec := NewInnerFrameCodec()
	buf := make([]byte, 0, 128)
	w := bytes.NewBuffer(buf)
	// 模拟binary.Write返回错误
	err := codec.Encode(&ReturnErrorWriter{
		W:  w,
		Wn: 1,
	}, []byte("hello"))
	if err == nil {
		t.Errorf("want non-nil, actual nil")
	}

	// 模拟w.Write返回错误
	err = codec.Encode(&ReturnErrorWriter{
		W:  w,
		Wn: 2,
	}, []byte("hello"))
	if err == nil {
		t.Errorf("want non-nil, actual nil")
	}
}

func TestDecodeWithReadFail(t *testing.T) {
	codec := NewInnerFrameCodec()
	data := []byte{0x0, 0x0, 0x0, 0x9, 'h', 'e', 'l', 'l', 'o'}
	// 模拟binary.Read返回错误
	_, err := codec.Decode(&ReturnErrorReader{
		R:  bytes.NewReader(data),
		Rn: 1,
	})
	if err == nil {
		t.Errorf("want non-nil, actual nil")
	}
	// 模拟io.ReadFull返回错误
	_, err = codec.Decode(&ReturnErrorReader{
		R:  bytes.NewReader(data),
		Rn: 2,
	})
	if err == nil {
		t.Errorf("want non-nil, actual nil")
	}
}
