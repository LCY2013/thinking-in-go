package context

import (
	"context"
	"testing"
	"time"
)

func TestContextValue(t *testing.T) {
	ctx := context.Background()
	valueCtx := context.WithValue(ctx, "magic", "fufeng")
	name := valueCtx.Value("magic")
	t.Log(name)
}

func TestContextDeadline(t *testing.T) {
	ctx := context.Background()
	//toDoCtx := context.TODO()
	//timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second)
	//defer cancelFunc()
	//deadline, ok := timeoutCtx.Deadline()
	deadline, ok := ctx.Deadline()
	t.Log(deadline, ok)
}

func TestContextCancel(t *testing.T) {
	ctx := context.Background()
	//toDoCtx := context.TODO()
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second)
	time.Sleep(time.Millisecond * 500)
	cancelFunc()
	err := timeoutCtx.Err()
	t.Log(err)
}

func TestContextTimeout(t *testing.T) {
	ctx := context.Background()
	//toDoCtx := context.TODO()
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second)
	defer cancelFunc()

	time.Sleep(time.Second * 2)
	err := timeoutCtx.Err()
	t.Log(err)
}

func SomeBusiness(ctx context.Context) {

}
