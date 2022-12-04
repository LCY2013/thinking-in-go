package reflect

import (
	"errors"
	"fmt"
	"github.com/LCY2013/thinking-in-go/middleware/concurrency/reflect/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
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
