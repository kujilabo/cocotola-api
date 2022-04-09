package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	"golang.org/x/xerrors"
)

// type WorkbookID uint

type Workbook interface {
	domain.WorkbookModel
	// GetWorkbookModel() domain.WorkbookModel

	// user.Model
	// GetSpaceID() user.SpaceID
	// GetOwnerID() user.AppUserID
	// GetName() string
	// GetProblemType() string
	// GetQuestionText() string
	// GetProperties() map[string]string

	FindProblems(ctx context.Context, operator domain.StudentModel, param ProblemSearchCondition) (ProblemSearchResult, error)

	FindAllProblems(ctx context.Context, operator domain.StudentModel) (ProblemSearchResult, error)

	FindProblemsByProblemIDs(ctx context.Context, operator domain.StudentModel, param ProblemIDsCondition) (ProblemSearchResult, error)

	FindProblemIDs(ctx context.Context, operator domain.StudentModel) ([]domain.ProblemID, error)

	FindProblemByID(ctx context.Context, operator domain.StudentModel, problemID domain.ProblemID) (Problem, error)

	AddProblem(ctx context.Context, operator domain.StudentModel, param ProblemAddParameter) (Added, domain.ProblemID, error)

	UpdateProblem(ctx context.Context, operator domain.StudentModel, id ProblemSelectParameter2, param ProblemUpdateParameter) (Added, Updated, error)

	RemoveProblem(ctx context.Context, operator domain.StudentModel, id ProblemSelectParameter2) error

	UpdateWorkbook(ctx context.Context, operator domain.StudentModel, version int, parameter WorkbookUpdateParameter) error

	RemoveWorkbook(ctx context.Context, operator domain.StudentModel, version int) error
}

type workbook struct {
	domain.WorkbookModel
	rf RepositoryFactory
	pf ProcessorFactory
	// WorkbookModel domain.WorkbookModel
	// user.Model
	// spaceID      user.SpaceID    `validate:"required"`
	// ownerID      user.AppUserID  `validate:"required"`
	// privileges   user.Privileges `validate:"required"`
	// Name         string          `validate:"required"`
	// ProblemType  string          `validate:"required"`
	// QuestionText string
	// Properties   map[string]string
}

func NewWorkbook(rf RepositoryFactory, pf ProcessorFactory,
	// model user.Model, spaceID user.SpaceID, ownerID user.AppUserID, privileges user.Privileges, name string, problemType string, questsionText string, properties map[string]string
	workbookModel domain.WorkbookModel) (Workbook, error) {
	m := &workbook{
		WorkbookModel: workbookModel,
		rf:            rf,
		pf:            pf,
		// WorkbookModel: workbookModel,
		// privileges: privileges,
		// Model:      model,
		// spaceID:    spaceID,
		// ownerID:    ownerID,
		// // Properties: workbookProperties{
		// // 	Name:         name,
		// // 	ProblemType:  problemType,
		// // 	QuestionText: questsionText,
		// // },

		// Name:         name,
		// ProblemType:  problemType,
		// QuestionText: questsionText,
		// Properties:   properties,
	}

	return m, lib.Validator.Struct(m)
}

func (m *workbook) GetWorkbookModel() domain.WorkbookModel {
	return m.WorkbookModel
}

// func (m *workbook) GetSpaceID() user.SpaceID {
// 	return m.spaceID
// }

// func (m *workbook) GetOwnerID() user.AppUserID {
// 	return m.ownerID
// }

// func (m *workbook) GetName() string {
// 	return m.Name
// }

// func (m *workbook) GetProblemType() string {
// 	return m.ProblemType
// }

// func (m *workbook) GetQuestionText() string {
// 	return m.QuestionText
// }

// func (m *workbook) GetProperties() map[string]string {
// 	return m.Properties
// }

func (m *workbook) FindProblems(ctx context.Context, operator domain.StudentModel, param ProblemSearchCondition) (ProblemSearchResult, error) {
	problemRepo, err := m.rf.NewProblemRepository(ctx, m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return nil, err
	}
	return problemRepo.FindProblems(ctx, operator, param)
}

