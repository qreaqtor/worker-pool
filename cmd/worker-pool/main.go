package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	workersapp "github.com/qreator/worker-pool/internal/app"
	"github.com/qreator/worker-pool/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// уведомляю канал, а не контекст, чтобы отменять контекст в том же месте где и создал
	closeCtx := make(chan os.Signal, 1)
	signal.Notify(closeCtx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-closeCtx
		cancel()
	}()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	app := workersapp.NewApp(ctx, cfg, closeCtx)

	err = app.Start()
	if err != nil {
		log.Fatalln(err)
	}

	errs := app.Wait()
	if len(errs) != 0 {
		log.Fatalln(errs)
	}
}
