package reflect

import (
	"github.com/LCY2013/thinking-in-go/middleware/feature/reflect/types"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterateFuncs(t *testing.T) {
	type args struct {
		val any
	}

	user := types.User{
		Name: "zhangsan",
	}

	tests := []struct {
		name    string
		args    args
		want    map[string]FuncInfo
		wantErr error
		Assert  assert.ComparisonAssertionFunc
	}{
		{
			name: "normal struct",
			args: args{
				val: user,
			},
			want: map[string]FuncInfo{
				"GetAge": {
					Name: "GetAge",
					In:   []reflect.Type{},
					Out:  []reflect.Type{reflect.TypeOf(int64(0))},
					Ret:  []any{int64(0)},
				},
			},
			wantErr: nil,
			Assert:  assert.Equal,
		},
		{
			name: "pointer",
			args: args{
				val: &user,
			},
			want: map[string]FuncInfo{
				"GetAge": {
					Name: "GetAge",
					In:   []reflect.Type{},
					Out:  []reflect.Type{reflect.TypeOf(int64(0))},
					Ret:  []any{int64(0)},
				},
				"ChangeName": {
					Name: "ChangeName",
					In:   []reflect.Type{reflect.TypeOf("")},
					Out:  []reflect.Type{},
					Ret:  []any{},
				},
				// 已经被 ChangeName 重新置为 ""
				"GetName": {
					Name: "GetName",
					In:   []reflect.Type{},
					Out:  []reflect.Type{reflect.TypeOf("")},
					Ret:  []any{""},
				},
			},
			wantErr: nil,
			Assert:  assert.Equal,
		},
		{
			name: "normal struct use pointer method",
			args: args{
				val: user,
			},
			want: map[string]FuncInfo{
				"GetAge": {
					Name: "GetAge",
					In:   []reflect.Type{},
					Out:  []reflect.Type{reflect.TypeOf(int64(0))},
					Ret:  []any{int64(0)},
				},
				"ChangeName": {
					Name: "ChangeName",
					In:   []reflect.Type{reflect.TypeOf("lisi")},
					Out:  []reflect.Type{},
					Ret:  []any{},
				},
				// 已经被 ChangeName 重新置为 ""
				"GetName": {
					Name: "GetName",
					In:   []reflect.Type{},
					Out:  []reflect.Type{reflect.TypeOf("")},
					Ret:  []any{""},
				},
			},
			wantErr: nil,
			Assert:  assert.NotEqual,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IterateFunc(tt.args.val)
			if err != nil && !assert.Equal(t, err, tt.wantErr) {
				t.Error(err)
			}
			tt.Assert(t, tt.want, got, "IterateFunc(%v)", tt.args.val)
		})
	}
}
