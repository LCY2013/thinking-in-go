package errors

import "errors"

var (
	ERR_LOCK_ALREADY_REQUIRED = errors.New("锁被占用")
)
