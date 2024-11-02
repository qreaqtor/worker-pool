package worker

type WorkerFunc[In, Out any] func(In) Out
