package bot

import (
	"fmt"

	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal"
)

func (b *TgBot) constructSubjectMsg(subject *internal.Subject) (message string) {
	message += fmt.Sprintf("- Институт: %s\n", subject.Institute)
	message += fmt.Sprintf("  Группа: %s\n", subject.Group)
	if subject.Discipline != nil {
		message += fmt.Sprintf("  Дисциплина: %s\n", *subject.Discipline)
	}
	if subject.Year != nil {
		message += fmt.Sprintf("  Курс: %d\n", *subject.Year)
	}
	if subject.Professor != nil {
		message += fmt.Sprintf("  Преподаватель: %s\n", *subject.Professor)
	}
	if subject.Date != nil {
		message += fmt.Sprintf("  Дата: %s\n", subject.Date.Format(internal.DateLayout))
	}
	if subject.TimeOfStart != nil {
		message += fmt.Sprintf("  Время: %s\n", subject.TimeOfStart.Format(internal.TimeOfStartLayout))
	}
	if subject.Classroom != nil {
		message += fmt.Sprintf("  Аудитория: %s\n", *subject.Classroom)
	}
	if subject.Comment != nil {
		message += fmt.Sprintf("  Примечание: %s\n", *subject.Comment)
	}
	message += "\n"
	return
}
