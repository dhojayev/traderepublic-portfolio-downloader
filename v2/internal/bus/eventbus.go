//go:generate go tool mockgen -source=eventbus.go -destination eventbus_mock_gen.go -package=bus

package bus

import (
	"log/slog"
	"sync"
)

type Event struct {
	Topic string
	ID    string
	Data  any
}

func NewEvent(topic, id string, data any) Event {
	return Event{
		Topic: topic,
		ID:    id,
		Data:  data,
	}
}

type EventHandler func(Event)

type EventBusInterface interface {
	Subscribe(string, EventHandler)
	Publish(Event)
}

type EventBus struct {
	subscribers map[string][]EventHandler
	mu          sync.RWMutex
}

func New() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

func (b *EventBus) Subscribe(topic string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[topic] = append(b.subscribers[topic], handler)
}

func (b *EventBus) Publish(event Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if handlers, found := b.subscribers[event.Topic]; found {
		for _, handler := range handlers {
			go handler(event)
		}
	}

	slog.Debug("event published", "topic", event.Topic, "id", event.ID)
}
