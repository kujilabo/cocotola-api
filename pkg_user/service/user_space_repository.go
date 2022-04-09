package service

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

// UserSpaceRepository mangages relationship between AppUser and Space
type UserSpaceRepository interface {
	Add(ctx context.Context, operator domain.AppUserModel, spaceID domain.SpaceID) error
}
