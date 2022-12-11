package unsafe

import (
	"github.com/LCY2013/thinking-in-go/middleware/feature/unsafe/types"
	"testing"
)

func TestPrintFieldOffset(t *testing.T) {
	tests := []struct {
		name   string
		entity any
	}{
		{
			name:   "user",
			entity: types.User{},
		},
		{
			name:   "userV1",
			entity: types.UserV1{},
		},
		{
			name:   "userV2",
			entity: types.UserV2{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintFieldOffset(tt.entity)
		})
	}
}
