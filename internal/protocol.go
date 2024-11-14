package internal

import "context"

type Bot interface {
	Start(context.Context, *Data) error
	Shutdown(ctx context.Context) error
	ErrChan() chan error
}

type Parser interface {
	Parse() (*Data, error)
}
