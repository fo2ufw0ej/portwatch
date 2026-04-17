package plugboard

import (
	"sort"
	"sync"
	"testing"
)

func TestPublish_CallsSubscribers(t *testing.T) {
	b := New()
	var got []string
	b.Subscribe("ports.opened", func(_ string, payload any) {
		got = append(got, payload.(string))
	})
	b.Publish("ports.opened", "8080")
	b.Publish("ports.opened", "9090")
	if len(got) != 2 || got[0] != "8080" || got[1] != "9090" {
		t.Fatalf("unexpected got: %v", got)
	}
}

func TestPublish_NoSubscribers_NoError(t *testing.T) {
	b := New()
	b.Publish("ports.closed", 443) // should not panic
}

func TestPublish_MultipleSubscribers(t *testing.T) {
	b := New()
	var mu sync.Mutex
	count := 0
	for i := 0; i < 3; i++ {
		b.Subscribe("scan.done", func(_ string, _ any) {
			mu.Lock()
			count++
			mu.Unlock()
		})
	}
	b.Publish("scan.done", nil)
	if count != 3 {
		t.Fatalf("expected 3 calls, got %d", count)
	}
}

func TestUnsubscribe_RemovesHandlers(t *testing.T) {
	b := New()
	called := false
	b.Subscribe("ports.opened", func(_ string, _ any) { called = true })
	b.Unsubscribe("ports.opened")
	b.Publish("ports.opened", nil)
	if called {
		t.Fatal("handler should not have been called after Unsubscribe")
	}
}

func TestTopics_ReturnsList(t *testing.T) {
	b := New()
	b.Subscribe("a", func(_ string, _ any) {})
	b.Subscribe("b", func(_ string, _ any) {})
	topics := b.Topics()
	sort.Strings(topics)
	if len(topics) != 2 || topics[0] != "a" || topics[1] != "b" {
		t.Fatalf("unexpected topics: %v", topics)
	}
}

func TestTopics_EmptyBus(t *testing.T) {
	b := New()
	if len(b.Topics()) != 0 {
		t.Fatal("expected no topics on empty bus")
	}
}
