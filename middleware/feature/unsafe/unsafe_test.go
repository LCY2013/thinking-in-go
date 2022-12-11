package unsafe

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsafeAccessor_GetIntField(t *testing.T) {
	tests := []struct {
		name string

		entity    any
		fieldName string

		want    int
		wantErr error
	}{
		{
			name: "normal case",
			entity: &User{
				Age: 18,
			},
			fieldName: "Age",
			want:      18,
		},
		{
			name: "invalid field",
			entity: &User{
				Age: 18,
			},
			fieldName: "age",
			want:      18,
			wantErr:   fmt.Errorf("invalid field: %s", "age"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessor, err := NewUnsafeAccessor(tt.entity)
			if err != nil {
				assert.Equal(t, err, tt.wantErr)
				return
			}

			got, err := accessor.GetIntField(tt.fieldName)
			if err != nil {
				assert.Equal(t, err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want)

			gotAny, err := accessor.GetAnyField(tt.fieldName)
			if err != nil {
				assert.Equal(t, err, tt.wantErr)
				return
			}
			assert.Equal(t, gotAny, tt.want)
		})
	}
}

func TestUnsafeAccessor_SetIntField(t *testing.T) {

	tests := []struct {
		name string

		entity    *User
		fieldName string
		val       int

		wantErr error
	}{
		{
			name:      "normal case",
			entity:    &User{},
			fieldName: "Age",
			val:       18,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessor, err := NewUnsafeAccessor(tt.entity)
			if err != nil {
				assert.Equal(t, err, tt.wantErr)
				return
			}
			/*err = accessor.SetIntField(tt.fieldName, tt.val)
			if err != nil {
				assert.Equal(t, err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.val, tt.entity.Age)
			*/
			err = accessor.SetAnyField(tt.fieldName, tt.val)
			if err != nil {
				assert.Equal(t, err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.val, tt.entity.Age)
		})
	}
}

type User struct {
	Age int
}
