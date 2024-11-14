package bot

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal"
)

type TgBot struct {
	server  *http.Server
	errChan chan error
}

func NewTgBot() internal.Bot {
	errChan := make(chan error, 1)
	s := &http.Server{
		ReadHeaderTimeout: 1 * time.Second,
	}
	return &TgBot{
		server:  s,
		errChan: errChan,
	}
}

func (b *TgBot) Start(_ context.Context, _ *internal.Data) error {
	go func() {
		if err := b.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			b.errChan <- err
		}
	}()
	return nil
}

func (b *TgBot) ErrChan() chan error {
	return b.errChan
}

func (b *TgBot) Shutdown(ctx context.Context) error {
	return b.server.Shutdown(ctx)
}
