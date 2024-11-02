package appserver

import (
	"context"
	"net"
	"sync/atomic"
)

type server interface {
	Close() error
	Serve(net.Listener) error
}

type AppServer struct {
	started atomic.Bool

	ctx context.Context

	server server

	addr string

	errChan chan error
}

// addr is a network address that must match the form "host:port"
func NewAppServer(ctx context.Context, server server, addr string) *AppServer {
	return &AppServer{
		ctx:     ctx,
		addr:    addr,
		server:  server,
		errChan: make(chan error, 1),
	}
}

// start listen
func (a *AppServer) Start() error {
	if a.started.Swap(true) {
		return ErrAlreadyStarted
	}

	l, err := net.Listen("tcp", a.addr)
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		err := a.server.Serve(l)
		if err != nil {
			errChan <- err
		}
	}()

	go func(ctx context.Context) {
		defer func() {
			err := l.Close()
			if err != nil {
				a.errChan <- err
			}

			err = a.server.Close()
			if err != nil {
				a.errChan <- err
			}

			close(a.errChan)
		}()

		select {
		case <-ctx.Done():
			return
		case err, ok := <-errChan:
			if ok {
				a.errChan <- err
			}
			return
		}
	}(a.ctx)

	return nil
}

// waiting when all goroutines is done and return serve errors
func (a *AppServer) Wait() []error {
	errs := make([]error, 0)

	for err := range a.errChan {
		errs = append(errs, err)
	}

	return errs
}
