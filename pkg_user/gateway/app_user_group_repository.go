package gateway

import (
	"context"
	"time"

	"gorm.io/gorm"

	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/service"
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

func (e *appUserGroupEntity) toAppUserGroup() (service.AppUserGroup, error) {
	model, err := domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	appUserGroupMdoel, err := domain.NewAppUserGroup(model, domain.OrganizationID(e.OrganizationID), e.Key, e.Name, e.Description)
	if err != nil {
		return nil, err
	}

	return service.NewAppUserGroup(appUserGroupMdoel)
}

func NewAppUserGroupRepository(db *gorm.DB) service.AppUserGroupRepository {
	return &appUserGroupRepository{
		db: db,
	}
}

func (r *appUserGroupRepository) FindPublicGroup(ctx context.Context, operator domain.SystemOwnerModel) (service.AppUserGroup, error) {
	appUserGroup := appUserGroupEntity{}
	if result := r.db.Where(&appUserGroupEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		Key:            "public",
	}).Find(&appUserGroup); result.Error != nil {
		return nil, result.Error
	}
	return appUserGroup.toAppUserGroup()
}

func (r *appUserGroupRepository) AddPublicGroup(ctx context.Context, operator domain.SystemOwnerModel) (domain.AppUserGroupID, error) {
	appUserGroup := appUserGroupEntity{
		Version:        1,
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		Key:            "public",
		Name:           "Public group",
	}
	if result := r.db.Create(&appUserGroup); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists)
	}
	return domain.AppUserGroupID(appUserGroup.ID), nil
}

func (r *appUserGroupRepository) AddPersonalGroup(ctx context.Context, operator domain.AppUserModel) (uint, error) {
	appUserGroup := appUserGroupEntity{
		Version:        1,
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		Key:            "#" + operator.GetLoginID(),
		Name:           "Personal group",
	}
	if result := r.db.Create(&appUserGroup); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists)
	}
	return appUserGroup.ID, nil
}
