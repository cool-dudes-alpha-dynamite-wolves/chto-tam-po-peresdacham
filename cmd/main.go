package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal/bot"
	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal/parser"
)

var (

	// хранить секреты в коде плохо, поэтому используем флаги
	// "1953480583:AAEU7eBaZnCUt525oUkCMRCQxK1TJmaoVd4"
	botToken = flag.String("tg.token", "", "token for telegram")

	// это не секрет, но для простоты тоже выносим под конфиг
	// "https://5872-95-165-1-28.eu.ngrok.io"
	WebhookURL = flag.String("tg.webhook", "", "webhook addr for telegram")

	// запуск выглядит так:
	// go run bot.go -tg.token="1953480583:AAEU7eBaZnCUt525oUkCMRCQxK1TJmaoVd4" -tg.webhook="https://5872-95-165-1-28.eu.ngrok.io"
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
	bot, err := bot.NewTgBot(*botToken)
	if err != nil {
		log.Fatal("failed to initialize bot", err)
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
