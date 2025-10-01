//go:generate go tool mockgen -source=publisher.go -destination publisher_mock_gen.go -package=traderepublic

package traderepublic

import (
	"log/slog"
	"sync"
)

type PublisherInterface interface {
	Subscribe(topic string) <-chan []byte
	Publish(msg []byte, topic string)
	Close(topic string)
}

type Publisher struct {
	subscribers map[string]chan []byte
	mu          *sync.RWMutex
}

func NewPublisher() *Publisher {
	return &Publisher{
		subscribers: make(map[string]chan []byte),
		mu:          &sync.RWMutex{},
	}
}

func (p *Publisher) Subscribe(topic string) <-chan []byte {
	p.mu.Lock()
	defer p.mu.Unlock()

	ch := make(chan []byte)
	p.subscribers[topic] = ch

	return ch
}

func (p *Publisher) Publish(msg []byte, topic string) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	ch, ok := p.subscribers[topic]
	if !ok {
		slog.Error("channel not found", "topic", topic)

		return
	}

	ch <- msg
}

func (p *Publisher) Close(topic string) {
	ch, ok := p.subscribers[topic]
	if !ok {
		return
	}

	close(ch)
	delete(p.subscribers, topic)
}
