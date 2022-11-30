package sync

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewTaskPool(t *testing.T) {
	type args struct {
		runSize   int
		queueSize int
		content   []string
	}
	tests := []struct {
		name string
		args args
		want *TaskPool
	}{
		{
			name: "test",
			args: args{
				runSize:   3,
				queueSize: 10,
				content:   []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"},
			},
		},
	}
	for _, tt := range tests {
		tp := NewTaskPool(tt.args.runSize, tt.args.queueSize)
		go func() {
			for _, ctx := range tt.args.content {
				temCtx := ctx
				if err := tp.Run(func() {
					time.Sleep(time.Second * 1)
					fmt.Println(temCtx)
				}); err != nil {
					break
				}
			}
		}()
		time.Sleep(time.Second * 1)
		tp.Stop(context.Background())
	}
}
