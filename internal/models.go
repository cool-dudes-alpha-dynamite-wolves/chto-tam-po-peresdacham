package internal

import "time"

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
