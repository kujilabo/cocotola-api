package domain

import "github.com/go-playground/validator"

type Recordbook interface {
	GetID() WorkbookID
	GetResults() map[ProblemID]int
}

type recordbook struct {
	ID      WorkbookID `validate:"required"`
	Results map[ProblemID]int
}

func NewRecordbook(id WorkbookID, results map[ProblemID]int) (Recordbook, error) {
	m := &recordbook{
		ID:      id,
		Results: results,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (m *recordbook) GetID() WorkbookID {
	return m.ID
}

func (m *recordbook) GetResults() map[ProblemID]int {
	return m.Results
}
