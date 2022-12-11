package unsafe

import (
	"fmt"
	"reflect"
	"unsafe"
)

// go 里面对齐是俺字长，64位机器字长为8，32位机器字长为4

// PrintFieldOffset 打印字段偏移量
// 探讨go内存布局
// 接受一个结构体作为内存布局
func PrintFieldOffset(entity any) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Struct {
		return
	}
	fmt.Sprintf("%s: offset is %d\n", typ.Name(), unsafe.Sizeof(entity))
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fmt.Printf("%s: %d\n", field.Name, field.Offset)
	}
}
