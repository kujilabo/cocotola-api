package gateway

import (
	"context"
	"time"

	"gorm.io/gorm"

	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

var (
	GroupUserTableName = "group_user"
)

type groupUserRepository struct {
	db *gorm.DB
}

type groupUserEntity struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      uint
	UpdatedBy      uint
	OrganizationID uint
	AppUserGroupID uint
	AppUserID      uint
}

func (u *groupUserEntity) TableName() string {
	return GroupUserTableName
}

func NewGroupUserRepository(db *gorm.DB) domain.GroupUserRepository {
	return &groupUserRepository{
		db: db,
	}
}

func (r *groupUserRepository) AddGroupUser(ctx context.Context, operator domain.AppUser, appUserGroupID domain.AppUserGroupID, appUserID domain.AppUserID) error {
	groupUser := groupUserEntity{
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		AppUserGroupID: uint(appUserGroupID),
		AppUserID:      uint(appUserID),
	}
	if result := r.db.Create(&groupUser); result.Error != nil {
		return libG.ConvertDuplicatedError(result.Error, domain.ErrAppUserAlreadyExists)
	}
	return nil
}
