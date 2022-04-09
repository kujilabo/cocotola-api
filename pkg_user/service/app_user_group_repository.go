package service

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type AppUserGroupRepository interface {
	FindPublicGroup(ctx context.Context, operator domain.SystemOwnerModel) (AppUserGroup, error)

	AddPublicGroup(ctx context.Context, operator domain.SystemOwnerModel) (domain.AppUserGroupID, error)
	// AddPersonalGroup(operator SystemOwner, studentID uint) (uint, error)
}
