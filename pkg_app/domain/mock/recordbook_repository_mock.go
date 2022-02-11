package domain_mock

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/stretchr/testify/mock"
)

type RecordbookRepositoryMock struct {
	mock.Mock
}

func (m *RecordbookRepositoryMock) FindStudyResults(ctx context.Context, operator domain.Student, workbookID domain.WorkbookID, studyType string) (map[domain.ProblemID]domain.StudyStatus, error) {
	args := m.Called(ctx, operator, workbookID, studyType)
	return args.Get(0).(map[domain.ProblemID]domain.StudyStatus), args.Error(1)
}

func (m *RecordbookRepositoryMock) SetResult(ctx context.Context, operator domain.Student, workbookID domain.WorkbookID, studyType string, problemType string, problemID domain.ProblemID, result, memorized bool) error {
	args := m.Called(ctx, operator, workbookID, studyType)
	return args.Error(0)
}
