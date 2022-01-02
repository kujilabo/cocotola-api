package domain

import (
	"context"

	"github.com/go-playground/validator/v10"
	"golang.org/x/xerrors"

	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

var ErrWorkbookNotFound = xerrors.New("Workbook not found")
var ErrWorkbookAlreadyExists = xerrors.New("Workbook already exists")
var ErrWorkbookPermissionDenied = xerrors.New("Permission denied")

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

	v := validator.New()
	return m, v.Struct(m)
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

	v := validator.New()
	return m, v.Struct(m)
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

	v := validator.New()
	return m, v.Struct(m)
}

func (p *workbookUpdateParameter) GetName() string {
	return p.Name
}

func (p *workbookUpdateParameter) GetQuestionText() string {
	return p.QuestionText
}

type WorkbookRepository interface {
	FindPersonalWorkbooks(ctx context.Context, operator Student, param WorkbookSearchCondition) (*WorkbookSearchResult, error)

	FindWorkbookByID(ctx context.Context, operator Student, id WorkbookID) (Workbook, error)

	FindWorkbookByName(ctx context.Context, operator Student, name string) (Workbook, error)

	AddWorkbook(ctx context.Context, operator Student, spaceID user.SpaceID, param WorkbookAddParameter) (WorkbookID, error)

	UpdateWorkbook(ctx context.Context, operator Student, workbookID WorkbookID, version int, param WorkbookUpdateParameter) error

	RemoveWorkbook(ctx context.Context, operator Student, workbookID WorkbookID, version int) error
}
