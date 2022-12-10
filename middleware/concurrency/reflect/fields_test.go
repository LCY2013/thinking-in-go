package reflect

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/LCY2013/thinking-in-go/middleware/concurrency/reflect/types"
	"github.com/stretchr/testify/assert"
)

// Test_iterateFields TDD
func Test_iterateFields(t *testing.T) {
	user := &types.User{
		Name: "zhangsan",
	}
	user1 := &user

	tests := []struct {
		name string

		input any

		want    map[string]any
		wantErr error
	}{
		{
			// 普通结构体
			name: "normal struct",
			input: types.User{
				Name: "zhangsan",
				//age:  18,
			},
			want: map[string]any{
				"Name": "zhangsan",
				"age":  0,
			},
		}, {
			// 指针结构体
			name: "pointer struct",
			input: &types.User{
				Name: "zhangsan",
			},
			want: map[string]any{
				"Name": "zhangsan",
				"age":  0,
			},
		}, {
			// 多重指针
			name:  "multiple pointer struct",
			input: &user1,
			want: map[string]any{
				"Name": "zhangsan",
				"age":  0,
			},
		}, {
			// 非法输入
			name:    "slice",
			input:   []string{},
			wantErr: errors.New("input must be a struct"),
		}, {
			// 非法输入指针
			name:    "pointer to map",
			input:   &map[string]string{},
			wantErr: errors.New("input must be a struct"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := iterateFields(tt.input)
			if err != nil && !assert.Equal(t, err, tt.wantErr) {
				t.Errorf("iterateFields() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			//assert.Equal(t, got, tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("iterateFields() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test_SetField TDD
func Test_SetField(t *testing.T) {
	tests := []struct {
		name string

		entity any
		field  string
		value  any

		wantErr error
	}{
		{
			name:    "struct",
			entity:  types.User{},
			field:   "Name",
			value:   "zhangsan",
			wantErr: errors.New("entity must be a pointer to a struct"),
		},
		{
			name:    "private name",
			entity:  &types.User{},
			field:   "age",
			value:   18,
			wantErr: errors.New(fmt.Sprintf("field %s must be exported", "age")),
		},
		{
			name:    "invalid field name",
			entity:  &types.User{},
			field:   "invalid field name",
			value:   "zhangsan",
			wantErr: errors.New(fmt.Sprintf("field %s not found", "invalid field name")),
		}, {
			name:   "pass",
			entity: &types.User{},
			field:  "Name",
			value:  "zhangsan",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetField(tt.entity, tt.field, tt.value)
			if err != nil && !assert.Equal(t, err, tt.wantErr) {
				t.Errorf("iterateFields() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			t.Logf("test name[%s]: %s", tt.name, tt.entity)
		})
	}
}

// TestLocalVariable 测试值是否可以通过反射被修改
func TestLocalVariable(t *testing.T) {
	i := 1
	val := reflect.ValueOf(i)
	ptrVal := reflect.ValueOf(&i)
	t.Log(val.CanSet())           // false
	t.Log(ptrVal.CanSet())        // false
	t.Log(ptrVal.Elem().CanSet()) // true
	ptrVal.Elem().Set(reflect.ValueOf(2))
	t.Log(i) // 2
}

// BenchmarkFieldValueOfIndexOrName 对比字段通过索引查询和通过名称查询的性能差异
// index 性能比 name 好很多
// go test -bench=. -benchmem -benchtime=3s -run=none
// run=none 忽略输出
func BenchmarkFieldValueOfIndexOrName(b *testing.B) {
	user := types.User{}
	val := reflect.ValueOf(user)
	ptrVal := reflect.ValueOf(&user).Elem()

	b.Run("val by index", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// 获取随机其中一个字段信息
			_ = val.Field(1)
		}
	})
	b.Run("ptrVal by index", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// 获取随机其中一个字段信息
			_ = ptrVal.Field(1)
		}
	})

	b.Run("val by name", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// 获取随机其中一个字段信息
			_ = val.FieldByName("Name")
		}
	})
	b.Run("ptrVal by name", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// 获取随机其中一个字段信息
			_ = ptrVal.FieldByName("Name")
		}
	})
}
