package eventbus_test

import (
	"sync/atomic"
	"testing"

	"github.com/dmokel/portwatch/internal/eventbus"
)

func TestPublish_CallsSubscriber(t *testing.T) {
	b := eventbus.New()
	var called int32
	b.Subscribe(eventbus.EventPortOpened, func(e eventbus.Event) {
		atomic.AddInt32(&called, 1)
	})
	b.Publish(eventbus.Event{Type: eventbus.EventPortOpened, Payload: 8080})
	if atomic.LoadInt32(&called) != 1 {
		t.Fatalf("expected handler called once, got %d", called)
	}
}

func TestPublish_MultipleSubscribers(t *testing.T) {
	b := eventbus.New()
	var count int32
	for i := 0; i < 3; i++ {
		b.Subscribe(eventbus.EventScanDone, func(e eventbus.Event) {
			atomic.AddInt32(&count, 1)
		})
	}
	b.Publish(eventbus.Event{Type: eventbus.EventScanDone})
	if atomic.LoadInt32(&count) != 3 {
		t.Fatalf("expected 3 calls, got %d", count)
	}
}

func TestPublish_NoSubscribers_NoError(t *testing.T) {
	b := eventbus.New()
	// Should not panic.
	b.Publish(eventbus.Event{Type: eventbus.EventPortClosed})
}

func TestPublish_WrongType_NotCalled(t *testing.T) {
	b := eventbus.New()
	var called int32
	b.Subscribe(eventbus.EventPortOpened, func(e eventbus.Event) {
		atomic.AddInt32(&called, 1)
	})
	b.Publish(eventbus.Event{Type: eventbus.EventPortClosed})
	if atomic.LoadInt32(&called) != 0 {
		t.Fatalf("handler should not have been called")
	}
}

func TestTopics_ReturnsList(t *testing.T) {
	b := eventbus.New()
	b.Subscribe(eventbus.EventPortOpened, func(e eventbus.Event) {})
	b.Subscribe(eventbus.EventScanDone, func(e eventbus.Event) {})
	topics := b.Topics()
	if len(topics) != 2 {
		t.Fatalf("expected 2 topics, got %d", len(topics))
	}
}

func TestReset_ClearsSubscribers(t *testing.T) {
	b := eventbus.New()
	var called int32
	b.Subscribe(eventbus.EventPortOpened, func(e eventbus.Event) {
		atomic.AddInt32(&called, 1)
	})
	b.Reset()
	b.Publish(eventbus.Event{Type: eventbus.EventPortOpened})
	if atomic.LoadInt32(&called) != 0 {
		t.Fatalf("expected no calls after reset")
	}
}
