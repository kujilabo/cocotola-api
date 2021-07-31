package domain

import (
	"context"

	"github.com/go-playground/validator"
	"golang.org/x/xerrors"
)

var ErrProblemAlreadyExists = xerrors.New("problem already exists")
var ErrProblemNotFound = xerrors.New("problem not found")
var ErrProblemOtherError = xerrors.New("problem other error")

type ProblemAddParameter interface {
	GetWorkbookID() WorkbookID
	GetNumber() int
	GetProblemType() string
	GetProperties() map[string]string
}

type problemAddParameter struct {
	WorkbookID  WorkbookID `validate:"required"`
	Number      int        `validate:"required"`
	ProblemType string     `validate:"required"`
	Properties  map[string]string
}

func NewProblemAddParameter(workbookID WorkbookID, number int, problemType string, properties map[string]string) (ProblemAddParameter, error) {
	m := &problemAddParameter{
		WorkbookID:  workbookID,
		Number:      number,
		ProblemType: problemType,
		Properties:  properties,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (p *problemAddParameter) GetWorkbookID() WorkbookID {
	return p.WorkbookID
}
func (p *problemAddParameter) GetNumber() int {
	return p.Number
}
func (p *problemAddParameter) GetProblemType() string {
	return p.ProblemType
}
func (p *problemAddParameter) GetProperties() map[string]string {
	return p.Properties
}

type ProblemUpdateParameter interface {
	GetNumber() int
	GetProperties() map[string]string
}

type problemUpdateParameter struct {
	Number     int `validate:"required"`
	Properties map[string]string
}

func NewProblemUpdateParameter(workbookID WorkbookID, number int, problemType string, properties map[string]string) (ProblemUpdateParameter, error) {
	m := &problemUpdateParameter{
		Number:     number,
		Properties: properties,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (p *problemUpdateParameter) GetNumber() int {
	return p.Number
}
func (p *problemUpdateParameter) GetProperties() map[string]string {
	return p.Properties
}

type ProblemSearchCondition interface {
	GetWorkbookID() WorkbookID
	GetPageNo() int
	GetPageSize() int
	GetKeyword() string
}

type problemSearchCondition struct {
	WorkbookID WorkbookID
	PageNo     int `validate:"required,gte=1"`
	PageSize   int `validate:"required,gte=1,lte=100"`
	Keyword    string
}

func NewProblemSearchCondition(workbookID WorkbookID, pageNo, pageSize int, keyword string) (ProblemSearchCondition, error) {
	m := &problemSearchCondition{
		WorkbookID: workbookID,
		PageNo:     pageNo,
		PageSize:   pageSize,
		Keyword:    keyword,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (c *problemSearchCondition) GetWorkbookID() WorkbookID {
	return c.WorkbookID
}

func (c *problemSearchCondition) GetPageNo() int {
	return c.PageNo
}

func (c *problemSearchCondition) GetPageSize() int {
	return c.PageSize
}

func (c *problemSearchCondition) GetKeyword() string {
	return c.Keyword
}

type ProblemIDsCondition interface {
	GetWorkbookID() WorkbookID
	GetIDs() []ProblemID
}

type problemIDsCondition struct {
	WorkbookID WorkbookID
	IDs        []ProblemID
}

func NewProblemIDsCondition(workbookID WorkbookID, ids []ProblemID) (ProblemIDsCondition, error) {
	m := &problemIDsCondition{
		WorkbookID: workbookID,
		IDs:        ids,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (c *problemIDsCondition) GetWorkbookID() WorkbookID {
	return c.WorkbookID
}

func (c *problemIDsCondition) GetIDs() []ProblemID {
	return c.IDs
}

type ProblemSearchResult struct {
	TotalCount int64
	Results    []Problem
}

type ProblemRepository interface {
	FindProblems(ctx context.Context, operator Student, param ProblemSearchCondition) (*ProblemSearchResult, error)

	FindProblemsByProblemIDs(ctx context.Context, operator Student, param ProblemIDsCondition) (*ProblemSearchResult, error)

	FindProblemIDs(ctx context.Context, operator Student, workbookID WorkbookID) ([]ProblemID, error)

	FindProblemByID(ctx context.Context, operator Student, workbookID WorkbookID, problemID ProblemID) (Problem, error)

	AddProblem(ctx context.Context, operator Student, param ProblemAddParameter) (ProblemID, error)

	// UpdateProblem(ctx context.Context, operator Student, param ProblemUpdateParameter) error

	RemoveProblem(ctx context.Context, operator Student, problemID ProblemID, version int) error
}
