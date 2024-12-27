package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/multierr"

	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal/bot"
	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal/parser"
)

var (
	// хранить секреты в коде плохо, поэтому используем флаги
	// "1953480583:AAEU7eBaZnCUt525oUkCMRCQxK1TJmaoVd4"
	botToken = flag.String("tg.token", "", "token for telegram")
	// запуск выглядит так:
	// go run bot.go -tg.token="1953480583:AAEU7eBaZnCUt525oUkCMRCQxK1TJmaoVd4"
)

func main() {
	flag.Parse()

	ctx := context.Background()

	// Инициализация парсера
	parser, err := parser.NewExcelParser("./retakes")
	if err != nil {
		log.Fatal("failed to initialize parser", err)
	}

	// Парсинг данных
	data, err := parser.Parse()
	if err != nil {
		log.Fatal("failed to parse input data", err.Error())
	}
	log.Println("Parsing completed!")

	// Инициализация и запуск бота
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Println("PORT:", port)
	if *botToken == "" {
		*botToken = os.Getenv("BOT_TOKEN")
	}
	bot, err := bot.NewTgBot(*botToken, port)
	if err != nil {
		log.Fatal("failed to initialize bot", err)
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	err = run(ctx, stop,
		func(ctx context.Context) error {
			return bot.Start(ctx, data)
		},
		func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				log.Println("got stop signal")
				return nil
			case err := <-bot.ErrChan():
				log.Println("got err from bot:", err.Error())
				return err
			}
		},
	)

	log.Println("application is shutting down...")
	if err != nil {
		log.Println("failed to shutdown application", err)
	}

	if err = bot.Shutdown(ctx); err != nil {
		log.Println("got error during bot shutdown:", err.Error())
	}
}

func run(ctx context.Context, shutdown func(), startUpFuncs ...func(ctx context.Context) error) error {
	errCh := make(chan error)
	for i := range startUpFuncs {
		go func(i int) { errCh <- startUpFuncs[i](ctx) }(i)
	}

	var err error
	var closed bool
	for range startUpFuncs {
		err = multierr.Append(err, <-errCh)
		if err != nil && !closed {
			shutdown()
			closed = true
		}
	}

	return err
}
