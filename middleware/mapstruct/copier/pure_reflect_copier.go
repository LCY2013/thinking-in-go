package copier

import "reflect"

// CopyTo 复制结构体，递归实现，src,dest 必须是结构体指针
func CopyTo(src, dest any) error {
	srcPrtTyp := reflect.TypeOf(src)
	if srcPrtTyp.Kind() != reflect.Ptr {
		return newErrTypeError(srcPrtTyp)
	}
	srcTyp := srcPrtTyp.Elem()
	if srcTyp.Kind() != reflect.Struct {
		return newErrTypeError(srcTyp)
	}

	destPrtTyp := reflect.TypeOf(dest)
	if destPrtTyp.Kind() != reflect.Ptr {
		return newErrTypeError(destPrtTyp)
	}
	destTyp := destPrtTyp.Elem()
	if destTyp.Kind() != reflect.Struct {
		return newErrTypeError(destTyp)
	}

	srcVal := reflect.ValueOf(src).Elem()
	destVal := reflect.ValueOf(dest).Elem()
	return copyStruct(srcTyp, srcVal, destTyp, destVal)
}

func copyStruct(srcTyp reflect.Type,
	srcVal reflect.Value,
	destTyp reflect.Type,
	destVal reflect.Value) error {
	srcFieldNameIndex := make(map[string]int, 0)
	for i := 0; i < srcTyp.NumField(); i++ {
		fieldTyp := srcTyp.Field(i)
		if !fieldTyp.IsExported() {
			continue
		}
		srcFieldNameIndex[fieldTyp.Name] = i
	}

	for i := 0; i < destTyp.NumField(); i++ {
		fieldTyp := destTyp.Field(i)
		if !fieldTyp.IsExported() {
			continue
		}
		if idx, ok := srcFieldNameIndex[fieldTyp.Name]; ok {
			if err := copyStructField(srcTyp, srcVal, idx, destTyp, destVal, i); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyStructField(srcTyp reflect.Type,
	srcVal reflect.Value,
	srcIdx int,
	destTyp reflect.Type,
	destVal reflect.Value,
	destIdx int) error {
	srcFieldTyp := srcTyp.Field(srcIdx)
	destFieldTyp := destTyp.Field(destIdx)
	if srcFieldTyp.Type.Kind() != destFieldTyp.Type.Kind() {
		return newErrKindNotMatchError(srcFieldTyp.Type.Kind(), destFieldTyp.Type.Kind(), srcFieldTyp.Name)
	}

	srcFieldVal := srcVal.Field(srcIdx)
	destFieldVal := destVal.Field(destIdx)
	if srcFieldVal.Type().Kind() == reflect.Ptr {
		if srcFieldVal.IsNil() {
			return nil
		}
		if destFieldVal.IsNil() {
			destFieldVal.Set(reflect.New(destFieldTyp.Type.Elem()))
		}
		return copyData(srcFieldTyp.Type.Elem(), srcFieldVal.Elem(), destFieldTyp.Type.Elem(), destFieldVal.Elem(), srcFieldTyp.Name)
	}
	return copyData(srcFieldTyp.Type, srcFieldVal, destFieldTyp.Type, destFieldVal, srcFieldTyp.Name)
}

func copyData(srcTyp reflect.Type, srcVal reflect.Value, destTyp reflect.Type, destVal reflect.Value, fieldName string) error {
	if srcTyp.Kind() == reflect.Ptr {
		return newErrMultiPointer(fieldName)
	}

	if srcTyp.Kind() != destTyp.Kind() {
		return newErrKindNotMatchError(srcTyp.Kind(), destTyp.Kind(), fieldName)
	}

	if isShadowCopyType(srcTyp.Kind()) {
		// 内置类型，但不匹配，如别名、map、slice
		if srcTyp != destTyp {
			return newErrTypeNotMatchError(srcTyp, destTyp, fieldName)
		}
		if destVal.CanSet() {
			destVal.Set(srcVal)
		}
	} else if srcTyp.Kind() == reflect.Struct {
		return copyStruct(srcTyp, srcVal, destTyp, destVal)
	}

	return nil
}
