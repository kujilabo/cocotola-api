package service

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

type Recordbook interface {
	GetStudent() Student
	GetWorkbookID() domain.WorkbookID
	GetResults(ctx context.Context) (map[domain.ProblemID]domain.StudyStatus, error)
	GetResultsSortedLevel(ctx context.Context) ([]domain.ProblemWithLevel, error)
	SetResult(ctx context.Context, problemType string, problemID domain.ProblemID, result, memorized bool) error
}

type recordbook struct {
	rf         RepositoryFactory
	student    Student
	workbookID domain.WorkbookID `validate:"required"`
	studyType  string
}

func NewRecordbook(rf RepositoryFactory, student Student, workbookID domain.WorkbookID, studyType string) (Recordbook, error) {
	m := &recordbook{
		rf:         rf,
		student:    student,
		workbookID: workbookID,
		studyType:  studyType,
	}

	return m, libD.Validator.Struct(m)
}

func (m *recordbook) GetStudent() Student {
	return m.student
}

func (m *recordbook) GetWorkbookID() domain.WorkbookID {
	return m.workbookID
}

func (m *recordbook) GetResults(ctx context.Context) (map[domain.ProblemID]domain.StudyStatus, error) {
	repo, err := m.rf.NewRecordbookRepository(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewRecordbookRepository. err: %w", err)
	}

	studyResults, err := repo.FindStudyResults(ctx, m.GetStudent(), m.workbookID, m.studyType)
	if err != nil {
		return nil, xerrors.Errorf("failed to FindStudyResults. err: %w", err)
	}

	workbookService, err := m.GetStudent().FindWorkbookByID(ctx, m.workbookID)
	if err != nil {
		return nil, xerrors.Errorf("failed to FindWorkbookByID. err: %w", err)
	}

	problemIDs, err := workbookService.FindProblemIDs(ctx, m.GetStudent())
	if err != nil {
		return nil, xerrors.Errorf("failed to FindProblemIDs. err: %w", err)
	}

	results := make(map[domain.ProblemID]domain.StudyStatus)
	for _, problemID := range problemIDs {
		if status, ok := studyResults[problemID]; ok {
			results[problemID] = status
		} else {
			results[problemID] = domain.StudyStatus{
				Level:     0,
				Memorized: false,
			}
		}
	}

	return results, nil
}

func (m *recordbook) GetResultsSortedLevel(ctx context.Context) ([]domain.ProblemWithLevel, error) {
	problems1, err := m.GetResults(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to GetResults. err: %w", err)
	}

	problems2 := make([]domain.ProblemWithLevel, len(problems1))
	i := 0
	for k, v := range problems1 {
		problems2[i] = domain.ProblemWithLevel{
			ProblemID: k,
			Level:     v.Level,
			Memorized: v.Memorized,
		}
		i++
	}

	return problems2, nil
}

func (m *recordbook) SetResult(ctx context.Context, problemType string, problemID domain.ProblemID, result, memorized bool) error {
	repo, err := m.rf.NewRecordbookRepository(ctx)
	if err != nil {
		return xerrors.Errorf("failed to NewStudyResultRepository. err: %w", err)
	}

	if err := repo.SetResult(ctx, m.GetStudent(), m.workbookID, m.studyType, problemType, problemID, result, memorized); err != nil {
		return xerrors.Errorf("failed to SetResult. err: %w", err)
	}

	return nil
}
