package service_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/src/user/domain"
)

type GroupUserRepositoryMock struct {
	mock.Mock
}

func (m *GroupUserRepositoryMock) AddGroupUser(ctx context.Context, operator domain.AppUserModel, appUserGroupID domain.AppUserGroupID, appUserID domain.AppUserID) error {
	args := m.Called(ctx, operator, appUserGroupID, appUserID)
	return args.Error(0)
}
