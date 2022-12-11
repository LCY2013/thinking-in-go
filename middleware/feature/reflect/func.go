package reflect

import (
	"errors"
	"reflect"
)

type Interface interface {
	Method()
}

// Struct 指针与结构体
type Struct struct {
}

var _ Interface = Struct{}  // ok
var _ Interface = &Struct{} // ok

func (s Struct) Method() {

}

var _ Interface = &Str{} // ok
//var _ Interface = Str{} // bad

type Str struct {
}

func (s *Str) Method() {

}

// IterateFunc 输出给定结构体的方法信息
func IterateFunc(val any) (map[string]FuncInfo, error) {
	if val == nil {
		return nil, errors.New("invalid input")
	}

	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Ptr && typ.Kind() != reflect.Struct {
		return nil, errors.New("val must be a pointer or a struct")
	}

	num := typ.NumMethod()
	ret := make(map[string]FuncInfo, num)
	for i := 0; i < num; i++ {
		fn := typ.Method(i)
		numIn := fn.Type.NumIn()
		params := make([]reflect.Value, 0, numIn)
		// 将第一个值给到对应的参数列表，第一个参数类似与其他语言的this或者self
		params = append(params, reflect.ValueOf(val))
		in := make([]reflect.Type, 0, numIn)
		// 从1开始，避开第一个参数
		for j := 1; j < numIn; j++ {
			inType := fn.Type.In(j)
			in = append(in, inType)
			params = append(params, reflect.Zero(inType))
		}

		// 开始方法调用
		call := fn.Func.Call(params)

		// 出参处理
		numOut := fn.Type.NumOut()
		out := make([]reflect.Type, 0, numOut)
		res := make([]any, 0, numOut)
		for j := 0; j < numOut; j++ {
			outType := fn.Type.Out(j)
			out = append(out, outType)
			res = append(res, call[j].Interface())
		}

		ret[fn.Name] = FuncInfo{
			Name: fn.Name,

			In:  in,
			Out: out,

			Ret: res,
		}
	}
	return ret, nil
}

type FuncInfo struct {
	Name string

	In  []reflect.Type
	Out []reflect.Type

	// 反射获取到的结果
	Ret []any
}
