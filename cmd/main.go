package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal/bot"
	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal/parser"
)

func main() {
	ctx := context.Background()

	parser, err := parser.NewExcelParser("./retakes")
	if err != nil {
		log.Fatal("failed to initialize parser", err)
	}
	bot := bot.NewTgBot()

	data, err := parser.Parse()
	if err != nil {
		log.Fatal("failed to parse input data", err.Error())
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	err = func() error {
		if err = bot.Start(ctx, data); err != nil {
			log.Println("failed to start bot:", err.Error())
			return err
		}
		for {
			select {
			case <-ctx.Done():
				log.Println("got stop signal")
				return nil
			case err := <-bot.ErrChan():
				log.Println("got err from bot:", err.Error())
				return err
			}
		}
	}()

	log.Println("app is shutting down...")
	stop()

	if err = bot.Shutdown(ctx); err != nil {
		log.Println("got error during bot shutdown:", err.Error())
	}
}
