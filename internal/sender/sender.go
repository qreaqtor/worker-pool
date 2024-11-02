package msgsender

import (
	"fmt"
	"log/slog"

	"github.com/qreator/worker-pool/internal/models"
)

type Sender[In, Out any] struct {
	output <-chan models.OutMsg[In, Out]
}

func NewSender[In, Out any](out <-chan models.OutMsg[In, Out]) *Sender[In, Out] {
	return &Sender[In, Out]{
		output: out,
	}
}

func (s *Sender[In, Out]) Run() {
	for msg := range s.output {
		if msg.Err != nil {
			slog.Error(msg.Err.Error(), slog.Int("worker", msg.Id))
			continue
		}

		slog.Info(fmt.Sprintf("worker=%d", msg.Id), slog.Any("data", msg.IncomigData), slog.Any("result", msg.Result))
	}
}
