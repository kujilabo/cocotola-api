package domain

import (
	"context"

	"github.com/go-playground/validator"
	"golang.org/x/xerrors"
)

var ErrProblemAlreadyExists = xerrors.New("problem already exists")
var ErrProblemNotFound = xerrors.New("problem not found")
var ErrProblemOtherError = xerrors.New("problem other error")

type ProblemAddParameter struct {
	WorkbookID  WorkbookID `validate:"required"`
	Number      int        `validate:"required"`
	ProblemType string     `validate:"required"`
	Properties  map[string]string
}

type ProblemUpdateParameter struct {
	Number     int `validate:"required"`
	Properties map[string]string
}

func NewProblemAddParameter(workbookID WorkbookID, number int, problemType string, properties map[string]string) (*ProblemAddParameter, error) {
	m := &ProblemAddParameter{
		WorkbookID:  workbookID,
		Number:      number,
		ProblemType: problemType,
		Properties:  properties,
	}

	v := validator.New()
	return m, v.Struct(m)
}

type ProblemSearchCondition struct {
	WorkbookID WorkbookID
	PageNo     int `validate:"required,gte=1"`
	PageSize   int `validate:"required,gte=1,lte=100"`
	Keyword    string
}

type ProblemSearchResult struct {
	TotalCount int64
	Results    []Problem
}

type ProblemIDsCondition struct {
	WorkbookID WorkbookID
	IDs        []ProblemID
}

type ProblemRepository interface {
	FindProblems(ctx context.Context, operator Student, param *ProblemSearchCondition) (*ProblemSearchResult, error)

	FindProblemsByProblemIDs(ctx context.Context, operator Student, param *ProblemIDsCondition) (*ProblemSearchResult, error)

	FindProblemIDs(ctx context.Context, operator Student, workbookID WorkbookID) ([]ProblemID, error)

	FindProblemByID(ctx context.Context, operator Student, workbookID WorkbookID, problemID ProblemID) (Problem, error)

	AddProblem(ctx context.Context, operator Student, param *ProblemAddParameter) (ProblemID, error)

	// UpdateProblem(ctx context.Context, operator Student, param *ProblemUpdateParameter) error

	RemoveProblem(ctx context.Context, operator Student, problemID ProblemID, version int) error
}
