package mapstruct

import "fmt"

// NewErrIndexOutOfRange 创建一个代表下标超出范围的错误
func NewErrIndexOutOfRange(length int, index int) error {
	return fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", length, index)
}

// NewErrInvalidType 创建一个代表类型转换失败的错误
func NewErrInvalidType(want, got string) error {
	return fmt.Errorf("ekit: 类型转换失败，want:%s, got:%s", want, got)
}
