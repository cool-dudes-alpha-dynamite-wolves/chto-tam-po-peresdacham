package parser

import "github.com/cool-dudes-alpha-dynamite-wolves/chto-tam-po-peresdacham/internal"

type ExcelParser struct{}

func NewExcelParser() internal.Parser {
	return &ExcelParser{}
}

func (b *ExcelParser) Parse() (*internal.Data, error) {
	return &internal.Data{}, nil
}
