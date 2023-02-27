package cgo_helper

//#include <stdio.h>
import "C"

type CChar C.char

func (p *CChar) GoString() string {
	return C.GoString((*C.char)(p))
}

func PrintCString(cs *C.char) {
	C.puts(cs)
}

// GenCChar 提供构造函数
func GenCChar(cc string) *C.char {
	return C.CString(cc)
}
