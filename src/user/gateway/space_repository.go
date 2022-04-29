package gateway

import (
	"context"
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"

	libG "github.com/kujilabo/cocotola-api/src/lib/gateway"
	"github.com/kujilabo/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola-api/src/user/service"
)

const SpaceTypeDefault = 1
const SpaceTypePersonal = 2
const SpaceTypeSystem = 3

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

func (e *spaceEntity) toSpace() (service.Space, error) {
	model, err := domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	spaceModel, err := domain.NewSpaceModel(model, domain.OrganizationID(e.OrganizationID), e.Type, e.Key, e.Name, e.Description)
	if err != nil {
		return nil, err
	}

	return service.NewSpace(spaceModel)
}

type spaceRepository struct {
	db *gorm.DB
}

func NewSpaceRepository(db *gorm.DB) service.SpaceRepository {
	return &spaceRepository{
		db: db,
	}
}

func (r *spaceRepository) FindDefaultSpace(ctx context.Context, operator domain.AppUserModel) (service.Space, error) {
	_, span := tracer.Start(ctx, "spaceRepository.FindDefaultSpace")
	defer span.End()

	space := spaceEntity{}
	result := r.db.Where(&spaceEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		Type:           SpaceTypeDefault,
		Key:            "default",
	}).First(&space)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrSpaceNotFound
		}
		return nil, result.Error
	}

	return space.toSpace()
}

func (r *spaceRepository) FindPersonalSpace(ctx context.Context, operator domain.AppUserModel) (service.Space, error) {
	_, span := tracer.Start(ctx, "spaceRepository.FindPersonalSpace")
	defer span.End()

	space := spaceEntity{}
	if result := r.db.Where(&spaceEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		Type:           SpaceTypePersonal,
		Key:            strconv.Itoa(int(operator.GetID())),
	}).First(&space); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrSpaceNotFound
		}
		return nil, result.Error
	}

	return space.toSpace()
}

func (r *spaceRepository) FindSystemSpace(ctx context.Context, operator domain.AppUserModel) (service.Space, error) {
	_, span := tracer.Start(ctx, "spaceRepository.FindSystemSpace")
	defer span.End()

	space := spaceEntity{}
	if result := r.db.Where(&spaceEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		Type:           SpaceTypeSystem,
		Key:            strconv.Itoa(int(operator.GetID())),
	}).First(&space); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrSpaceNotFound
		}
		return nil, result.Error
	}

	return space.toSpace()
}

func (r *spaceRepository) AddDefaultSpace(ctx context.Context, operator domain.SystemOwnerModel) (uint, error) {
	_, span := tracer.Start(ctx, "spaceRepository.AddDefaultSpace")
	defer span.End()

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
		return 0, libG.ConvertDuplicatedError(result.Error, service.ErrSpaceAlreadyExists)
	}
	return space.ID, nil
}

func (r *spaceRepository) AddPersonalSpace(ctx context.Context, operator domain.SystemOwnerModel, appUser domain.AppUserModel) (domain.SpaceID, error) {
	_, span := tracer.Start(ctx, "spaceRepository.AddPersonalSpace")
	defer span.End()

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

	if result := r.db.Create(&space); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, service.ErrSpaceAlreadyExists)
	}
	return domain.SpaceID(space.ID), nil
}

func (r *spaceRepository) AddSystemSpace(ctx context.Context, operator domain.SystemOwnerModel) (domain.SpaceID, error) {
	_, span := tracer.Start(ctx, "spaceRepository.AddSystemSpace")
	defer span.End()

	space := spaceEntity{
		Version:        1,
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		Type:           SpaceTypeSystem,
		Key:            strconv.Itoa(int(operator.GetID())),
		Name:           "System",
		Description:    "",
	}

	if result := r.db.Create(&space); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, service.ErrSpaceAlreadyExists)
	}
	return domain.SpaceID(space.ID), nil
}
