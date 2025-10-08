package app

import (
	"sync"
	"sync/atomic"
)

// Broadcast provides a simple pub/sub hub for a specific value type.
// Subscribers receive values through buffered channels. Publishing never blocks;
// if a subscriber channel is full the value is dropped for that subscriber only.
type Broadcast[T any] struct {
	mu          sync.RWMutex
	subscribers map[int64]chan T
	nextID      atomic.Int64
	closed      bool
}

// NewBroadcast creates a new Broadcast instance.
func NewBroadcast[T any]() *Broadcast[T] {
	return &Broadcast[T]{
		subscribers: make(map[int64]chan T),
	}
}

// Subscribe registers a new subscriber with the provided buffer size.
// It returns a receive-only channel and a cancel function to unsubscribe.
func (b *Broadcast[T]) Subscribe(buffer int) (<-chan T, func()) {
	ch := make(chan T, buffer)

	b.mu.Lock()
	if b.closed {
		close(ch)
		b.mu.Unlock()
		return ch, func() {}
	}

	id := b.nextID.Add(1)
	b.subscribers[id] = ch
	b.mu.Unlock()

	cancel := func() {
		b.mu.Lock()
		if sub, ok := b.subscribers[id]; ok {
			delete(b.subscribers, id)
			close(sub)
		}
		b.mu.Unlock()
	}

	return ch, cancel
}

// Publish sends value to all current subscribers.
func (b *Broadcast[T]) Publish(value T) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return
	}

	for _, ch := range b.subscribers {
		select {
		case ch <- value:
		default:
		}
	}
}

// Close closes the broadcast and all subscriber channels.
func (b *Broadcast[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	b.closed = true
	for id, ch := range b.subscribers {
		close(ch)
		delete(b.subscribers, id)
	}
}
