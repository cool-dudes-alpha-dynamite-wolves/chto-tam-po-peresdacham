package internal

import "context"

type Bot interface {
	Start(context.Context, []*Subject) error
	Shutdown(ctx context.Context) error
	ErrChan() chan error
}

type Parser interface {
	Parse() ([]*Subject, error)
}
