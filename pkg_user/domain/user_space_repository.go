package domain

import "context"

// UserSpaceRepository mangages relationship between AppUser and Space
type UserSpaceRepository interface {
	Add(ctx context.Context, operator AppUser, spaceID SpaceID) error
}
