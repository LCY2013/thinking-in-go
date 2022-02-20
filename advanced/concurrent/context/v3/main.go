package main

import (
	"context"
	"fmt"
	"time"
)

const shortDuration = 1 * time.Millisecond

func main() {
	d := time.Now().Add(shortDuration)
	ctx, cancelFunc := context.WithDeadline(context.Background(), d)

	defer cancelFunc()
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}
}
