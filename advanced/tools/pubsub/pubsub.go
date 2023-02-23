// Package pubsub implements a simple multi-topic pub-sub library.
package pubsub

import (
	"sync"
	"time"
)

type (
	subscriber chan any         // 订阅者为一个通道
	topicFunc  func(v any) bool // 主题的一个过滤器
)

// Publisher 发布者对象
type Publisher struct {
	mu          sync.RWMutex             // 读写锁
	buffer      int                      // 订阅队列的缓存大小
	timeout     time.Duration            // 发布超时时间
	subscribers map[subscriber]topicFunc // 订阅者信息
}

// NewPublisher 构建一个发布者对象，可以设置发布超时时间和缓存队列的长度
func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
	return &Publisher{
		timeout:     publishTimeout,
		buffer:      buffer,
		subscribers: make(map[subscriber]topicFunc),
	}
}

// Subscribe 添加新的订阅者，订阅所有主题
func (p *Publisher) Subscribe() chan any {
	return p.SubscribeTopic(nil)
}

// SubscribeTopic 添加一个新的订阅者，订阅过滤器筛选后的主题
func (p *Publisher) SubscribeTopic(topic topicFunc) chan any {
	ch := make(chan any, p.buffer)
	p.mu.Lock()
	defer p.mu.Unlock()

	p.subscribers[ch] = topic
	return ch
}

// Evict 退出订阅
func (p *Publisher) Evict(sub chan any) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.subscribers, sub)
	close(sub)
}

// Publish 发布
func (p *Publisher) Publish(v any) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var wg sync.WaitGroup

	for sub, topic := range p.subscribers {
		wg.Add(1)
		go func(sub subscriber, topic topicFunc) {
			p.sendTopic(sub, topic, v, &wg)
		}(sub, topic)
	}

	wg.Wait()
}

// Close 关闭发布者对象，同时关闭所有的订阅者通道
func (p *Publisher) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}

// sendTopic 发送主题，可以容忍一定的超时
func (p *Publisher) sendTopic(
	sub subscriber,
	topic topicFunc,
	v any,
	wg *sync.WaitGroup) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}

	select {
	case sub <- v:
	case <-time.After(p.timeout):
	}
}
