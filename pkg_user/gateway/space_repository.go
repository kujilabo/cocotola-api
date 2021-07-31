package gateway

import (
	"context"
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"

	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

const SpaceTypeDefault = 1
const SpaceTypePersonal = 2

type spaceEntity struct {
	ID             uint
	Version        int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      uint
	UpdatedBy      uint
	OrganizationID uint
	Type           int
	Key            string
	Name           string
	Description    string
}

func (e *spaceEntity) TableName() string {
	return "space"
}

func (e *spaceEntity) toSpace() (domain.Space, error) {
	model, err := domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	return domain.NewSpace(model, domain.OrganizationID(e.OrganizationID), e.Type, e.Key, e.Name, e.Description)
}

type spaceRepository struct {
	db *gorm.DB
}

func NewSpaceRepository(db *gorm.DB) domain.SpaceRepository {
	return &spaceRepository{
		db: db,
	}
}

func (r *spaceRepository) FindDefaultSpace(ctx context.Context, operator domain.AppUser) (domain.Space, error) {
	logger := log.FromContext(ctx)

	space := spaceEntity{}
	result := r.db.Where(&spaceEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		Type:           SpaceTypeDefault,
		Key:            "default",
	}).First(&space)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSpaceNotFound
		}
		return nil, result.Error
	}

	logger.Debugf("%d, %+v", result.RowsAffected, space)
	return space.toSpace()
}

func (r *spaceRepository) FindPersonalSpace(ctx context.Context, operator domain.AppUser) (domain.Space, error) {
	logger := log.FromContext(ctx)
	logger.Infof("operator %+v", operator)

	space := spaceEntity{}
	if result := r.db.Where(&spaceEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		Type:           SpaceTypePersonal,
		Key:            strconv.Itoa(int(operator.GetID())),
	}).First(&space); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSpaceNotFound
		}
		return nil, result.Error
	}

	return space.toSpace()
}

func (r *spaceRepository) AddDefaultSpace(ctx context.Context, operator domain.SystemOwner) (uint, error) {
	space := spaceEntity{
		Version:        1,
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		Type:           SpaceTypeDefault,
		Key:            "default",
		Name:           "Default",
		Description:    "",
	}
	if result := r.db.Create(&space); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, domain.ErrSpaceAlreadyExists)
	}
	return space.ID, nil
}

func (r *spaceRepository) AddPersonalSpace(ctx context.Context, operator domain.SystemOwner, appUser domain.AppUser) (domain.SpaceID, error) {
	logger := log.FromContext(ctx)
	space := spaceEntity{
		Version:        1,
		CreatedBy:      appUser.GetID(),
		UpdatedBy:      appUser.GetID(),
		OrganizationID: uint(appUser.GetOrganizationID()),
		Type:           SpaceTypePersonal,
		Key:            strconv.Itoa(int(appUser.GetID())),
		Name:           "Default",
		Description:    "",
	}
	logger.Infof("space %+v", space)
	if result := r.db.Create(&space); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, domain.ErrSpaceAlreadyExists)
	}
	return domain.SpaceID(space.ID), nil
}
