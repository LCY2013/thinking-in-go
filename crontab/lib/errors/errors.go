package errors

import "errors"

var (
	ErrLockAlreadyRequired = errors.New("锁被占用")
)
