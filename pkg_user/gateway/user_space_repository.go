package gateway

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type userSpaceRepository struct {
	db *gorm.DB
	rf domain.RepositoryFactory
}

type userSpaceEntity struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      uint
	UpdatedBy      uint
	OrganizationID uint
	AppUserID      uint
	SpaceID        uint
}

func (e *userSpaceEntity) TableName() string {
	return "user_space"
}

func NewUserSpaceRepository(rf domain.RepositoryFactory, db *gorm.DB) domain.UserSpaceRepository {
	return &userSpaceRepository{
		db: db,
		rf: rf,
	}
}

func (r *userSpaceRepository) Add(ctx context.Context, operator domain.AppUser, spaceID domain.SpaceID) error {
	if result := r.db.Create(&userSpaceEntity{
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		AppUserID:      operator.GetID(),
		SpaceID:        uint(spaceID),
	}); result.Error != nil {
		return libG.ConvertDuplicatedError(result.Error, domain.ErrAppUserAlreadyExists)
	}

	return nil
}

func (r *userSpaceRepository) Remove(ctx context.Context, operator domain.AppUser, spaceID uint) error {
	if result := r.db.Where(userSpaceEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		AppUserID:      operator.GetID(),
		SpaceID:        spaceID,
	}).Delete(userSpaceEntity{}); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userSpaceRepository) IsBelongedTo(ctx context.Context, operator domain.AppUser, spaceID uint) (bool, error) {
	entity := userSpaceEntity{}
	if result := r.db.Where(userSpaceEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		AppUserID:      operator.GetID(),
		SpaceID:        spaceID,
	}).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, result.Error
	}

	return true, nil
}
