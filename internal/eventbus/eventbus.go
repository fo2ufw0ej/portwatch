// Package eventbus provides a simple typed event bus for broadcasting
// port change events to multiple consumers within portwatch.
package eventbus

import "sync"

// EventType identifies the kind of event.
type EventType string

const (
	EventPortOpened EventType = "port.opened"
	EventPortClosed EventType = "port.closed"
	EventScanDone   EventType = "scan.done"
)

// Event carries an event type and an optional payload.
type Event struct {
	Type    EventType
	Payload any
}

// Handler is a function that receives an Event.
type Handler func(Event)

// Bus is a simple synchronous event bus.
type Bus struct {
	mu       sync.RWMutex
	subs     map[EventType][]Handler
}

// New returns an initialised Bus.
func New() *Bus {
	return &Bus{subs: make(map[EventType][]Handler)}
}

// Subscribe registers a handler for the given event type.
func (b *Bus) Subscribe(t EventType, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[t] = append(b.subs[t], h)
}

// Publish dispatches an event to all registered handlers synchronously.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	handlers := make([]Handler, len(b.subs[e.Type]))
	copy(handlers, b.subs[e.Type])
	b.mu.RUnlock()
	for _, h := range handlers {
		h(e)
	}
}

// Topics returns the list of event types that have at least one subscriber.
func (b *Bus) Topics() []EventType {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]EventType, 0, len(b.subs))
	for t := range b.subs {
		out = append(out, t)
	}
	return out
}

// Reset removes all subscribers.
func (b *Bus) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs = make(map[EventType][]Handler)
}
