package worker

import (
	"time"
)

type SleepWorker struct {
	dur time.Duration
}

func NewSleepWorker(d time.Duration) *SleepWorker {
	return &SleepWorker{
		dur: d,
	}
}

func (sw *SleepWorker) SleepWorkerFunc() any {
	return WorkerFunc[string](func(msg string) string {

		time.Sleep(sw.dur)
		return msg
	})
}
