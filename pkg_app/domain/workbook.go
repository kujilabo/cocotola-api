package domain

import (
	"context"
	"errors"

	"github.com/go-playground/validator"
	"golang.org/x/xerrors"

	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type WorkbookID uint

type WorkbookSearchResult struct {
	TotalCount int64
	Results    []Workbook
}

//
type workbookProperties struct {
	Name         string `validate:"required"`
	ProblemType  string `validate:"required"`
	QuestionText string
}

type Workbook interface {
	user.Model
	GetSpaceID() user.SpaceID
	GetOwnerID() user.AppUserID
	GetName() string
	GetProblemType() string
	GetQuestionText() string

	FindProblems(ctx context.Context, operator Student, param ProblemSearchCondition) (*ProblemSearchResult, error)

	FindProblemsByProblemIDs(ctx context.Context, operator Student, param ProblemIDsCondition) (*ProblemSearchResult, error)

	FindProblemIDs(ctx context.Context, operator Student) ([]ProblemID, error)

	FindProblemByID(ctx context.Context, operator Student, problemID ProblemID) (Problem, error)

	AddProblem(ctx context.Context, operator Student, param ProblemAddParameter) (ProblemID, error)

	RemoveProblem(ctx context.Context, operator Student, problemID ProblemID, version int) error

	UpdateWorkbook(ctx context.Context, operator Student, version int, parameter WorkbookUpdateParameter) error

	RemoveWorkbook(ctx context.Context, operator Student, version int) error
}

type workbook struct {
	repo             RepositoryFactory
	processorFactory ProcessorFactory
	user.Model
	spaceID    user.SpaceID    `validate:"required"`
	ownerID    user.AppUserID  `validate:"required"`
	privileges user.Privileges `validate:"required"`
	Properties workbookProperties
}

func NewWorkbook(repo RepositoryFactory, processorFactory ProcessorFactory, model user.Model, spaceID user.SpaceID, ownerID user.AppUserID, privileges user.Privileges, name string, problemType string, questsionText string) (Workbook, error) {
	m := &workbook{
		repo:             repo,
		processorFactory: processorFactory,
		privileges:       privileges,
		Model:            model,
		spaceID:          spaceID,
		ownerID:          ownerID,
		Properties: workbookProperties{
			Name:         name,
			ProblemType:  problemType,
			QuestionText: questsionText,
		},
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (m *workbook) GetSpaceID() user.SpaceID {
	return m.spaceID
}

func (m *workbook) GetOwnerID() user.AppUserID {
	return m.ownerID
}

func (m *workbook) GetName() string {
	return m.Properties.Name
}

func (m *workbook) GetProblemType() string {
	return m.Properties.ProblemType
}

func (m *workbook) GetQuestionText() string {
	return m.Properties.QuestionText
}

func (m *workbook) FindProblems(ctx context.Context, operator Student, param ProblemSearchCondition) (*ProblemSearchResult, error) {
	problemRepo, err := m.repo.NewProblemRepository(ctx, m.GetProblemType())
	if err != nil {
		return nil, err
	}
	return problemRepo.FindProblems(ctx, operator, param)
}

func (m *workbook) FindProblemsByProblemIDs(ctx context.Context, operator Student, param ProblemIDsCondition) (*ProblemSearchResult, error) {
	problemRepo, err := m.repo.NewProblemRepository(ctx, m.GetProblemType())
	if err != nil {
		return nil, err
	}
	return problemRepo.FindProblemsByProblemIDs(ctx, operator, param)
}

func (m *workbook) FindProblemIDs(ctx context.Context, operator Student) ([]ProblemID, error) {
	problemRepo, err := m.repo.NewProblemRepository(ctx, m.GetProblemType())
	if err != nil {
		return nil, err
	}
	return problemRepo.FindProblemIDs(ctx, operator, WorkbookID(m.GetID()))
}

func (m *workbook) FindProblemByID(ctx context.Context, operator Student, problemID ProblemID) (Problem, error) {
	problemRepo, err := m.repo.NewProblemRepository(ctx, m.GetProblemType())
	if err != nil {
		return nil, err
	}
	return problemRepo.FindProblemByID(ctx, operator, WorkbookID(m.GetID()), problemID)
}

func (m *workbook) AddProblem(ctx context.Context, operator Student, param ProblemAddParameter) (ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("workbook.AddProblem")

	if !m.privileges.HasPrivilege(PrivilegeUpdate) {
		return 0, errors.New("no update privilege")
	}

	processor, err := m.processorFactory.NewProblemAddProcessor(param.GetProblemType())
	if err != nil {
		return 0, xerrors.Errorf("processor not found. problemType: %s, err: %w", param.GetProblemType(), err)
	}

	logger.Infof("processor.AddProblem")
	return processor.AddProblem(ctx, m.repo, operator, param)
}

func (m *workbook) RemoveProblem(ctx context.Context, operator Student, problemID ProblemID, version int) error {
	logger := log.FromContext(ctx)
	logger.Infof("workbook.RemoveProblem")

	if !m.privileges.HasPrivilege(PrivilegeUpdate) {
		return errors.New("no update privilege")
	}

	processor, err := m.processorFactory.NewProblemRemoveProcessor(m.GetProblemType())
	if err != nil {
		return xerrors.Errorf("processor not found. problemType: %s, err: %w", m.GetProblemType(), err)
	}

	return processor.RemoveProblem(ctx, m.repo, operator, problemID, version)

}

func (m *workbook) UpdateWorkbook(ctx context.Context, operator Student, version int, parameter WorkbookUpdateParameter) error {
	if !m.privileges.HasPrivilege(PrivilegeUpdate) {
		return ErrWorkbookPermissionDenied
	}

	workbookRepo, err := m.repo.NewWorkbookRepository(ctx)
	if err != nil {
		return xerrors.Errorf("failed to NewWorkbookRepository. err: %v", err)
	}

	return workbookRepo.UpdateWorkbook(ctx, operator, WorkbookID(m.GetID()), version, parameter)
}

func (m *workbook) RemoveWorkbook(ctx context.Context, operator Student, version int) error {
	if !m.privileges.HasPrivilege(PrivilegeRemove) {
		return ErrWorkbookPermissionDenied
	}

	workbookRepo, err := m.repo.NewWorkbookRepository(ctx)
	if err != nil {
		return xerrors.Errorf("failed to NewWorkbookRepository. err: %v", err)
	}

	return workbookRepo.RemoveWorkbook(ctx, operator, WorkbookID(m.GetID()), version)
}
