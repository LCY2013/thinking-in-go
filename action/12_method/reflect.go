package util

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNotStructType = errors.New("not a struct type")
)

// SetStringFieldByFieldName 设置string 类型字段的值
func SetStringFieldByFieldName(any interface{}, fieldName, fieldValue string) error {
	defer func() {
		if e := recover(); e == nil {
			fmt.Printf("set field value error %+v\n", any)
		}
	}()
	dynTyp := reflect.TypeOf(any)
	if dynTyp.Kind() == reflect.Ptr {
		dynTyp = dynTyp.Elem()
	}

	if dynTyp.Kind() != reflect.Struct {
		return ErrNotStructType
	}

	dynValue := reflect.ValueOf(any)
	if dynValue.Kind() == reflect.Ptr {
		dynValue = dynValue.Elem()
	}

	fieldValueReflect := dynValue.FieldByName(fieldName)
	if fieldValueReflect.Kind() == reflect.Ptr {
		fieldValueReflect = fieldValueReflect.Elem()
	}
	// 设置新值
	fieldValueReflect.SetString(fieldValue)
	return nil
}

// GetStringFieldByFieldName 获取string 类型字段的值
func GetStringFieldByFieldName(any interface{}, fieldName string) (string, error) {
	defer func() {
		if e := recover(); e == nil {
			fmt.Printf("get field value error")
		}
	}()
	dynTyp := reflect.TypeOf(any)
	if dynTyp.Kind() == reflect.Ptr {
		dynTyp = dynTyp.Elem()
	}

	if dynTyp.Kind() != reflect.Struct {
		fmt.Println("not a struct type")
		return "", ErrNotStructType
	}

	dynValue := reflect.ValueOf(any)
	if dynValue.Kind() == reflect.Ptr {
		dynValue = dynValue.Elem()
	}

	fieldValueReflect := dynValue.FieldByName(fieldName)
	if fieldValueReflect.Kind() == reflect.Ptr {
		fieldValueReflect = fieldValueReflect.Elem()
	}

	return fieldValueReflect.FieldByName(fieldName).String(), nil
}

// SetIntFieldByFieldName 设置int 类型字段的值
func SetIntFieldByFieldName(any interface{}, fieldName string, fieldValue int64) error {
	defer func() {
		if e := recover(); e == nil {
			fmt.Printf("set field value error %+v", any)
		}
	}()
	dynTyp := reflect.TypeOf(any)
	if dynTyp.Kind() == reflect.Ptr {
		dynTyp = dynTyp.Elem()
	}

	if dynTyp.Kind() != reflect.Struct {
		return ErrNotStructType
	}

	dynValue := reflect.ValueOf(any)
	if dynValue.Kind() == reflect.Ptr {
		dynValue = dynValue.Elem()
	}

	fieldValueReflect := dynValue.FieldByName(fieldName)
	if fieldValueReflect.Kind() == reflect.Ptr {
		fieldValueReflect = fieldValueReflect.Elem()
	}
	// 设置新值
	fieldValueReflect.SetInt(fieldValue)
	return nil
}

// GetIntFieldByFieldName 获取int 类型字段的值
func GetIntFieldByFieldName(any interface{}, fieldName string) (int64, error) {
	defer func() {
		if e := recover(); e == nil {
			fmt.Printf("get field value error")
		}
	}()
	dynTyp := reflect.TypeOf(any)
	if dynTyp.Kind() == reflect.Ptr {
		dynTyp = dynTyp.Elem()
	}

	if dynTyp.Kind() != reflect.Struct {
		fmt.Printf("not a struct type")
		return 0, ErrNotStructType
	}

	dynValue := reflect.ValueOf(any)
	if dynValue.Kind() == reflect.Ptr {
		dynValue = dynValue.Elem()
	}

	fieldValueReflect := dynValue.FieldByName(fieldName)
	if fieldValueReflect.Kind() == reflect.Ptr {
		fieldValueReflect = fieldValueReflect.Elem()
	}

	return fieldValueReflect.FieldByName(fieldName).Int(), nil
}

// SetFloatFieldByFieldName 设置float 类型字段的值
func SetFloatFieldByFieldName(any interface{}, fieldName string, fieldValue float64) error {
	defer func() {
		if e := recover(); e == nil {
			fmt.Printf("set field value error %+v", any)
		}
	}()
	dynTyp := reflect.TypeOf(any)
	if dynTyp.Kind() == reflect.Ptr {
		dynTyp = dynTyp.Elem()
	}

	if dynTyp.Kind() != reflect.Struct {
		return ErrNotStructType
	}

	dynValue := reflect.ValueOf(any)
	if dynValue.Kind() == reflect.Ptr {
		dynValue = dynValue.Elem()
	}

	fieldValueReflect := dynValue.FieldByName(fieldName)
	if fieldValueReflect.Kind() == reflect.Ptr {
		fieldValueReflect = fieldValueReflect.Elem()
	}
	// 设置新值
	fieldValueReflect.SetFloat(fieldValue)
	return nil
}

// GetFloatFieldByFieldName 获取float 类型字段的值
func GetFloatFieldByFieldName(any interface{}, fieldName string) (float64, error) {
	defer func() {
		if e := recover(); e == nil {
			fmt.Printf("get field value error")
		}
	}()
	dynTyp := reflect.TypeOf(any)
	if dynTyp.Kind() == reflect.Ptr {
		dynTyp = dynTyp.Elem()
	}

	if dynTyp.Kind() != reflect.Struct {
		fmt.Printf("not a struct type")
		return 0, ErrNotStructType
	}

	dynValue := reflect.ValueOf(any)
	if dynValue.Kind() == reflect.Ptr {
		dynValue = dynValue.Elem()
	}

	fieldValueReflect := dynValue.FieldByName(fieldName)
	if fieldValueReflect.Kind() == reflect.Ptr {
		fieldValueReflect = fieldValueReflect.Elem()
	}

	return fieldValueReflect.FieldByName(fieldName).Float(), nil
}
