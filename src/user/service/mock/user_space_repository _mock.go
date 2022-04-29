package service_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/src/user/domain"
)

type UserSpaceRepositoryMock struct {
	mock.Mock
}

func (m *UserSpaceRepositoryMock) Add(ctx context.Context, operator domain.AppUserModel, spaceID domain.SpaceID) error {
	args := m.Called(ctx, operator, spaceID)
	return args.Error(0)
}
