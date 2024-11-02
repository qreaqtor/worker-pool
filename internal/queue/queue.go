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
		// если есть активные writers, то просто выхожу из функции
		select {
		case q.writers <- struct{}{}:
			for {
				q.mu.Lock()

				if len(q.items) == 0 {
					q.mu.Unlock()
					break
				}

				item := q.items[0] // беру следующий элемент, который нужно отправить

				q.mu.Unlock()

				q.input <- item

				q.mu.Lock()
				q.items = q.items[1:] // выполняю смещение только здесь, чтобы в GetJobs(), помимо всех предстоящих был и следующий элемент, который будет отправлен
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
