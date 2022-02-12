package domain_mock

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/stretchr/testify/mock"
)

type UserQuotaRepositoryMock struct {
	mock.Mock
}

func (m *UserQuotaRepositoryMock) IsExceeded(ctx context.Context, operator domain.Student, name string, unit domain.QuotaUnit, limit int) (bool, error) {
	args := m.Called(ctx, operator, name, unit, limit)
	return args.Bool(0), args.Error(1)
}

func (m *UserQuotaRepositoryMock) Increment(ctx context.Context, operator domain.Student, name string, unit domain.QuotaUnit, limit int, count int) (bool, error) {
	args := m.Called(ctx, operator, name, unit, limit, count)
	return args.Bool(0), args.Error(1)
}
