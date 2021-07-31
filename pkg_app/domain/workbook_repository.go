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

type WorkbookAddParameter struct {
	ProblemType  string
	Name         string
	QuestionText string
}

func NewWorkbookAddParameter(problemType string, name, questionText string) (*WorkbookAddParameter, error) {
	m := &WorkbookAddParameter{
		ProblemType:  problemType,
		Name:         name,
		QuestionText: questionText,
	}
	v := validator.New()
	return m, v.Struct(m)
}

type WorkbookUpdateParameter struct {
	Name         string
	QuestionText string
}

func NewWorkbookUpdateParameter(name, questionText string) (*WorkbookUpdateParameter, error) {
	m := &WorkbookUpdateParameter{
		Name:         name,
		QuestionText: questionText,
	}
	v := validator.New()
	return m, v.Struct(m)
}

type WorkbookRepository interface {
	FindWorkbooks(ctx context.Context, operator Student, param *WorkbookSearchCondition) (*WorkbookSearchResult, error)

	FindWorkbookByID(ctx context.Context, operator Student, id WorkbookID) (Workbook, error)

	AddWorkbook(ctx context.Context, operator Student, spaceID user.SpaceID, param *WorkbookAddParameter) (WorkbookID, error)

	UpdateWorkbook(ctx context.Context, operator Student, workbookID WorkbookID, version int, param *WorkbookUpdateParameter) error

	RemoveWorkbook(ctx context.Context, operator Student, workbookID WorkbookID, version int) error
}
