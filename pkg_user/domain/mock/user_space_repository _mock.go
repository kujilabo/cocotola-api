package domain_mock

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/stretchr/testify/mock"
)

type UserSpaceRepositoryMock struct {
	mock.Mock
}

func (m *UserSpaceRepositoryMock) Add(ctx context.Context, operator domain.AppUser, spaceID domain.SpaceID) error {
	args := m.Called(ctx, operator, spaceID)
	return args.Error(0)
}
