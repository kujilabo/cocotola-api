package domain

import (
	"context"

	"github.com/go-playground/validator"
	"golang.org/x/xerrors"
)

type Recordbook interface {
	GetWorkbookID() WorkbookID
	GetResults(ctx context.Context) (map[ProblemID]StudyStatus, error)
	GetResultsSortedLevel(ctx context.Context) ([]ProblemWithLevel, error)
	SetResult(ctx context.Context, problemType string, problemID ProblemID, result, memorized bool) error
}

type recordbook struct {
	rf         RepositoryFactory
	student    Student
	workbookID WorkbookID `validate:"required"`
	studyType  string
}

func NewRecordbook(rf RepositoryFactory, student Student, workbookID WorkbookID, studyType string) (Recordbook, error) {
	m := &recordbook{
		rf:         rf,
		student:    student,
		workbookID: workbookID,
		studyType:  studyType,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (m *recordbook) GetWorkbookID() WorkbookID {
	return m.workbookID
}

func (m *recordbook) GetResults(ctx context.Context) (map[ProblemID]StudyStatus, error) {
	repo, err := m.rf.NewRecordbookRepository(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewRecordbookRepository. err: %w", err)
	}

	studyResults, err := repo.FindStudyResults(ctx, m.student, m.workbookID, m.studyType)
	if err != nil {
		return nil, xerrors.Errorf("failed to FindStudyResults. err: %w", err)
	}

	workbook, err := m.student.FindWorkbookByID(ctx, m.workbookID)
	if err != nil {
		return nil, xerrors.Errorf("failed to FindWorkbookByID. err: %w", err)
	}

	problemIDs, err := workbook.FindProblemIDs(ctx, m.student)
	if err != nil {
		return nil, xerrors.Errorf("failed to FindProblemIDs. err: %w", err)
	}

	results := make(map[ProblemID]StudyStatus)
	for _, problemID := range problemIDs {
		if level, ok := studyResults[problemID]; ok {
			results[problemID] = level
		} else {
			results[problemID] = StudyStatus{
				Level:     0,
				Memorized: false,
			}
		}
	}

	return results, nil
}

func (m *recordbook) GetResultsSortedLevel(ctx context.Context) ([]ProblemWithLevel, error) {
	problems1, err := m.GetResults(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to GetResults. err: %w", err)
	}

	problems2 := make([]ProblemWithLevel, len(problems1))
	i := 0
	for k, v := range problems1 {
		problems2[i] = ProblemWithLevel{
			ProblemID: k,
			Level:     v.Level,
			Memorized: v.Memorized,
		}
		i++
	}

	return problems2, nil
}

func (m *recordbook) SetResult(ctx context.Context, problemType string, problemID ProblemID, result, memorized bool) error {
	repo, err := m.rf.NewRecordbookRepository(ctx)
	if err != nil {
		return xerrors.Errorf("failed to NewStudyResultRepository. err: %w", err)
	}

	if err := repo.SetResult(ctx, m.student, m.workbookID, m.studyType, problemType, problemID, result, memorized); err != nil {
		return xerrors.Errorf("failed to SetResult. err: %w", err)
	}

	return nil
}
