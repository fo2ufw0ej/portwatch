// Package plugboard provides a simple event bus for broadcasting port change
// events to multiple registered handlers.
package plugboard

import "sync"

// Handler is a function that receives a named event and an arbitrary payload.
type Handler func(event string, payload any)

// Bus is a thread-safe event bus.
type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

// New returns an initialised Bus.
func New() *Bus {
	return &Bus{handlers: make(map[string][]Handler)}
}

// Subscribe registers h to be called whenever event is published.
func (b *Bus) Subscribe(event string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[event] = append(b.handlers[event], h)
}

// Publish delivers payload to every handler subscribed to event.
// Handlers are called synchronously in registration order.
func (b *Bus) Publish(event string, payload any) {
	b.mu.RLock()
	handlers := make([]Handler, len(b.handlers[event]))
	copy(handlers, b.handlers[event])
	b.mu.RUnlock()

	for _, h := range handlers {
		h(event, payload)
	}
}

// Unsubscribe removes all handlers for event.
func (b *Bus) Unsubscribe(event string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.handlers, event)
}

// Topics returns the list of events that have at least one subscriber.
func (b *Bus) Topics() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	topics := make([]string, 0, len(b.handlers))
	for k := range b.handlers {
		topics = append(topics, k)
	}
	return topics
}
