package domain_mock

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/stretchr/testify/mock"
)

type RepositoryFactoryMock struct {
	mock.Mock
}

func (m *RepositoryFactoryMock) NewWorkbookRepository(ctx context.Context) (domain.WorkbookRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.WorkbookRepository), args.Error(1)
}

func (m *RepositoryFactoryMock) NewProblemRepository(ctx context.Context, problemType string) (domain.ProblemRepository, error) {
	args := m.Called(ctx, problemType)
	return args.Get(0).(domain.ProblemRepository), args.Error(1)
}

func (m *RepositoryFactoryMock) NewProblemTypeRepository(ctx context.Context) (domain.ProblemTypeRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.ProblemTypeRepository), args.Error(1)
}

func (m *RepositoryFactoryMock) NewStudyTypeRepository(ctx context.Context) (domain.StudyTypeRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.StudyTypeRepository), args.Error(1)
}

func (m *RepositoryFactoryMock) NewAudioRepository(ctx context.Context) (domain.AudioRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.AudioRepository), args.Error(1)
}

func (m *RepositoryFactoryMock) NewRecordbookRepository(ctx context.Context) (domain.RecordbookRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.RecordbookRepository), args.Error(1)
}

func (m *RepositoryFactoryMock) NewUserQuotaRepository(ctx context.Context) (domain.UserQuotaRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.UserQuotaRepository), args.Error(1)
}
