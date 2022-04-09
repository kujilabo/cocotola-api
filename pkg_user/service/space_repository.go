package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

var ErrSpaceNotFound = errors.New("space not found")
var ErrSpaceAlreadyExists = errors.New("space already exists")

type SpaceRepository interface {
	FindDefaultSpace(ctx context.Context, operator domain.AppUserModel) (Space, error)

	FindPersonalSpace(ctx context.Context, operator domain.AppUserModel) (Space, error)

	FindSystemSpace(ctx context.Context, operator domain.AppUserModel) (Space, error)

	AddDefaultSpace(ctx context.Context, operator domain.SystemOwnerModel) (uint, error)

	AddPersonalSpace(ctx context.Context, operator domain.SystemOwnerModel, appUser domain.AppUserModel) (domain.SpaceID, error)

	AddSystemSpace(ctx context.Context, operator domain.SystemOwnerModel) (domain.SpaceID, error)
}
