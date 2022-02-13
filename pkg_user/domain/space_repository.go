package domain

import (
	"context"
	"errors"
)

var ErrSpaceNotFound = errors.New("space not found")
var ErrSpaceAlreadyExists = errors.New("space already exists")

type SpaceRepository interface {
	FindDefaultSpace(ctx context.Context, operator AppUser) (Space, error)

	FindPersonalSpace(ctx context.Context, operator AppUser) (Space, error)

	FindSystemSpace(ctx context.Context, operator AppUser) (Space, error)

	AddDefaultSpace(ctx context.Context, operator SystemOwner) (uint, error)

	AddPersonalSpace(ctx context.Context, operator SystemOwner, appUser AppUser) (SpaceID, error)

	AddSystemSpace(ctx context.Context, operator SystemOwner) (SpaceID, error)
}
