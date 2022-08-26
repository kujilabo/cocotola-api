//go:generate mockery --output mock --name Recordbook
package service

import (
	"context"

	"github.com/kujilabo/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	liberrors "github.com/kujilabo/cocotola-api/src/lib/errors"
)

type Recordbook interface {
	GetStudent() Student

	GetWorkbookID() domain.WorkbookID

	GetResults(ctx context.Context) (map[domain.ProblemID]domain.StudyRecord, error)

	GetResultsSortedLevel(ctx context.Context) ([]domain.StudyRecordWithProblemID, error)

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

func (m *recordbook) GetResults(ctx context.Context) (map[domain.ProblemID]domain.StudyRecord, error) {
	repo := m.rf.NewRecordbookRepository(ctx)

	studyResults, err := repo.FindStudyRecords(ctx, m.GetStudent(), m.workbookID, m.studyType)
	if err != nil {
		return nil, liberrors.Errorf("failed to FindStudyResults. err: %w", err)
	}

	workbookService, err := m.GetStudent().FindWorkbookByID(ctx, m.workbookID)
	if err != nil {
		return nil, liberrors.Errorf("failed to FindWorkbookByID. err: %w", err)
	}

	problemIDs, err := workbookService.FindProblemIDs(ctx, m.GetStudent())
	if err != nil {
		return nil, liberrors.Errorf("failed to FindProblemIDs. err: %w", err)
	}

	results := make(map[domain.ProblemID]domain.StudyRecord)
	for _, problemID := range problemIDs {
		if status, ok := studyResults[problemID]; ok {
			results[problemID] = status
		} else {
			results[problemID] = domain.StudyRecord{
				Level:          0,
				ResultPrev1:    false,
				Memorized:      false,
				LastAnsweredAt: nil,
			}
		}
	}

	return results, nil
}

func (m *recordbook) GetResultsSortedLevel(ctx context.Context) ([]domain.StudyRecordWithProblemID, error) {
	problems1, err := m.GetResults(ctx)
	if err != nil {
		return nil, liberrors.Errorf("failed to GetResults. err: %w", err)
	}

	problems2 := make([]domain.StudyRecordWithProblemID, len(problems1))
	i := 0
	for k, v := range problems1 {
		problems2[i] = domain.StudyRecordWithProblemID{
			ProblemID: k,
			StudyRecord: domain.StudyRecord{
				Level:          v.Level,
				ResultPrev1:    v.ResultPrev1,
				Memorized:      v.Memorized,
				LastAnsweredAt: v.LastAnsweredAt,
			},
		}
		i++
	}

	return problems2, nil
}

func (m *recordbook) SetResult(ctx context.Context, problemType string, problemID domain.ProblemID, result, memorized bool) error {
	repo := m.rf.NewRecordbookRepository(ctx)

	if err := repo.SetResult(ctx, m.GetStudent(), m.workbookID, m.studyType, problemType, problemID, result, memorized); err != nil {
		return liberrors.Errorf("failed to SetResult. err: %w", err)
	}

	return nil
}

type RecordbookSummary interface {
	GetCompletionRate(ctx context.Context) (map[string]int, error)
}

type recordbookSummary struct {
	rf         RepositoryFactory
	student    Student
	workbookID domain.WorkbookID `validate:"required"`
}

func (m *recordbookSummary) GetStudent() Student {
	return m.student
}

func (m *recordbookSummary) GetWorkbookID() domain.WorkbookID {
	return m.workbookID
}

func NewRecordbookSummary(rf RepositoryFactory, student Student, workbookID domain.WorkbookID) (RecordbookSummary, error) {
	m := &recordbookSummary{
		rf:         rf,
		student:    student,
		workbookID: workbookID,
	}

	return m, libD.Validator.Struct(m)
}
func (m *recordbookSummary) GetCompletionRate(ctx context.Context) (map[string]int, error) {
	rateMax := 100
	repo := m.rf.NewRecordbookRepository(ctx)

	numberOfMemorizedProblemsMap, err := repo.CountMemorizedProblem(ctx, m.GetStudent(), m.workbookID)
	if err != nil {
		return nil, liberrors.Errorf("failed to SetResult. err: %w", err)
	}

	workbookService, err := m.GetStudent().FindWorkbookByID(ctx, m.workbookID)
	if err != nil {
		return nil, liberrors.Errorf("failed to FindWorkbookByID. err: %w", err)
	}

	numberOfProblems, err := workbookService.CountProblems(ctx, m.GetStudent())
	if err != nil {
		return nil, liberrors.Errorf("failed to SetResult. err: %w", err)
	}

	completionRateMap := map[string]int{}
	for studyType, numberOfMemorizedProblems := range numberOfMemorizedProblemsMap {
		if numberOfProblems == 0 {
			completionRateMap[studyType] = 0
		} else {
			completionRateMap[studyType] = numberOfMemorizedProblems * rateMax / numberOfProblems
		}
	}

	return completionRateMap, nil
}
