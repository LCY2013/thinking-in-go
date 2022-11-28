package sync

import "fmt"

// newErrIndexOutOfRange 创建一个索引越界error
func newErrIndexOutOfRange(length int, index int) error {
	return fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", length, index)
}
