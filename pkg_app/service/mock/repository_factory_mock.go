package service_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/pkg_app/service"
)

type RepositoryFactoryMock struct {
	mock.Mock
}

func (m *RepositoryFactoryMock) NewWorkbookRepository(ctx context.Context) (service.WorkbookRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(service.WorkbookRepository), args.Error(1)
}

func (m *RepositoryFactoryMock) NewProblemRepository(ctx context.Context, problemType string) (service.ProblemRepository, error) {
	args := m.Called(ctx, problemType)
	return args.Get(0).(service.ProblemRepository), args.Error(1)
}

func (m *RepositoryFactoryMock) NewProblemTypeRepository(ctx context.Context) (service.ProblemTypeRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(service.ProblemTypeRepository), args.Error(1)
}

func (m *RepositoryFactoryMock) NewStudyTypeRepository(ctx context.Context) (service.StudyTypeRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(service.StudyTypeRepository), args.Error(1)
}

func (m *RepositoryFactoryMock) NewAudioRepository(ctx context.Context) (service.AudioRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(service.AudioRepository), args.Error(1)
}

func (m *RepositoryFactoryMock) NewRecordbookRepository(ctx context.Context) (service.RecordbookRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(service.RecordbookRepository), args.Error(1)
}

func (m *RepositoryFactoryMock) NewUserQuotaRepository(ctx context.Context) (service.UserQuotaRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(service.UserQuotaRepository), args.Error(1)
}
