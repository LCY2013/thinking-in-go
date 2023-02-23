package pubsub

import (
	"fufeng.org/advanced/async"
	"github.com/sirupsen/logrus"
	"strings"
	"testing"
	"time"
)

var (
	p *Publisher
)

func init() {
	p = NewPublisher(100*time.Millisecond, 100)
}

func TestPubSub(t *testing.T) {
	defer p.Close()

	allSub := p.Subscribe()
	golangSub := p.SubscribeTopic(func(v any) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})

	p.Publish("hello world")
	p.Publish("hello golang")

	async.GO(func() {
		for msg := range allSub {
			logrus.
				WithField("sub", "all").
				Info(msg)
		}
	})

	async.GO(func() {
		for msg := range golangSub {
			logrus.
				WithField("sub", "golang").
				Info(msg)
		}
	})

	p.Publish("hello world")
	p.Publish("hello golang")

	// 运行一段时间后退出
	time.Sleep(3 * time.Second)
}
