package msgsender

import (
	"fmt"
	"log/slog"

	"github.com/qreator/worker-pool/internal/models"
)

type Sender[T any] struct {
	output <-chan models.OutMsg[T]
}

func NewSender[T any](out <-chan models.OutMsg[T]) *Sender[T] {
	return &Sender[T]{
		output: out,
	}
}

func (s *Sender[T]) Run() {
	for msg := range s.output {
		if msg.Err != nil {
			slog.Error(msg.Err.Error(), slog.Int("worker", msg.Id))
			continue
		}

		slog.Info(fmt.Sprintf("worker=%d", msg.Id), slog.String("data", msg.Data), slog.Any("result", msg.Result))
	}
}
