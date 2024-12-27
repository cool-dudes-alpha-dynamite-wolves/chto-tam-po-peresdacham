package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal"
)

type TgBot struct {
	server   *http.Server
	errChan  chan error
	bot      *telegram.BotAPI
	subjects []*internal.Subject // Данные из парсера
}

func NewTgBot(token, port string) (*TgBot, error) {
	errChan := make(chan error, 1)
	bot, err := telegram.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Telegram Bot API: %v", err)
	}

	// bot.Debug = true

	return &TgBot{
		server: &http.Server{
			ReadHeaderTimeout: 1 * time.Second,
			Addr:              ":" + port,
		},
		errChan: errChan,
		bot:     bot,
	}, nil
}

func (b *TgBot) Start(_ context.Context, subjects []*internal.Subject) error {
	b.subjects = subjects
	updateConfig := telegram.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := b.bot.GetUpdatesChan(updateConfig)

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			switch update.Message.Command() {
			case "start":
				b.sendMessage(update.Message.Chat.ID, "Добро пожаловать! Введите /help для списка команд.")
			case "help":
				b.sendMessage(update.Message.Chat.ID, "Список команд: /start, /help, /schedule")
			case "schedule":
				args := strings.Fields(update.Message.Text)
				if len(args) < 3 {
					b.sendMessage(update.Message.Chat.ID, "Введите команду в формате: /schedule <Институт> <Группа>")
					return
				}
				instituteFilter := args[1]
				groupFilter := args[2]
				b.handleSchedule(update.Message.Chat.ID, instituteFilter, groupFilter)
			default:
				b.sendMessage(update.Message.Chat.ID, "Неизвестная команда. Введите /help.")
			}
		}
	}()

	go func() {
		if err := b.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			b.errChan <- err
		}
	}()
	return nil
}

func (b *TgBot) handleSchedule(chatID int64, instituteFilter, groupFilter string) {
	// TODO: ИТКН и ИКН - один институт, был ребрендинг ;-)
	if instituteFilter == "ИТКН" {
		instituteFilter = "ИКН"
	}

	message := "Расписание пересдач:\n"
	for _, subj := range b.subjects {
		if !(subj.Institute == instituteFilter && subj.Group == groupFilter) {
			continue
		}
		message += b.constructSubjectMsg(subj)
	}
	if message == "Расписание пересдач:\n" {
		message = "Нет данных для выбранных фильтров."
	}
	b.sendMessage(chatID, message)
}

func (b *TgBot) sendMessage(chatID int64, text string) {
	msg := telegram.NewMessage(chatID, text)
	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v\n", err)
	}
}

func (b *TgBot) ErrChan() chan error {
	return b.errChan
}

func (b *TgBot) Shutdown(ctx context.Context) error {
	return b.server.Shutdown(ctx)
}
