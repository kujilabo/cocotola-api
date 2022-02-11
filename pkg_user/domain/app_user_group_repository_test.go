package domain_test

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/stretchr/testify/mock"
)

type AppUserGroupRepositoryMock struct {
	mock.Mock
}

func (m *AppUserGroupRepositoryMock) FindPublicGroup(ctx context.Context, operator domain.SystemOwner) (domain.AppUserGroup, error) {
	args := m.Called()
	return args.Get(0).(domain.AppUserGroup), args.Error(1)
}

func (m *AppUserGroupRepositoryMock) AddPublicGroup(ctx context.Context, operator domain.SystemOwner) (domain.AppUserGroupID, error) {
	args := m.Called()
	return args.Get(0).(domain.AppUserGroupID), args.Error(1)
}
