package internal

import "time"

const (
	DateLayout        = "02.01.2006"
	TimeOfStartLayout = "15:04"
)

type Subject struct {
	Discipline  *string
	Institute   string
	Department  *string
	Year        *int
	Group       string
	Professor   *string
	Date        *time.Time
	TimeOfStart *time.Time
	Classroom   *string
	Comment     *string
}
