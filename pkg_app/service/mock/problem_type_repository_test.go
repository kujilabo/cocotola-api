package service_mock

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/stretchr/testify/mock"
)

type ProblemTypeRepositoryMock struct {
	mock.Mock
}

func (m *ProblemTypeRepositoryMock) FindAllProblemTypes(ctx context.Context) ([]domain.ProblemType, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.ProblemType), args.Error(1)
}
