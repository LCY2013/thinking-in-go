package reflect

import (
	"errors"
	"fmt"
	"reflect"
)

// IterateFields 迭代结构体里面的所有字段值信息，非公开字段以零值填充
func IterateFields(val any) {
	// 负责逻辑
	fields, err := iterateFields(val)

	// 简单逻辑
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range fields {
		fmt.Println(k, v)
	}
}

// iterateFields 返回所有的字段名称
// val 只能是结构体，或者结构体指针，或者多重指针
func iterateFields(input any) (map[string]any, error) {
	// 类型
	typ := reflect.TypeOf(input)
	// 值
	val := reflect.ValueOf(input)

	// 处理指针，拿到指针指向的东西
	// 这里处理多重指针的效果
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	// 如果不是结构体就返回error
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct")
	}

	num := typ.NumField()
	res := make(map[string]any, num)
	for i := 0; i < num; i++ {
		field := typ.Field(i)
		name := field.Name
		fieldVal := val.FieldByName(field.Name)
		if field.IsExported() {
			res[name] = fieldVal.Interface()
		} else {
			// 如果是一个非公开的字段，就用零值填充
			res[name] = reflect.Zero(field.Type).Interface()
		}
	}

	return res, nil
}

// SetField 设置某个结构体的字段的值信息，允许多重指针
func SetField(entity any, field string, newVal any) error {
	val := reflect.ValueOf(entity)
	typ := val.Type()

	// 如果最后的不是结构体，就提示报错
	if typ.Kind() != reflect.Ptr {
		return errors.New("entity must be a pointer to a struct")
	}

	for typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// 如果最后的不是结构体，就提示报错
	if typ.Kind() != reflect.Struct {
		return errors.New("entity must be a struct")
	}

	// 查询到对应的字段信息存不存在
	if fieldTyp, ok := typ.FieldByName(field); !ok {
		return errors.New(fmt.Sprintf("field %s not found", field))
	} else if !fieldTyp.IsExported() {
		return errors.New(fmt.Sprintf("field %s must be exported", field))
	}

	fd := val.FieldByName(field)

	// 判断该字段是否可以设值
	if !fd.CanSet() {
		return errors.New(fmt.Sprintf("field %s can`t set value", field))
	}
	fd.Set(reflect.ValueOf(newVal))
	return nil
}
