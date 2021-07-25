package domain

import "context"

type AppUserGroupRepository interface {
	FindPublicGroup(ctx context.Context, operator SystemOwner) (AppUserGroup, error)

	AddPublicGroup(ctx context.Context, operator SystemOwner) (AppUserGroupID, error)
	// AddPersonalGroup(operator SystemOwner, studentID uint) (uint, error)
}
