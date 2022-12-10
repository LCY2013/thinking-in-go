package sync

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnceClose_Close(t *testing.T) {
	type fields struct {
		close sync.Once
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "close",
			wantErr: assert.NoError,
			fields: fields{
				close: sync.Once{},
			},
		},
	}

	oc := &OnceClose{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, oc.Close(), fmt.Sprintf("Close()"))
		})
	}
}
