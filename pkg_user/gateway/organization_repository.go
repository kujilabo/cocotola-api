package gateway

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/service"
)

type organizationRepository struct {
	db *gorm.DB
}

type organizationEntity struct {
	ID        uint
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uint
	UpdatedBy uint
	Name      string
}

func (e *organizationEntity) TableName() string {
	return "organization"
}

func (e *organizationEntity) toModel() (service.Organization, error) {
	model, err := domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	organizationModel, err := domain.NewOrganizationModel(model, e.Name)
	if err != nil {
		return nil, err
	}

	return service.NewOrganization(organizationModel)
}

func NewOrganizationRepository(db *gorm.DB) service.OrganizationRepository {
	return &organizationRepository{
		db: db,
	}
}

func (r *organizationRepository) GetOrganization(ctx context.Context, operator domain.AppUserModel) (service.Organization, error) {
	_, span := tracer.Start(ctx, "organizationRepository.GetOrganization")
	defer span.End()

	organization := organizationEntity{}

	if result := r.db.Where(organizationEntity{
		ID: uint(operator.GetOrganizationID()),
	}).First(&organization); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrOrganizationNotFound
		}
		return nil, result.Error
	}

	return organization.toModel()
}

func (r *organizationRepository) FindOrganizationByName(ctx context.Context, operator domain.SystemAdminModel, name string) (service.Organization, error) {
	_, span := tracer.Start(ctx, "organizationRepository.FindOrganizationByName")
	defer span.End()

	organization := organizationEntity{}

	if result := r.db.Where(organizationEntity{
		Name: name,
	}).First(&organization); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrOrganizationNotFound
		}
		return nil, result.Error
	}

	return organization.toModel()
}

func (r *organizationRepository) FindOrganizationByID(ctx context.Context, operator domain.SystemAdminModel, id domain.OrganizationID) (service.Organization, error) {
	_, span := tracer.Start(ctx, "organizationRepository.FindOrganizationByID")
	defer span.End()

	organization := organizationEntity{}

	if result := r.db.Where(organizationEntity{
		ID: uint(id),
	}).First(&organization); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrOrganizationNotFound
		}
		return nil, result.Error
	}

	return organization.toModel()
}

func (r *organizationRepository) AddOrganization(ctx context.Context, operator domain.SystemAdminModel, param service.OrganizationAddParameter) (domain.OrganizationID, error) {
	_, span := tracer.Start(ctx, "organizationRepository.AddOrganization")
	defer span.End()

	organization := organizationEntity{
		Version:   1,
		CreatedBy: operator.GetID(),
		UpdatedBy: operator.GetID(),
		Name:      param.GetName(),
	}

	if result := r.db.Create(&organization); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, service.ErrOrganizationAlreadyExists)
	}

	return domain.OrganizationID(organization.ID), nil
}
