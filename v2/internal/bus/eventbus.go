package bus

import "sync"

type Event struct {
	Name string
	Data any
}

type EventHandler func(Event)

type EventBus struct {
	subscribers map[string][]EventHandler
	mu          sync.RWMutex
}

func New() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

func (b *EventBus) Subscribe(eventName string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[eventName] = append(b.subscribers[eventName], handler)
}

func (b *EventBus) Publish(topic string, data any) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if handlers, found := b.subscribers[topic]; found {
		event := Event{Name: topic, Data: data}

		for _, handler := range handlers {
			go handler(event)
		}
	}
}
