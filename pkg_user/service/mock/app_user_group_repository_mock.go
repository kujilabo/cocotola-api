package service_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/service"
)

type AppUserGroupRepositoryMock struct {
	mock.Mock
}

func (m *AppUserGroupRepositoryMock) FindPublicGroup(ctx context.Context, operator domain.SystemOwnerModel) (service.AppUserGroup, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(service.AppUserGroup), args.Error(1)
}

func (m *AppUserGroupRepositoryMock) AddPublicGroup(ctx context.Context, operator domain.SystemOwnerModel) (domain.AppUserGroupID, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(domain.AppUserGroupID), args.Error(1)
}
