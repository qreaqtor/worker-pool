package app

import (
	"context"
	"os"

	"github.com/gorilla/mux"
	apiconsole "github.com/qreator/worker-pool/internal/api/console"
	apiweb "github.com/qreator/worker-pool/internal/api/web"
	"github.com/qreator/worker-pool/internal/config"
	"github.com/qreator/worker-pool/internal/models"
	msgsender "github.com/qreator/worker-pool/internal/sender"
	"github.com/qreator/worker-pool/internal/worker"
	workerpool "github.com/qreator/worker-pool/internal/worker-pool"
	appserver "github.com/qreator/worker-pool/pkg/appServer"
	dummyserver "github.com/qreator/worker-pool/pkg/dummyServer"
	httpserver "github.com/qreator/worker-pool/pkg/httpServer"
)

type sender interface {
	Run()
}

type server interface {
	Start() error
	Wait() []error
}

type App struct {
	server server

	sender sender
}

func NewApp(ctx context.Context, cfg *config.Config, closeCtx chan<- os.Signal) *App {
	output := make(chan models.OutMsg[string], 1)

	createWorker := worker.NewSleepWorker(cfg.Workers.Sleeper.Sleep).SleepWorkerFunc

	workerSrv := workerpool.NewWorkerPoolSrv[string](ctx, output, createWorker)

	app := &App{
		sender: msgsender.NewSender(output),
	}

	if cfg.AppType == "web" {
		api := apiweb.NewWebWorkersAPI(workerSrv)
		r := mux.NewRouter()
		api.Register(r)
		app.server = appserver.NewAppServer(ctx, httpserver.NewHTTPServer(r), cfg.Addr.ToString())

	} else {
		api := apiconsole.NewConsoleWorkersAPI(workerSrv)
		api.Register(os.Stdin, closeCtx)
		app.server = dummyserver.NewDummyServer(ctx)
	}

	initLogger()

	return app
}

func (a *App) Start() error {
	go a.sender.Run()

	return a.server.Start()
}

func (a *App) Wait() []error {
	return a.server.Wait()
}
