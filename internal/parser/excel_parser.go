package parser

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal"
	"github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/pkg"
)

type ExcelParser struct {
	pathToRetakes string
}

type subjectOption func(subject *subject, value string) error

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
		if filePath != "retakes/matematika_3-kurs.xlsx" {
			continue
		}
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
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

				opts := getSubjectOpts(rows[rowWithHeaders])

				for i, row := range rows {
					if i <= rowWithHeaders || len(row) != len(rows[rowWithHeaders]) {
						continue
					}
					subj := &subject{}
					isSuccess := true
					for idx, cell := range row {
						if err = opts[idx](subj, cell); err != nil {
							isSuccess = false
							break
						}
					}
					if !isSuccess {
						log.Printf("got bad subject data: %s", err)
						continue
					}
					subjChan <- subj
				}
			}
		}(filePath)
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(subjChan)
	}(wg)

	result := make([]*internal.Subject, 0)
	for subj := range subjChan {
		fmt.Printf("%+#v\n", subj)
		result = append(result, subj.toDomain())
	}

	return result, nil
}

func getSubjectOpts(headers []string) (opts []subjectOption) {
	for _, cell := range headers {
		if _, ok := validDepartmentFields[cell]; ok {
			opts = append(opts, func(s *subject, value string) error {
				s.department = pkg.Pointer(value)
				return nil
			})
		} else if _, ok = validInstituteFields[cell]; ok {
			opts = append(opts, func(s *subject, value string) error {
				s.institute = pkg.Pointer(value)
				return nil
			})
		} else if _, ok = validDisciplineFields[cell]; ok {
			opts = append(opts, func(s *subject, value string) error {
				s.discipline = pkg.Pointer(value)
				return nil
			})
		} else if _, ok = validGroupFields[cell]; ok {
			opts = append(opts, func(s *subject, value string) error {

				s.group = pkg.Pointer(value)
				return nil
			})
		} else if _, ok = validProfessorFields[cell]; ok {
			opts = append(opts, func(s *subject, value string) error {
				s.professor = pkg.Pointer(value)
				return nil
			})
		} else if _, ok = validClassroomFields[cell]; ok {
			opts = append(opts, func(s *subject, value string) error {
				s.classroom = pkg.Pointer(value)
				return nil
			})
		} else if _, ok = validYearFields[cell]; ok {
			opts = append(opts, func(s *subject, value string) error {
				if v, err := strconv.Atoi(value); err == nil {
					s.year = pkg.Pointer(v)
					return nil
				}
				return fmt.Errorf("can not parse year %s", value)
			})
		} else if _, ok = validDateFields[cell]; ok {
			opts = append(opts, func(s *subject, value string) error {
				if v, err := time.Parse(dateLayout, value); err == nil {
					s.date = pkg.Pointer(v)
					return nil
				}
				return fmt.Errorf("can not parse date %s", value)
			})
		} else if _, ok = validTimeOfStartFields[cell]; ok {
			opts = append(opts, func(s *subject, value string) error {
				if v, err := time.Parse(timeOfStartLayout, value); err == nil {
					s.timeOfStart = pkg.Pointer(v)
					return nil
				}
				return fmt.Errorf("can not parse time of start %s", value)
			})
		} else if _, ok = validCommentFields[cell]; ok {
			opts = append(opts, func(s *subject, value string) error {
				s.comment = pkg.Pointer(value)
				return nil
			})
		} else {
			log.Printf("we do not know field type %s\n", cell)
			opts = append(opts, func(_ *subject, _ string) error { return nil })
		}
	}
	return
}
