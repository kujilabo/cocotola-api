package domain_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type SpaceRepositoryMock struct {
	mock.Mock
}

func (m *SpaceRepositoryMock) FindDefaultSpace(ctx context.Context, operator domain.AppUser) (domain.Space, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(domain.Space), args.Error(1)
}

func (m *SpaceRepositoryMock) FindPersonalSpace(ctx context.Context, operator domain.AppUser) (domain.Space, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(domain.Space), args.Error(1)
}

func (m *SpaceRepositoryMock) FindSystemSpace(ctx context.Context, operator domain.AppUser) (domain.Space, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(domain.Space), args.Error(1)
}

func (m *SpaceRepositoryMock) AddDefaultSpace(ctx context.Context, operator domain.SystemOwner) (uint, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(uint), args.Error(1)
}

func (m *SpaceRepositoryMock) AddPersonalSpace(ctx context.Context, operator domain.SystemOwner, appUser domain.AppUser) (domain.SpaceID, error) {
	args := m.Called(ctx, operator, appUser)
	return args.Get(0).(domain.SpaceID), args.Error(1)
}

func (m *SpaceRepositoryMock) AddSystemSpace(ctx context.Context, operator domain.SystemOwner) (domain.SpaceID, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(domain.SpaceID), args.Error(1)
}
