package workerpool

import (
	"context"
	"sync"

	"github.com/qreator/worker-pool/internal/models"
	"github.com/qreator/worker-pool/internal/worker"
)

type WorkerPool[Out any] struct {
	input  chan string
	output chan<- models.OutMsg[Out]

	next int

	mu *sync.RWMutex

	ctx context.Context

	workerIDs map[int]context.CancelFunc

	pool *sync.Pool
}

func NewWorkerPoolSrv[Out any](ctx context.Context, output chan<- models.OutMsg[Out], createWorker func() any) *WorkerPool[Out] {
	pool := &sync.Pool{
		New: createWorker,
	}

	return &WorkerPool[Out]{
		input:     make(chan string, 1),
		output:    output,
		next:      1,
		mu:        &sync.RWMutex{},
		ctx:       ctx,
		workerIDs: make(map[int]context.CancelFunc),
		pool:      pool,
	}
}

func (w *WorkerPool[Out]) Delete(ids []int) {
	for _, id := range ids {
		w.deleteOne(id)
	}
}

func (w *WorkerPool[Out]) deleteOne(id int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if cancel, ok := w.workerIDs[id]; ok {
		cancel()
		delete(w.workerIDs, id)
	}
}

func (w *WorkerPool[Out]) Add(n int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for range n {
		ctx, cancel := context.WithCancel(w.ctx)

		go w.startWorker(ctx, w.next)

		w.workerIDs[w.next] = cancel

		w.next++
	}
}

// Return slice with id alive's workers.
func (w *WorkerPool[Out]) Alive() []int {
	w.mu.RLock()
	defer w.mu.RUnlock()

	workers := make([]int, 0, len(w.workerIDs))

	for id := range w.workerIDs {
		workers = append(workers, id)
	}

	return workers
}

func (w *WorkerPool[Out]) startWorker(ctx context.Context, id int) {
	defer w.deleteOne(id)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-w.input:
			if !ok {
				return
			}

			outMsg := models.OutMsg[Out]{
				Id:   id,
				Data: msg,
			}

			worker, ok := w.pool.Get().(worker.WorkerFunc[Out])
			if !ok {
				outMsg.Err = errBadWorkerFuncType
			} else {
				outMsg.Result = worker(msg)
			}

			w.output <- outMsg

			if outMsg.Err == nil {
				w.pool.Put(worker)
			}
		}
	}
}

func (w *WorkerPool[Out]) Work(jobs []string) {
	for _, job := range jobs {
		go func(msg string) {
			w.input <- msg
		}(job)
	}
}
