package gateway

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
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

func (e *organizationEntity) toModel() (domain.Organization, error) {
	return domain.NewOrganization(domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy), e.Name)
}

func NewOrganizationRepository(db *gorm.DB) domain.OrganizationRepository {
	return &organizationRepository{
		db: db,
	}
}

func (r *organizationRepository) GetOrganization(ctx context.Context, operator domain.AppUser) (domain.Organization, error) {
	organization := organizationEntity{}

	if result := r.db.Where(organizationEntity{
		ID: uint(operator.GetOrganizationID()),
	}).First(&organization); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOrganizationNotFound
		}
		return nil, result.Error
	}

	return organization.toModel()
}

func (r *organizationRepository) FindOrganizationByName(ctx context.Context, operator domain.SystemAdmin, name string) (domain.Organization, error) {
	organization := organizationEntity{}

	if result := r.db.Where(organizationEntity{
		Name: name,
	}).First(&organization); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOrganizationNotFound
		}
		return nil, result.Error
	}

	return organization.toModel()
}

func (r *organizationRepository) AddOrganization(ctx context.Context, operator domain.SystemAdmin, param domain.OrganizationAddParameter) (domain.OrganizationID, error) {
	organization := organizationEntity{
		CreatedBy: operator.GetID(),
		UpdatedBy: operator.GetID(),
		Name:      param.GetName(),
	}

	if result := r.db.Create(&organization); result.Error != nil {
		return 0, convertDuplicatedError(result.Error, domain.ErrOrganizationAlreadyExists)
	}

	return domain.OrganizationID(organization.ID), nil
}
