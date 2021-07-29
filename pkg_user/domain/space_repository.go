package domain

import (
	"context"

	"golang.org/x/xerrors"
)

var ErrSpaceNotFound = xerrors.New("Space not found")
var ErrSpaceAlreadyExists = xerrors.New("Space already exists")

type SpaceRepository interface {
	FindDefaultSpace(ctx context.Context, operator AppUser) (Space, error)

	FindPersonalSpace(ctx context.Context, operator AppUser) (Space, error)

	AddDefaultSpace(ctx context.Context, operator SystemOwner) (uint, error)

	AddPersonalSpace(ctx context.Context, operator SystemOwner, appUser AppUser) (SpaceID, error)
}
