package domain

import (
	"context"
)

type GroupUserRepository interface {
	AddGroupUser(ctx context.Context, operator AppUser, appUserGroupID AppUserGroupID, appUserID AppUserID) error
}
