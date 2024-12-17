package parser

import (
	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal"
)

func (s *subject) toDomain() *internal.Subject {
	return &internal.Subject{
		Discipline:  s.discipline,
		Institute:   s.institute,
		Department:  s.department,
		Year:        s.year,
		Group:       s.group,
		Professor:   s.professor,
		Date:        s.date,
		TimeOfStart: s.timeOfStart,
		Classroom:   s.classroom,
		Comment:     s.comment,
	}
}
