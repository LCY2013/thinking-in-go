package errors

import "errors"

var (
	ErrLockAlreadyRequired = errors.New("锁被占用")

	// ERR_NO_LOCAL_IP_FOUND ipv4 not found error
	ERR_NO_LOCAL_IP_FOUND = errors.New("not found ipv4 address")

	ErrServesConfigNotFound = errors.New("serves config not found")

	ErrSlaveNotFound = errors.New("worker not found")
)
