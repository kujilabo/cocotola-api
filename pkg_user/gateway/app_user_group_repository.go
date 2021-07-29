package gateway

import (
	"context"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

var (
	AppUserGroupTableName = "app_user_group"
)

type appUserGroupRepository struct {
	db *gorm.DB
}

type appUserGroupEntity struct {
	ID             uint
	Version        int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      uint
	UpdatedBy      uint
	OrganizationID uint
	Key            string
	Name           string
	Description    string
}

func (e *appUserGroupEntity) TableName() string {
	return AppUserGroupTableName
}

func (e *appUserGroupEntity) toAppUserGroup() domain.AppUserGroup {
	return domain.NewAppUserGroup(domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy), domain.OrganizationID(e.OrganizationID), e.Key, e.Name, e.Description)
}

func NewAppUserGroupRepository(db *gorm.DB) domain.AppUserGroupRepository {
	return &appUserGroupRepository{
		db: db,
	}
}

func (r *appUserGroupRepository) FindPublicGroup(ctx context.Context, operator domain.SystemOwner) (domain.AppUserGroup, error) {
	appUserGroup := appUserGroupEntity{}
	if result := r.db.Where(&appUserGroupEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		Key:            "public",
	}).Find(&appUserGroup); result.Error != nil {
		return nil, result.Error
	}
	return appUserGroup.toAppUserGroup(), nil
}

func (r *appUserGroupRepository) AddPublicGroup(ctx context.Context, operator domain.SystemOwner) (domain.AppUserGroupID, error) {
	appUserGroup := appUserGroupEntity{
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		Key:            "public",
		Name:           "Public group",
	}
	if result := r.db.Create(&appUserGroup); result.Error != nil {
		dbErr, ok := result.Error.(*mysql.MySQLError)
		if ok && dbErr.Number == 1062 {
			return 0, domain.ErrAppUserAlreadyExists
		}
		return 0, result.Error
	}
	return domain.AppUserGroupID(appUserGroup.ID), nil
}

func (r *appUserGroupRepository) AddPersonalGroup(ctx context.Context, operator domain.AppUser) (uint, error) {
	appUserGroup := appUserGroupEntity{
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		Key:            "#" + operator.GetLoginID(),
		Name:           "Personal group",
	}
	if result := r.db.Create(&appUserGroup); result.Error != nil {
		dbErr, ok := result.Error.(*mysql.MySQLError)
		if ok && dbErr.Number == 1062 {
			return 0, domain.ErrAppUserAlreadyExists
		}
		return 0, result.Error
	}
	return appUserGroup.ID, nil
}
