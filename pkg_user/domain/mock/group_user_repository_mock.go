package domain_mock

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/stretchr/testify/mock"
)

type GroupUserRepositoryMock struct {
	mock.Mock
}

func (m *GroupUserRepositoryMock) AddGroupUser(ctx context.Context, operator domain.AppUser, appUserGroupID domain.AppUserGroupID, appUserID domain.AppUserID) error {
	args := m.Called(ctx, operator, appUserGroupID, appUserID)
	return args.Error(0)
}
