package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"
)

type save interface {
	Save(writer io.Writer, data []byte) error
}

type FileSave struct {
	save
}

func (FileSave) Save(writer io.Writer, data []byte) error {
	fmt.Printf("%s", string(data))
	return nil
}

// 在这段代码中，我们通过 bytes.NewBuffer 创建了一个 *bytes.Buffer 类型变量 buf，由于 bytes.Buffer 实现了 Write 方法，进而实现了 io.Writer 接口，我们可以合法地将变量 buf 传递给 Save 函数。之后我们可以从 buf 中取出 Save 函数写入的数据内容与预期的数据做比对，就可以达到对 Save 函数进行单元测试的目的了。在整个测试过程中，我们不需要创建任何磁盘文件或建立任何网络连接。
func TestSave(t *testing.T) {
	b := make([]byte, 0, 128)
	buf := bytes.NewBuffer(b)
	data := []byte("hello, golang")
	fileSave := FileSave{}
	err := fileSave.Save(buf, data)
	if err != nil {
		t.Errorf("want nil, actual %s", err.Error())
	}
	saved := buf.Bytes()
	if !reflect.DeepEqual(saved, data) {
		t.Errorf("want %s, actual %s", string(data), string(saved))
	}
}
