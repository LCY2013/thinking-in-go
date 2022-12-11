package mapstruct

import "reflect"

type AnyValue struct {
	Val any
	Err error
}

// Int 返回 int 数据
func (av AnyValue) Int() (int, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(int)
	if !ok {
		return 0, NewErrInvalidType("int", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

// IntOrDefault 返回 int 数据，或者默认值
func (a AnyValue) IntOrDefault(def int) int {
	val, err := a.Int()
	if err != nil {
		return def
	}
	return val
}

// Uint 返回 uint 数据
func (av AnyValue) Uint() (uint, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(uint)
	if !ok {
		return 0, NewErrInvalidType("uint", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

// UintOrDefault 返回 uint 数据，或者默认值
func (a AnyValue) UintOrDefault(def uint) uint {
	val, err := a.Uint()
	if err != nil {
		return def
	}
	return val
}

// Int32 返回 int32 数据
func (av AnyValue) Int32() (int32, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(int32)
	if !ok {
		return 0, NewErrInvalidType("int32", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

// Int32OrDefault 返回 int32 数据，或者默认值
func (a AnyValue) Int32OrDefault(def int32) int32 {
	val, err := a.Int32()
	if err != nil {
		return def
	}
	return val
}

// Uint32 返回 uint32 数据
func (av AnyValue) Uint32() (uint32, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(uint32)
	if !ok {
		return 0, NewErrInvalidType("uint32", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

// Uint32OrDefault 返回 uint32 数据，或者默认值
func (a AnyValue) Uint32OrDefault(def uint32) uint32 {
	val, err := a.Uint32()
	if err != nil {
		return def
	}
	return val
}

// Int64 返回 int64 数据
func (av AnyValue) Int64() (int64, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(int64)
	if !ok {
		return 0, NewErrInvalidType("int64", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

// Int64OrDefault 返回 int64 数据，或者默认值
func (a AnyValue) Int64OrDefault(def int64) int64 {
	val, err := a.Int64()
	if err != nil {
		return def
	}
	return val
}

// Uint64 返回 uint64 数据
func (av AnyValue) Uint64() (uint64, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(uint64)
	if !ok {
		return 0, NewErrInvalidType("uint64", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

// Uint64OrDefault 返回 uint64 数据，或者默认值
func (a AnyValue) Uint64OrDefault(def uint64) uint64 {
	val, err := a.Uint64()
	if err != nil {
		return def
	}
	return val
}

// Float32 返回 float32 数据
func (av AnyValue) Float32() (float32, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(float32)
	if !ok {
		return 0, NewErrInvalidType("float32", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

// Float32OrDefault 返回 float32 数据，或者默认值
func (a AnyValue) Float32OrDefault(def float32) float32 {
	val, err := a.Float32()
	if err != nil {
		return def
	}
	return val
}

// Float64 返回 float64 数据
func (av AnyValue) Float64() (float64, error) {
	if av.Err != nil {
		return 0, av.Err
	}
	val, ok := av.Val.(float64)
	if !ok {
		return 0, NewErrInvalidType("float64", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

// Float64OrDefault 返回 float64 数据，或者默认值
func (a AnyValue) Float64OrDefault(def float64) float64 {
	val, err := a.Float64()
	if err != nil {
		return def
	}
	return val
}

// String 返回 string 数据
func (av AnyValue) String() (string, error) {
	if av.Err != nil {
		return "", av.Err
	}
	val, ok := av.Val.(string)
	if !ok {
		return "", NewErrInvalidType("string", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

// StringOrDefault 返回 string 数据，或者默认值
func (a AnyValue) StringOrDefault(def string) string {
	val, err := a.String()
	if err != nil {
		return def
	}
	return val
}

// Bytes 返回 []byte 数据
func (av AnyValue) Bytes() ([]byte, error) {
	if av.Err != nil {
		return nil, av.Err
	}
	val, ok := av.Val.([]byte)
	if !ok {
		return nil, NewErrInvalidType("[]byte", reflect.TypeOf(av.Val).String())
	}
	return val, nil
}

// BytesOrDefault 返回 []byte 数据，或者默认值
func (a AnyValue) BytesOrDefault(def []byte) []byte {
	val, err := a.Bytes()
	if err != nil {
		return def
	}
	return val
}
