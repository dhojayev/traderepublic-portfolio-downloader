package publisher

import "log/slog"

type Publisher struct {
	logger      *slog.Logger
	subscribers map[string]chan []byte
}

func NewPublisher(logger *slog.Logger) *Publisher {
	return &Publisher{
		logger:      logger,
		subscribers: make(map[string]chan []byte),
	}
}

func (p *Publisher) Subscribe(topic string) <-chan []byte {
	ch := make(chan []byte)
	p.subscribers[topic] = ch

	return ch
}

func (p *Publisher) Publish(msg []byte, topic string) {
	ch, ok := p.subscribers[topic]
	if !ok {
		p.logger.Error("channel not found", "topic", topic)

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
