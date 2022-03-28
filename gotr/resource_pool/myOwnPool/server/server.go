package main

import "context"

type (
	// Server can be started and stopped
	Server interface {
		Start() error
		Stop()
		Shutdown(ctx context.Context) error
	}
)
