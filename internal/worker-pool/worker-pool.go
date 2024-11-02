package workerpool

import (
	"context"
	"sync"

	"github.com/qreator/worker-pool/internal/models"
	"github.com/qreator/worker-pool/internal/queue"
	"github.com/qreator/worker-pool/internal/worker"
)

type WorkerPool[In, Out any] struct {
	input  chan In
	output chan<- models.OutMsg[In, Out]

	next int

	mu *sync.RWMutex

	ctx context.Context

	workerIDs map[int]context.CancelFunc

	pool *sync.Pool

	queue *queue.Queue[In]
}

type WorkerPoolParams[In, Out any] struct {
	Ctx          context.Context
	CreateWorker func() any
	Output chan<- models.OutMsg[In, Out]
}

func NewWorkerPoolSrv[In, Out any](params WorkerPoolParams[In, Out]) *WorkerPool[In, Out] {
	input := make(chan In)

	pool := &sync.Pool{
		New: params.CreateWorker,
	}

	queue := queue.NewQueue(input)

	return &WorkerPool[In, Out]{
		mu:              &sync.RWMutex{},
		workerIDs:       make(map[int]context.CancelFunc),
		ctx:             params.Ctx,
		input:           input,
		pool:            pool,
		queue:           queue,
		output:          params.Output,
	}
}

func (w *WorkerPool[In, Out]) Delete(ids []int) {
	for _, id := range ids {
		w.deleteOne(id)
	}
}

func (w *WorkerPool[In, Out]) deleteOne(id int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if cancel, ok := w.workerIDs[id]; ok {
		cancel()
		delete(w.workerIDs, id)
	}
}

func (w *WorkerPool[In, Out]) Add(n int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for range n {
		ctx, cancel := context.WithCancel(w.ctx)

		w.next++

		go w.startWorker(ctx, w.next)

		w.workerIDs[w.next] = cancel
	}
}

// Return slice with id alive's workers.
func (w *WorkerPool[In, Out]) Alive() []int {
	w.mu.RLock()
	defer w.mu.RUnlock()

	workers := make([]int, 0, len(w.workerIDs))

	for id := range w.workerIDs {
		workers = append(workers, id)
	}

	return workers
}

func (w *WorkerPool[In, Out]) startWorker(ctx context.Context, id int) {
	defer w.deleteOne(id)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-w.input:
			if !ok {
				return
			}

			outMsg := models.OutMsg[In, Out]{
				Id:   id,
				Data: msg,
			}

			worker, ok := w.pool.Get().(worker.WorkerFunc[In, Out])
			if !ok {
				outMsg.Err = errBadWorkerFuncType
			} else {
				outMsg.Result = worker(msg)
				w.pool.Put(worker)
			}

			w.output <- outMsg
		}
	}
}

func (w *WorkerPool[In, Out]) Work(jobs []In) {
	w.queue.Append(jobs)
}

func (w *WorkerPool[In, Out]) GetJobs() []In {
	return w.queue.GetJobs()
}
