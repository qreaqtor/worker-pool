package queue

import "sync"

type Queue[T any] struct {
	input   chan<- T
	writers chan struct{}

	mu *sync.RWMutex

	items []T
}

func NewQueue[T any](input chan<- T) *Queue[T] {
	return &Queue[T]{
		items:   make([]T, 0),
		input:   input,
		mu:      &sync.RWMutex{},
		writers: make(chan struct{}, 1),
	}
}

func (q *Queue[T]) Append(newItems []T) {
	q.mu.Lock()
	q.items = append(q.items, newItems...)
	q.mu.Unlock()

	go func() {
		select {
		case q.writers <- struct{}{}:
			for {
				q.mu.Lock()

				if len(q.items) == 0 {
					q.mu.Unlock()
					break
				}

				item := q.items[0]

				q.mu.Unlock()

				q.input <- item

				q.mu.Lock()
				q.items = q.items[1:]
				q.mu.Unlock()
			}
			<-q.writers
		default:
			return
		}
	}()
}

func (q *Queue[T]) GetJobs() []T {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.items
}
