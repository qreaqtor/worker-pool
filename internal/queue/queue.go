package queue

import (
	"context"
	"sync"
)

type Queue[T any] struct {
	input   chan<- T
	writers chan struct{}

	mu *sync.RWMutex

	items []T

	ctx context.Context
}

func NewQueue[T any](ctx context.Context,input chan<- T) *Queue[T] {
	return &Queue[T]{
		items:   make([]T, 0),
		input:   input,
		mu:      &sync.RWMutex{},
		writers: make(chan struct{}, 1),
		ctx: ctx,
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
			defer func() {
				<-q.writers
			}()

			for {
				q.mu.Lock()

				if len(q.items) == 0 {
					q.mu.Unlock()
					return
				}

				item := q.items[0] // беру следующий элемент, который нужно отправить

				q.mu.Unlock()

				select {
				case q.input <- item:
					q.mu.Lock()
					q.items = q.items[1:] // выполняю смещение только здесь, чтобы в GetJobs(), помимо всех предстоящих был и следующий элемент, который будет отправлен
					q.mu.Unlock()
				case <-q.ctx.Done(): // если контекст закрыт, то выхожу без записи в канал
					return
				}
			}
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
