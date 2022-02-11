package domain_mock

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/stretchr/testify/mock"
)

type StudyTypeRepositoryMock struct {
	mock.Mock
}

func (m *StudyTypeRepositoryMock) FindAllStudyTypes(ctx context.Context) ([]domain.StudyType, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.StudyType), args.Error(1)
}
