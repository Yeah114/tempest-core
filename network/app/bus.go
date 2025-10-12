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
	subscribers map[int64]*broadcastSubscriber[T]
	nextID      atomic.Int64
	closed      bool
}

// NewBroadcast creates a new Broadcast instance.
func NewBroadcast[T any]() *Broadcast[T] {
	return &Broadcast[T]{
		subscribers: make(map[int64]*broadcastSubscriber[T]),
	}
}

// Subscribe registers a new subscriber with the provided buffer size.
// It returns a receive-only channel and a cancel function to unsubscribe.
func (b *Broadcast[T]) Subscribe(buffer int) (<-chan T, func()) {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		ch := make(chan T)
		close(ch)
		return ch, func() {}
	}

	sub := newBroadcastSubscriber[T](buffer)
	id := b.nextID.Add(1)
	b.subscribers[id] = sub
	b.mu.Unlock()

	cancel := func() {
		b.mu.Lock()
		if existing, ok := b.subscribers[id]; ok && existing == sub {
			delete(b.subscribers, id)
			sub.close()
		}
		b.mu.Unlock()
	}

	return sub.channel(), cancel
}

// Publish sends value to all current subscribers.
func (b *Broadcast[T]) Publish(value T) {
	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return
	}
	subs := make([]*broadcastSubscriber[T], 0, len(b.subscribers))
	for _, sub := range b.subscribers {
		subs = append(subs, sub)
	}
	b.mu.RUnlock()

	for _, sub := range subs {
		sub.publish(value)
	}
}

// Close closes the broadcast and all subscriber channels.
func (b *Broadcast[T]) Close() {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return
	}
	b.closed = true
	subs := b.subscribers
	b.subscribers = make(map[int64]*broadcastSubscriber[T])
	b.mu.Unlock()

	for _, sub := range subs {
		sub.close()
	}
}

type broadcastSubscriber[T any] struct {
	queue *EventQueue[T]
	ch    chan T
}

func newBroadcastSubscriber[T any](buffer int) *broadcastSubscriber[T] {
	if buffer <= 0 {
		buffer = 1
	}
	sub := &broadcastSubscriber[T]{
		queue: NewEventQueue[T](buffer),
		ch:    make(chan T, buffer),
	}
	go sub.run()
	return sub
}

func (s *broadcastSubscriber[T]) run() {
	for {
		value, ok := s.queue.Pop()
		if !ok {
			close(s.ch)
			return
		}
		s.ch <- value
	}
}

func (s *broadcastSubscriber[T]) publish(value T) {
	s.queue.Push(value)
}

func (s *broadcastSubscriber[T]) close() {
	s.queue.Close()
}

func (s *broadcastSubscriber[T]) channel() <-chan T {
	return s.ch
}
