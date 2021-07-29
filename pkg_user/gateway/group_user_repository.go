package gateway

import (
	"context"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

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
		dbErr, ok := result.Error.(*mysql.MySQLError)
		if ok && dbErr.Number == 1062 {
			return domain.ErrAppUserAlreadyExists
		}
		return result.Error
	}
	return nil
}
