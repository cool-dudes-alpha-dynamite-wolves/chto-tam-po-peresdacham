package parser

import (
	"time"
)

const (
	requiredFieldsInARowMinimum = 3

	dateLayout        = "02.01.2006"
	timeOfStartLayout = "15:04"
)

type subject struct {
	department  *string
	institute   *string
	discipline  *string
	year        *int
	group       *string
	professor   *string
	date        *time.Time
	timeOfStart *time.Time
	classroom   *string
	comment     *string
}

var (
	validDepartmentFields = map[string]struct{}{
		"Кафедра": {},
	}
	validInstituteFields = map[string]struct{}{
		"Институт": {},
	}
	validDisciplineFields = map[string]struct{}{
		"Дисциплина": {},
	}
	validGroupFields = map[string]struct{}{
		"Группа": {},
	}
	validProfessorFields = map[string]struct{}{
		"Преподаватель": {},
	}
	validYearFields = map[string]struct{}{
		"Курс": {},
	}
	validClassroomFields = map[string]struct{}{
		"Аудитория": {},
	}
	validTimeOfStartFields = map[string]struct{}{
		"Время": {},
	}
	validDateFields = map[string]struct{}{
		"Дата": {},
	}
	validCommentFields = map[string]struct{}{
		"Примечание": {},
	}
)
