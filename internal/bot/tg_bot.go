package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	server   *http.Server
	errChan  chan error
	bot      *telegram.BotAPI
	subjects []*internal.Subject // Данные из парсера
}
type Subject struct {
	Name string    // Название предмета
	Date time.Time // Дата пересдачи
}

func NewTgBot() *TgBot {
	errChan := make(chan error, 1)
	botAPI, err := telegram.NewBotAPI(token)
	if err != nil {
		log.Fatalf("failed to initialize Telegram Bot API: %v", err)
	}

	return &TgBot{
		server:  &http.Server{ReadHeaderTimeout: 1 * time.Second},
		errChan: errChan,
		bot:     botAPI,
	}
}

func (b *TgBot) Start(ctx context.Context, subjects []*internal.Subject) error {
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
	if len(b.subjects) == 0 {
		b.sendMessage(chatID, "Нет данных о пересдачах.")
		return
	}

	message := "Расписание пересдач:\n"
	for _, subj := range b.subjects {
		if subj.Institute != instituteFilter {
			continue
		}
		if !strings.HasPrefix(subj.Group, groupFilter) {
			continue
		}

		date := "не указана"
		timeOfStart := "не указано"

		if subj.Date != nil {
			date = subj.Date.Format("02.01.2006")
		}
		if subj.TimeOfStart != nil {
			timeOfStart = subj.TimeOfStart.Format("15:04")
		}

		message += fmt.Sprintf(
			"- Дисциплина: %s\n  Институт: %s\n  Группа: %s\n  Дата: %s\n  Время: %s\n  Аудитория: %s\n\n",
			getOrDefault(subj.Discipline, "не указана"),
			subj.Institute,
			subj.Group,
			date,
			timeOfStart,
			getOrDefault(subj.Classroom, "не указана"),
		)
	}
	if message == "Расписание пересдач:\n" {
		message = "Нет данных для выбранных фильтров."
	}
	b.sendMessage(chatID, message)
}

func getOrDefault(value *string, defaultValue string) string {
	if value != nil {
		return *value
	}
	return defaultValue
}

func (b *TgBot) sendMessage(chatID int64, text string) {
	msg := telegram.NewMessage(chatID, text)
	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

func (b *TgBot) ErrChan() chan error {
	return b.errChan
}

func (b *TgBot) Shutdown(ctx context.Context) error {
	return b.server.Shutdown(ctx)
}
