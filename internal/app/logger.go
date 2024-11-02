package app

import (
	"log/slog"
	"os"

	"github.com/qreator/worker-pool/pkg/logging/pretty"
)

func initLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	slog.SetDefault(slog.New(pretty.NewPrettyHandler(os.Stdout, opts)))
}
