package appserver

import "errors"

var (
	ErrAlreadyStarted = errors.New("server already started")
)
