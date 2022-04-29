package service_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola-api/src/user/service"
)

type SpaceRepositoryMock struct {
	mock.Mock
}

func (m *SpaceRepositoryMock) FindDefaultSpace(ctx context.Context, operator domain.AppUserModel) (service.Space, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(service.Space), args.Error(1)
}

func (m *SpaceRepositoryMock) FindPersonalSpace(ctx context.Context, operator domain.AppUserModel) (service.Space, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(service.Space), args.Error(1)
}

func (m *SpaceRepositoryMock) FindSystemSpace(ctx context.Context, operator domain.AppUserModel) (service.Space, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(service.Space), args.Error(1)
}

func (m *SpaceRepositoryMock) AddDefaultSpace(ctx context.Context, operator domain.SystemOwnerModel) (uint, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(uint), args.Error(1)
}

func (m *SpaceRepositoryMock) AddPersonalSpace(ctx context.Context, operator domain.SystemOwnerModel, appUser domain.AppUserModel) (domain.SpaceID, error) {
	args := m.Called(ctx, operator, appUser)
	return args.Get(0).(domain.SpaceID), args.Error(1)
}

func (m *SpaceRepositoryMock) AddSystemSpace(ctx context.Context, operator domain.SystemOwnerModel) (domain.SpaceID, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(domain.SpaceID), args.Error(1)
}
