package message

import "sync"

// InMemoryQueue implements a queue that stores messages in memory. Messages
// can be enqueued concurrently but the returned iterator is not safe for
// concurrent access.
type InMemoryQueue struct {
	mu sync.Mutex
	msgs []Message
	
	latchedMsg Message
}

// NewInMemoryQueue creates a new in-memory queue instance. This function can serve as a Queuefactory.
func NewInMemoryQueue() Queue {
	return new(InMemoryQueue)
}

// Enqueue implements Queue 
func (q *InMemoryQueue) Enqueue(msg Message) error {
	q.mu.Lock()
	q.msgs = append(q.msgs, msg)
	q.mu.Unlock()
	return nil
}

// PendingMessages implements Queue
func (q *InMemoryQueue) PendingMessages() bool {
	q.mu.Lock()
	pending := len(q.msgs) != 0
	q.mu.Unlock()
	return pending
}

// DiscardMessages implements Queue
func (q *InMemoryQueue) DiscardMessages() error {
	q.mu.Lock()
	q.msgs = q.msgs[:0]
	q.latchedMsg = nil
	q.mu.Unlock()
	return nil
}

// Close implements Queue
func (*InMemoryQueue) Close() error { return nil }

// Messages implements Queue.
func (q *InMemoryQueue) Messages() Iterator { return q }

// Next implements Iterator.
func (q *InMemoryQueue) Next() bool {
	q.mu.Lock()
	qLen := len(q.msgs)
	if 	qLen == 0 {
		q.mu.Unlock()
		return false
	}
	
	// Dequeue message from the tail of the queue.
	q.latchedMsg = q.msgs[qLen-1]
	q.msgs = q.msgs[:qLen-1]
	q.mu.Unlock()
	return true
}

// Message implements Iterator.
func (q *InMemoryQueue) Message() Message {
	q.mu.Lock()
	msg := q.latchedMsg
	q.mu.Unlock()
	return msg
}

// Error implements Iterator.
func (q *InMemoryQueue) Error() error { return nil }
