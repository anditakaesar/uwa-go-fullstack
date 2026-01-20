package common

import "fmt"

type Pagination struct {
	Page int
	Size int
}

func (p *Pagination) GetOffset() int {
	offset := (p.Page - 1) * p.Size
	return offset
}

type SortDirection string

const (
	SORT_ASC  SortDirection = "ASC"
	SORT_DESC SortDirection = "DESC"
)

type Sort struct {
	Field     string
	Direction SortDirection
}

func (s *Sort) ToSQLSort() string {
	return fmt.Sprintf("%s %s", s.Field, s.Direction)
}
