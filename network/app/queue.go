package app

import "sync"

// EventQueue provides a bounded FIFO queue with blocking push/pop operations.
type EventQueue[T any] struct {
	mu     sync.Mutex
	cond   *sync.Cond
	queue  []T
	max    int
	closed bool
}

// NewEventQueue creates a new EventQueue with the provided maximum capacity.
// A max value <= 0 uses 1 as the capacity.
func NewEventQueue[T any](max int) *EventQueue[T] {
	if max <= 0 {
		max = 1
	}
	q := &EventQueue[T]{max: max}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Push enqueues value, blocking if the queue reached capacity.
// Returns false if the queue was closed before the value could be enqueued.
func (q *EventQueue[T]) Push(value T) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.queue) >= q.max && !q.closed {
		q.cond.Wait()
	}
	if q.closed {
		return false
	}
	q.queue = append(q.queue, value)
	q.cond.Signal()
	return true
}

// Pop dequeues the next value, blocking until one is available or the queue closes.
// The ok result is false when the queue has been closed and is empty.
func (q *EventQueue[T]) Pop() (value T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.queue) == 0 && !q.closed {
		q.cond.Wait()
	}
	if len(q.queue) == 0 {
		var zero T
		return zero, false
	}
	value = q.queue[0]
	var zero T
	q.queue[0] = zero
	q.queue = q.queue[1:]
	q.cond.Signal()
	return value, true
}

// Close signals all waiters and prevents future pushes.
func (q *EventQueue[T]) Close() {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return
	}
	q.closed = true
	q.queue = nil
	q.cond.Broadcast()
	q.mu.Unlock()
}
