package dummyserver

import (
	"context"
)

type DummyServer struct {
	ctx context.Context
}

func NewDummyServer(ctx context.Context) *DummyServer {
	return &DummyServer{ctx}
}

func (d *DummyServer) Start() error {
	return nil
}

func (d *DummyServer) Wait() []error {
	<-d.ctx.Done()
	return nil
}
