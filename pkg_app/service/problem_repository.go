package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

var ErrProblemAlreadyExists = errors.New("problem already exists")
var ErrProblemNotFound = errors.New("problem not found")
var ErrProblemOtherError = errors.New("problem other error")

type ProblemAddParameter interface {
	GetWorkbookID() domain.WorkbookID
	GetNumber() int
	GetProperties() map[string]string
}

type problemAddParameter struct {
	WorkbookID domain.WorkbookID `validate:"required"`
	Number     int               `validate:"required"`
	Properties map[string]string
}

func NewProblemAddParameter(workbookID domain.WorkbookID, number int, properties map[string]string) (ProblemAddParameter, error) {
	m := &problemAddParameter{
		WorkbookID: workbookID,
		Number:     number,
		Properties: properties,
	}

	return m, libD.Validator.Struct(m)
}

func (p *problemAddParameter) GetWorkbookID() domain.WorkbookID {
	return p.WorkbookID
}
func (p *problemAddParameter) GetNumber() int {
	return p.Number
}
func (p *problemAddParameter) GetProperties() map[string]string {
	return p.Properties
}

type ProblemSelectParameter1 interface {
	GetWorkbookID() domain.WorkbookID
	GetProblemID() domain.ProblemID
}

type problemSelectParameter1 struct {
	WorkbookID domain.WorkbookID
	ProblemID  domain.ProblemID
}

func NewProblemSelectParameter1(WorkbookID domain.WorkbookID, problemID domain.ProblemID) (ProblemSelectParameter1, error) {
	m := &problemSelectParameter1{
		WorkbookID: WorkbookID,
		ProblemID:  problemID,
	}

	return m, libD.Validator.Struct(m)
}

func (p *problemSelectParameter1) GetWorkbookID() domain.WorkbookID {
	return p.WorkbookID
}
func (p *problemSelectParameter1) GetProblemID() domain.ProblemID {
	return p.ProblemID
}

type ProblemSelectParameter2 interface {
	GetWorkbookID() domain.WorkbookID
	GetProblemID() domain.ProblemID
	GetVersion() int
}

type problemSelectParameter2 struct {
	WorkbookID domain.WorkbookID
	ProblemID  domain.ProblemID
	Version    int
}

func NewProblemSelectParameter2(WorkbookID domain.WorkbookID, problemID domain.ProblemID, version int) (ProblemSelectParameter2, error) {
	m := &problemSelectParameter2{
		WorkbookID: WorkbookID,
		ProblemID:  problemID,
		Version:    version,
	}

	return m, libD.Validator.Struct(m)
}

func (p *problemSelectParameter2) GetWorkbookID() domain.WorkbookID {
	return p.WorkbookID
}
func (p *problemSelectParameter2) GetProblemID() domain.ProblemID {
	return p.ProblemID
}

func (p *problemSelectParameter2) GetVersion() int {
	return p.Version
}

type ProblemUpdateParameter interface {
	GetNumber() int
	GetProperties() map[string]string
	GetStringProperty(name string) (string, error)
	GetIntProperty(name string) (int, error)
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
func (p *problemUpdateParameter) GetStringProperty(name string) (string, error) {
	s, ok := p.Properties[name]
	if !ok {
		return "", errors.New("key not found")
	}
	return s, nil
}
func (p *problemUpdateParameter) GetIntProperty(name string) (int, error) {
	i, err := strconv.Atoi(p.Properties[name])
	if err != nil {
		return 0, err
	}
	return i, nil
}

type ProblemSearchCondition interface {
	GetWorkbookID() domain.WorkbookID
	GetPageNo() int
	GetPageSize() int
	GetKeyword() string
}

type problemSearchCondition struct {
	WorkbookID domain.WorkbookID
	PageNo     int `validate:"required,gte=1"`
	PageSize   int `validate:"required,gte=1,lte=1000"`
	Keyword    string
}

func NewProblemSearchCondition(workbookID domain.WorkbookID, pageNo, pageSize int, keyword string) (ProblemSearchCondition, error) {
	m := &problemSearchCondition{
		WorkbookID: workbookID,
		PageNo:     pageNo,
		PageSize:   pageSize,
		Keyword:    keyword,
	}

	return m, libD.Validator.Struct(m)
}

func (c *problemSearchCondition) GetWorkbookID() domain.WorkbookID {
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
	GetWorkbookID() domain.WorkbookID
	GetIDs() []domain.ProblemID
}

type problemIDsCondition struct {
	WorkbookID domain.WorkbookID
	IDs        []domain.ProblemID
}

func NewProblemIDsCondition(workbookID domain.WorkbookID, ids []domain.ProblemID) (ProblemIDsCondition, error) {
	m := &problemIDsCondition{
		WorkbookID: workbookID,
		IDs:        ids,
	}

	return m, libD.Validator.Struct(m)
}

func (c *problemIDsCondition) GetWorkbookID() domain.WorkbookID {
	return c.WorkbookID
}

func (c *problemIDsCondition) GetIDs() []domain.ProblemID {
	return c.IDs
}

type ProblemSearchResult interface {
	GetTotalCount() int
	GetResults() []domain.ProblemModel
}
type problemSearchResult struct {
	TotalCount int
	Results    []domain.ProblemModel
}

func NewProblemSearchResult(totalCount int, results []domain.ProblemModel) (ProblemSearchResult, error) {
	m := &problemSearchResult{
		TotalCount: totalCount,
		Results:    results,
	}

	return m, libD.Validator.Struct(m)
}
func (m *problemSearchResult) GetTotalCount() int {
	return m.TotalCount
}

func (m *problemSearchResult) GetResults() []domain.ProblemModel {
	return m.Results
}

type ProblemRepository interface {
	// FindProblems searches for problems based on search condition
	FindProblems(ctx context.Context, operator domain.StudentModel, param ProblemSearchCondition) (ProblemSearchResult, error)

	FindAllProblems(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID) (ProblemSearchResult, error)

	FindProblemsByProblemIDs(ctx context.Context, operator domain.StudentModel, param ProblemIDsCondition) (ProblemSearchResult, error)

	FindProblemsByCustomCondition(ctx context.Context, operator domain.StudentModel, condition interface{}) ([]domain.ProblemModel, error)

	FindProblemByID(ctx context.Context, operator domain.StudentModel, id ProblemSelectParameter1) (Problem, error)

	FindProblemIDs(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID) ([]domain.ProblemID, error)

	// AddProblem register a new problem
	AddProblem(ctx context.Context, operator domain.StudentModel, param ProblemAddParameter) (domain.ProblemID, error)

	UpdateProblem(ctx context.Context, operator domain.StudentModel, id ProblemSelectParameter2, param ProblemUpdateParameter) error

	RemoveProblem(ctx context.Context, operator domain.StudentModel, id ProblemSelectParameter2) error
}
