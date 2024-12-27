package parser

import (
	"fmt"
	"strings"

	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal"
)

func (s *rawSubject) extend() []*subject {
	result := make([]*subject, 0, len(s.groups))

	for _, g := range s.groups {
		subj := &subject{
			department:  s.department,
			institute:   s.institute,
			discipline:  s.discipline,
			year:        s.year,
			group:       g,
			professor:   s.professor,
			date:        s.date,
			timeOfStart: s.timeOfStart,
			classroom:   s.classroom,
			comment:     s.comment,
		}
		result = append(result, subj)
	}

	return result
}

func (s *subject) validate() error {
	// restriction: we consider only those subjects, that have "institute" and "group" != nil.
	if s.institute == nil || s.group == nil {
		return fmt.Errorf("subject with empty institute or group are not supported")
	}

	groups := institute2groupMapping[*s.institute]

	found := false
	for _, groupPrefix := range groups {
		if strings.HasPrefix(string(*s.group), groupPrefix) {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("subject has unknown group")
	}

	return nil
}

func (s *subject) toDomain() *internal.Subject {
	retSubj := &internal.Subject{
		Discipline:  s.discipline,
		Institute:   string(*s.institute),
		Department:  s.department,
		Year:        s.year,
		Group:       string(*s.group),
		Professor:   s.professor,
		Date:        s.date,
		TimeOfStart: s.timeOfStart,
		Classroom:   s.classroom,
		Comment:     s.comment,
	}

	// TODO: В домене, институт ИТКН = ИКН
	if *s.institute == itknInstitute {
		retSubj.Institute = string(iknInstitute)
	}
	return retSubj
}
