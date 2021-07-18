package v1

import "errors"

// Named Type Error

// errorString create named type for our new error type
type errorString string

// Implement the error interface
func (err errorString) Error() string {
	return string(err)
}

// New creates interface value of type error
func New(text string) error {
	return errorString(text)
}

// ErrNamedType 自定义的 EOF 异常
var ErrNamedType = New("EOF")

// ErrStructType 基础API创建 EOF 异常
var ErrStructType = errors.New("EOF")
