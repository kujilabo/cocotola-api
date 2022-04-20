//go:generate mockery --output mock --name WorkbookRepository
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

var ErrWorkbookNotFound = errors.New("workbook not found")
var ErrWorkbookAlreadyExists = errors.New("workbook already exists")
var ErrWorkbookPermissionDenied = errors.New("permission denied")

type WorkbookSearchCondition interface {
	GetPageNo() int
	GetPageSize() int
	GetSpaceIDs() []user.SpaceID
}

type workbookSearchCondition struct {
	PageNo   int
	PageSize int
	SpaceIDs []user.SpaceID
}

func NewWorkbookSearchCondition(pageNo, pageSize int, spaceIDs []user.SpaceID) (WorkbookSearchCondition, error) {
	m := &workbookSearchCondition{
		PageNo:   pageNo,
		PageSize: pageSize,
		SpaceIDs: spaceIDs,
	}

	return m, libD.Validator.Struct(m)
}

func (p *workbookSearchCondition) GetPageNo() int {
	return p.PageNo
}

func (p *workbookSearchCondition) GetPageSize() int {
	return p.PageSize
}

func (p *workbookSearchCondition) GetSpaceIDs() []user.SpaceID {
	return p.SpaceIDs
}

type WorkbookSearchResult interface {
	GetTotalCount() int
	GetResults() []domain.WorkbookModel
}

type workbookSearchResult struct {
	TotalCount int
	Results    []domain.WorkbookModel
}

func NewWorkbookSearchResult(totalCount int, results []domain.WorkbookModel) (WorkbookSearchResult, error) {
	m := &workbookSearchResult{
		TotalCount: totalCount,
		Results:    results,
	}

	return m, libD.Validator.Struct(m)
}
func (m *workbookSearchResult) GetTotalCount() int {
	return m.TotalCount
}

func (m *workbookSearchResult) GetResults() []domain.WorkbookModel {
	return m.Results
}

type WorkbookAddParameter interface {
	GetProblemType() string
	GetName() string
	GetQuestionText() string
	GetProperties() map[string]string
}

type workbookAddParameter struct {
	ProblemType  string
	Name         string
	QuestionText string
	Properties   map[string]string
}

func NewWorkbookAddParameter(problemType string, name, questionText string, properties map[string]string) (WorkbookAddParameter, error) {
	m := &workbookAddParameter{
		ProblemType:  problemType,
		Name:         name,
		QuestionText: questionText,
		Properties:   properties,
	}

	return m, libD.Validator.Struct(m)
}

func (p *workbookAddParameter) GetProblemType() string {
	return p.ProblemType
}

func (p *workbookAddParameter) GetName() string {
	return p.Name
}

func (p *workbookAddParameter) GetQuestionText() string {
	return p.QuestionText
}

func (p *workbookAddParameter) GetProperties() map[string]string {
	return p.Properties
}

type WorkbookUpdateParameter interface {
	GetName() string
	GetQuestionText() string
}

type workbookUpdateParameter struct {
	Name         string
	QuestionText string
}

func NewWorkbookUpdateParameter(name, questionText string) (WorkbookUpdateParameter, error) {
	m := &workbookUpdateParameter{
		Name:         name,
		QuestionText: questionText,
	}

	return m, libD.Validator.Struct(m)
}

func (p *workbookUpdateParameter) GetName() string {
	return p.Name
}

func (p *workbookUpdateParameter) GetQuestionText() string {
	return p.QuestionText
}

type WorkbookRepository interface {
	FindPersonalWorkbooks(ctx context.Context, operator domain.StudentModel, param WorkbookSearchCondition) (WorkbookSearchResult, error)

	FindWorkbookByID(ctx context.Context, operator domain.StudentModel, id domain.WorkbookID) (Workbook, error)

	FindWorkbookByName(ctx context.Context, operator user.AppUserModel, spaceID user.SpaceID, name string) (Workbook, error)

	AddWorkbook(ctx context.Context, operator user.AppUserModel, spaceID user.SpaceID, param WorkbookAddParameter) (domain.WorkbookID, error)

	UpdateWorkbook(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, version int, param WorkbookUpdateParameter) error

	RemoveWorkbook(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, version int) error
}
