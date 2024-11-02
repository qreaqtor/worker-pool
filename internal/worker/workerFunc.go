package worker

type WorkerFunc[Out any] func(string) Out