func (m *workbook) FindAllProblems(ctx context.Context, operator domain.StudentModel) (ProblemSearchResult, error) {
	problemRepo, err := m.rf.NewProblemRepository(ctx, m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return nil, err
	}
	return problemRepo.FindAllProblems(ctx, operator, domain.WorkbookID(m.GetWorkbookModel().GetID()))
}

func (m *workbook) FindProblemsByProblemIDs(ctx context.Context, operator domain.StudentModel, param ProblemIDsCondition) (ProblemSearchResult, error) {
	problemRepo, err := m.rf.NewProblemRepository(ctx, m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return nil, err
	}
	return problemRepo.FindProblemsByProblemIDs(ctx, operator, param)
}

func (m *workbook) FindProblemIDs(ctx context.Context, operator domain.StudentModel) ([]domain.ProblemID, error) {
	problemRepo, err := m.rf.NewProblemRepository(ctx, m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return nil, err
	}
	return problemRepo.FindProblemIDs(ctx, operator, domain.WorkbookID(m.GetWorkbookModel().GetID()))
}

func (m *workbook) FindProblemByID(ctx context.Context, operator domain.StudentModel, problemID domain.ProblemID) (Problem, error) {
	problemRepo, err := m.rf.NewProblemRepository(ctx, m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return nil, err
	}
	id, err := NewProblemSelectParameter1(domain.WorkbookID(m.GetWorkbookModel().GetID()), problemID)
	if err != nil {
		return nil, err
	}
	return problemRepo.FindProblemByID(ctx, operator, id)
}

func (m *workbook) AddProblem(ctx context.Context, operator domain.StudentModel, param ProblemAddParameter) (Added, domain.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("workbook.AddProblem")

	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return 0, 0, errors.New("no update privilege")
	}

	processor, err := m.pf.NewProblemAddProcessor(m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return 0, 0, xerrors.Errorf("processor not found. problemType: %s, err: %w", m.GetWorkbookModel().GetProblemType(), err)
	}

	logger.Infof("processor.AddProblem")
	return processor.AddProblem(ctx, m.rf, operator, m.GetWorkbookModel(), param)
}

func (m *workbook) UpdateProblem(ctx context.Context, operator domain.StudentModel, id ProblemSelectParameter2, param ProblemUpdateParameter) (Added, Updated, error) {
	logger := log.FromContext(ctx)
	logger.Infof("workbook.UpdateProblem")

	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return 0, 0, errors.New("no update privilege")
	}

	processor, err := m.pf.NewProblemUpdateProcessor(m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return 0, 0, xerrors.Errorf("processor not found. problemType: %s, err: %w", m.GetWorkbookModel().GetProblemType(), err)
	}

	return processor.UpdateProblem(ctx, m.rf, operator, m.GetWorkbookModel(), id, param)
}

func (m *workbook) RemoveProblem(ctx context.Context, operator domain.StudentModel, id ProblemSelectParameter2) error {
	logger := log.FromContext(ctx)
	logger.Infof("workbook.RemoveProblem")

	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return errors.New("no update privilege")
	}

	processor, err := m.pf.NewProblemRemoveProcessor(m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return xerrors.Errorf("processor not found. problemType: %s, err: %w", m.GetWorkbookModel().GetProblemType(), err)
	}

	return processor.RemoveProblem(ctx, m.rf, operator, id)

}

func (m *workbook) UpdateWorkbook(ctx context.Context, operator domain.StudentModel, version int, parameter WorkbookUpdateParameter) error {
	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return ErrWorkbookPermissionDenied
	}

	workbookRepo, err := m.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return xerrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

	return workbookRepo.UpdateWorkbook(ctx, operator, domain.WorkbookID(m.GetWorkbookModel().GetID()), version, parameter)
}

func (m *workbook) RemoveWorkbook(ctx context.Context, operator domain.StudentModel, version int) error {
	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeRemove) {
		return ErrWorkbookPermissionDenied
	}

	workbookRepo, err := m.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return xerrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

	return workbookRepo.RemoveWorkbook(ctx, operator, domain.WorkbookID(m.GetWorkbookModel().GetID()), version)
}
