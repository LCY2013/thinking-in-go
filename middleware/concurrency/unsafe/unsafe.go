package unsafe

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

type FieldAccessor interface {
	Field(field string) (int, error)
	SetField(field string, value int) error
}

type FieldMeta struct {
	typ reflect.Type
	// offset go的组合这些，或者复杂类型字段的时候，它的含义表达相当于对外层结构体的偏移量
	offset uintptr
}

type UnsafeAccessor struct {
	Fields     map[string]FieldMeta
	entityAddr unsafe.Pointer
}

func NewUnsafeAccessor(entity any) (*UnsafeAccessor, error) {
	if entity == nil {
		return nil, errors.New("invalid entity")
	}

	val := reflect.ValueOf(entity)
	typ := reflect.TypeOf(entity)
	//val.UnsafeAddr()
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errors.New("invalid entity")
	}
	fields := make(map[string]FieldMeta, typ.Elem().NumField())
	elemType := typ.Elem()
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		fields[field.Name] = FieldMeta{
			typ:    field.Type,
			offset: field.Offset,
		}
	}
	return &UnsafeAccessor{entityAddr: val.UnsafePointer(), Fields: fields}, nil
}

// GetIntField 明确类型可以使用该类函数提升性能
func (ua *UnsafeAccessor) GetIntField(fieldName string) (int, error) {
	fieldMeta, ok := ua.Fields[fieldName]
	if !ok {
		return 0, fmt.Errorf("invalid field: %s", fieldName)
	}
	ptr := unsafe.Pointer(uintptr(ua.entityAddr) + fieldMeta.offset)
	if ptr == nil {
		return 0, fmt.Errorf("invalid address of the field: %s", fieldName)
	}
	res := *(*int)(ptr)
	return res, nil
}

func (ua *UnsafeAccessor) SetIntField(fieldName string, val int) error {
	fieldMeta, ok := ua.Fields[fieldName]
	if !ok {
		return fmt.Errorf("invalid field: %s", fieldName)
	}
	ptr := unsafe.Pointer(uintptr(ua.entityAddr) + fieldMeta.offset)
	if ptr == nil {
		return fmt.Errorf("invalid address of the field: %s", fieldName)
	}
	*(*int)(ptr) = val
	return nil
}

// GetAnyField 任意类型支持，在不清楚类型情况下使用
func (ua *UnsafeAccessor) GetAnyField(fieldName string) (any, error) {
	fieldMeta, ok := ua.Fields[fieldName]
	if !ok {
		return nil, fmt.Errorf("invalid field: %s", fieldName)
	}
	ptr := unsafe.Pointer(uintptr(ua.entityAddr) + fieldMeta.offset)
	if ptr == nil {
		return nil, fmt.Errorf("invalid address of the field: %s", fieldName)
	}
	// 计算任意字段的地址
	res := reflect.NewAt(fieldMeta.typ, unsafe.Pointer(uintptr(ua.entityAddr)+fieldMeta.offset)).Elem()
	// 最终返回
	return res.Interface(), nil
}

func (ua *UnsafeAccessor) SetAnyField(fieldName string, val any) error {
	fieldMeta, ok := ua.Fields[fieldName]
	if !ok {
		return fmt.Errorf("invalid field: %s", fieldName)
	}
	ptr := unsafe.Pointer(uintptr(ua.entityAddr) + fieldMeta.offset)
	if ptr == nil {
		return fmt.Errorf("invalid address of the field: %s", fieldName)
	}
	// 计算任意字段的地址
	res := reflect.NewAt(fieldMeta.typ, unsafe.Pointer(uintptr(ua.entityAddr)+fieldMeta.offset)).Elem()
	// 最终设置新的值
	if res.CanSet() {
		res.Set(reflect.ValueOf(val))
	}
	return nil
}
