package parser

import (
	"regexp"
	"time"
)

const (
	requiredFieldsInARowMinimum = 3

	dateLayout        = "02.01.2006"
	timeOfStartLayout = "15:04"

	disciplinePatternRegex = "^[А-Я]{3,4}-\\d{2}-\\d$"
)

type institute string

const (
	itknInstitute   institute = "ИТКН"
	iknInstitute    institute = "ИКН" // equal to "ИТКН"
	ekotehInstitute institute = "ЭкоТех"
	inminInstitute  institute = "ИНМиН"
	euppInstitute   institute = "ЭУПП"
	iboInstitute    institute = "ИБО"
	inobrInstitute  institute = "ИНОБР"
	giInstitute     institute = "ГИ"
)

func (i institute) isValid() bool {
	_, ok := institute2groupMapping[i]
	return ok
}

type group string

func (d group) isValid() bool {
	if ok, err := regexp.MatchString(disciplinePatternRegex, string(d)); err == nil && ok {
		return true
	}
	return false
}

// rawSubject represents subject, that we get directly
// from the table
type rawSubject struct {
	department  *string
	institute   *institute
	discipline  *string
	year        *int
	groups      []*group
	professor   *string
	date        *time.Time
	timeOfStart *time.Time
	classroom   *string
	comment     *string
}

type subject struct {
	department  *string
	institute   *institute
	discipline  *string
	year        *int
	group       *group
	professor   *string
	date        *time.Time
	timeOfStart *time.Time
	classroom   *string
	comment     *string
}

var (
	institute2groupMapping = map[institute][]string{
		itknInstitute: {
			"БИВТ",
			"БПМ",
			"ББИ",
		},
		iknInstitute: {
			"БИВТ",
			"БПМ",
			"ББИ",
		},
		ekotehInstitute: {
			"БМТ",
			"БТМО",
		},
		inminInstitute: {
			"БМТМ",
			"БФЗ",
			"БНТМ",
			"БЭН",
		},
		euppInstitute: {
			"БЭК",
			"БТД",
		},
		iboInstitute: {
			"БЛГ",
		},
		giInstitute: {
			"СГД",
			"БЭЭ",
		},
		inobrInstitute: {},
	}

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
