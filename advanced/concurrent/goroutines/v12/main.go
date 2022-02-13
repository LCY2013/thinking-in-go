package main

import (
	"context"
	"log"
	"time"
)

func main() {
	tr := NewTracker()
	go tr.Run()
	_ = tr.Event(context.Background(), "test1")
	_ = tr.Event(context.Background(), "test2")
	_ = tr.Event(context.Background(), "test3")
	time.Sleep(time.Second * 3)
	deadlineCtx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))
	defer cancelFunc()
	tr.Shutdown(deadlineCtx)
}

func NewTracker() *Tracker {
	return &Tracker{ch: make(chan string, 10)}
}

type Tracker struct {
	ch   chan string
	stop chan struct{}
}

func (tracker *Tracker) Event(ctx context.Context, data string) error {
	select {
	case tracker.ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (tracker *Tracker) Run() {
	for data := range tracker.ch {
		time.Sleep(time.Second)
		log.Println(data)
	}
	tracker.stop <- struct{}{}
}

func (tracker *Tracker) Shutdown(ctx context.Context) {
	close(tracker.ch)
	select {
	case <-tracker.stop:
	case <-ctx.Done():
	}
}
