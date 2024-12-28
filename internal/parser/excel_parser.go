package parser

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal"
	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/pkg"
)

type ExcelParser struct {
	pathToRetakes string
}

type rawSubjectOption func(rawSubject *rawSubject, value string) error

func NewExcelParser(pathToRetakes string) (*ExcelParser, error) {
	fileInfo, err := os.Stat(pathToRetakes)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("provided path to retakes is not directory")
	}
	return &ExcelParser{
		pathToRetakes: pathToRetakes,
	}, nil
}

func (p *ExcelParser) Parse() ([]*internal.Subject, error) {
	var files []string
	err := filepath.WalkDir(p.pathToRetakes, func(path string, d fs.DirEntry, _ error) error {
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	subjChan := make(chan *subject)
	wg := &sync.WaitGroup{}
	for _, filePath := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			p.processFile(filePath, subjChan)
		}(filePath)
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(subjChan)
	}(wg)

	result := make([]*internal.Subject, 0)
	for subj := range subjChan {
		result = append(result, subj.toDomain())
	}

	return result, nil
}

func (p *ExcelParser) processFile(filePath string, subjChan chan<- *subject) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Printf("failed to open .xlsx file %s\n", filePath)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close .xlsx file %s\n", err)
		}
	}()

	for _, sheet := range f.GetSheetList() {
		rows, err := f.GetRows(sheet)
		if err != nil {
			log.Printf("could not open .xlsx file %s\n", filePath)
			return
		}
		// at least there should be 2 rows - headers and 1 line of data
		if len(rows) < 2 {
			log.Printf("there is no data on sheet %s\n", sheet)
			continue
		}

		rowWithHeaders := -1
		for i, row := range rows {
			if len(row) >= requiredFieldsInARowMinimum {
				rowWithHeaders = i
				break
			}
		}
		if rowWithHeaders == -1 {
			log.Printf("there is no data on sheet %s\n", sheet)
			continue
		}

		opts := p.getSubjectOpts(rows[rowWithHeaders])

		for i, row := range rows {
			if i <= rowWithHeaders || len(row) != len(rows[rowWithHeaders]) {
				continue
			}
			rSubj := &rawSubject{}
			isSuccess := true
			for idx, cell := range row {
				if err = opts[idx](rSubj, cell); err != nil {
					isSuccess = false
					break
				}
			}
			if !isSuccess {
				log.Printf("got bad raw subject data: %s", err)
				continue
			}

			for _, subj := range rSubj.extend() {
				if err = subj.validate(); err != nil {
					log.Printf("got invalid subject data: %s", err)
					continue
				}
				subjChan <- subj
			}
		}
	}
}

func (p *ExcelParser) getSubjectOpts(headers []string) (opts []rawSubjectOption) {
	for _, cell := range headers {
		if _, ok := validDepartmentFields[cell]; ok {
			opts = append(opts, func(s *rawSubject, value string) error {
				s.department = pkg.Pointer(value)
				return nil
			})
		} else if _, ok = validInstituteFields[cell]; ok {
			opts = append(opts, func(s *rawSubject, value string) error {
				if v := institute(value); v.isValid() {
					s.institute = pkg.Pointer(v)
					return nil
				}
				return fmt.Errorf("got invalid institute data %s", value)
			})
		} else if _, ok = validDisciplineFields[cell]; ok {
			opts = append(opts, func(s *rawSubject, value string) error {
				s.discipline = pkg.Pointer(value)
				return nil
			})
		} else if _, ok = validGroupFields[cell]; ok {
			opts = append(opts, func(s *rawSubject, value string) error {
				if v := group(value); v.isValid() {
					s.groups = []*group{&v}
					return nil
				}
				if groups := strings.Split(value, ","); len(groups) > 1 {
					for idx, rawGroup := range groups {
						groups[idx] = strings.Trim(rawGroup, "\r\n")
					}
					s.groups = pkg.Map(groups, func(g string) *group {
						return pkg.Pointer(group(g))
					})
					return nil
				}
				return fmt.Errorf("got invalid group data %s", value)
			})
		} else if _, ok = validProfessorFields[cell]; ok {
			opts = append(opts, func(s *rawSubject, value string) error {
				s.professor = pkg.Pointer(value)
				return nil
			})
		} else if _, ok = validClassroomFields[cell]; ok {
			opts = append(opts, func(s *rawSubject, value string) error {
				s.classroom = pkg.Pointer(value)
				return nil
			})
		} else if _, ok = validYearFields[cell]; ok {
			opts = append(opts, func(s *rawSubject, value string) error {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("can not parse year %s; err: %s", value, err)
				}
				s.year = pkg.Pointer(v)
				return nil
			})
		} else if _, ok = validDateFields[cell]; ok {
			opts = append(opts, func(s *rawSubject, value string) error {
				v, err := time.Parse(internal.DateLayout, value)
				if err != nil {
					return fmt.Errorf("can not parse date %s; err: %s", value, err)
				}
				s.date = pkg.Pointer(v)
				return nil
			})
		} else if _, ok = validTimeOfStartFields[cell]; ok {
			opts = append(opts, func(s *rawSubject, value string) error {
				v, err := time.Parse(internal.TimeOfStartLayout, value)
				if err != nil {
					return fmt.Errorf("can not parse time of start %s; err: %s", value, err)
				}
				s.timeOfStart = pkg.Pointer(v)
				return nil
			})
		} else if _, ok = validCommentFields[cell]; ok {
			opts = append(opts, func(s *rawSubject, value string) error {
				s.comment = pkg.Pointer(value)
				return nil
			})
		} else {
			log.Printf("we do not know field type %s\n", cell)
			opts = append(opts, func(_ *rawSubject, _ string) error { return nil })
		}
	}
	return
}
