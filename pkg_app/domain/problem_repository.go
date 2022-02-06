package domain

import (
	"context"
	"errors"

	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

var ErrProblemAlreadyExists = errors.New("problem already exists")
var ErrProblemNotFound = errors.New("problem not found")
var ErrProblemOtherError = errors.New("problem other error")

type ProblemAddParameter interface {
	GetWorkbookID() WorkbookID
	GetNumber() int
	GetProperties() map[string]string
}

type problemAddParameter struct {
	WorkbookID WorkbookID `validate:"required"`
	Number     int        `validate:"required"`
	Properties map[string]string
}

func NewProblemAddParameter(workbookID WorkbookID, number int, properties map[string]string) (ProblemAddParameter, error) {
	m := &problemAddParameter{
		WorkbookID: workbookID,
		Number:     number,
		Properties: properties,
	}

	return m, libD.Validator.Struct(m)
}

func (p *problemAddParameter) GetWorkbookID() WorkbookID {
	return p.WorkbookID
}
func (p *problemAddParameter) GetNumber() int {
	return p.Number
}
func (p *problemAddParameter) GetProperties() map[string]string {
	return p.Properties
}

type ProblemSelectParameter1 interface {
	GetWorkbookID() WorkbookID
	GetProblemID() ProblemID
}

type problemSelectParameter1 struct {
	WorkbookID WorkbookID
	ProblemID  ProblemID
}

func NewProblemSelectParameter1(WorkbookID WorkbookID, problemID ProblemID) (ProblemSelectParameter1, error) {
	m := &problemSelectParameter1{
		WorkbookID: WorkbookID,
		ProblemID:  problemID,
	}

	return m, libD.Validator.Struct(m)
}

func (p *problemSelectParameter1) GetWorkbookID() WorkbookID {
	return p.WorkbookID
}
func (p *problemSelectParameter1) GetProblemID() ProblemID {
	return p.ProblemID
}

type ProblemSelectParameter2 interface {
	GetWorkbookID() WorkbookID
	GetProblemID() ProblemID
	GetVersion() int
}

type problemSelectParameter2 struct {
	WorkbookID WorkbookID
	ProblemID  ProblemID
	Version    int
}

func NewProblemSelectParameter2(WorkbookID WorkbookID, problemID ProblemID, version int) (ProblemSelectParameter2, error) {
	m := &problemSelectParameter2{
		WorkbookID: WorkbookID,
		ProblemID:  problemID,
		Version:    version,
	}

	return m, libD.Validator.Struct(m)
}

func (p *problemSelectParameter2) GetWorkbookID() WorkbookID {
	return p.WorkbookID
}
func (p *problemSelectParameter2) GetProblemID() ProblemID {
	return p.ProblemID
}

func (p *problemSelectParameter2) GetVersion() int {
	return p.Version
}

type ProblemUpdateParameter interface {
	GetNumber() int
	GetProperties() map[string]string
}

type problemUpdateParameter struct {
	Number     int `validate:"required"`
	Properties map[string]string
}

func NewProblemUpdateParameter(number int, properties map[string]string) (ProblemUpdateParameter, error) {
	m := &problemUpdateParameter{
		Number:     number,
		Properties: properties,
	}

	return m, libD.Validator.Struct(m)
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
	PageSize   int `validate:"required,gte=1,lte=1000"`
	Keyword    string
}

func NewProblemSearchCondition(workbookID WorkbookID, pageNo, pageSize int, keyword string) (ProblemSearchCondition, error) {
	m := &problemSearchCondition{
		WorkbookID: workbookID,
		PageNo:     pageNo,
		PageSize:   pageSize,
		Keyword:    keyword,
	}

	return m, libD.Validator.Struct(m)
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

	return m, libD.Validator.Struct(m)
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

	FindAllProblems(ctx context.Context, operator Student, workbookID WorkbookID) (*ProblemSearchResult, error)

	FindProblemsByProblemIDs(ctx context.Context, operator Student, param ProblemIDsCondition) (*ProblemSearchResult, error)

	FindProblemIDs(ctx context.Context, operator Student, workbookID WorkbookID) ([]ProblemID, error)

	FindProblemByID(ctx context.Context, operator Student, id ProblemSelectParameter1) (Problem, error)

	AddProblem(ctx context.Context, operator Student, param ProblemAddParameter) (ProblemID, error)

	UpdateProblem(ctx context.Context, operator Student, id ProblemSelectParameter2, param ProblemUpdateParameter) error

	RemoveProblem(ctx context.Context, operator Student, id ProblemSelectParameter2) error
}
