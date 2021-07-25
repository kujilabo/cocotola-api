package domain

import (
	"context"

	"golang.org/x/xerrors"
)

var ErrOrganizationNotFound = xerrors.New("Organization not found")

type FirstOwnerAddParameter struct {
	LoginID  string
	Password string
	Username string
}

type OrganizationAddParameter struct {
	Name       string
	FirstOwner *FirstOwnerAddParameter
}

func NewOrganizationAddParameter(name string, firstOwner *FirstOwnerAddParameter) *OrganizationAddParameter {
	return &OrganizationAddParameter{
		Name:       name,
		FirstOwner: firstOwner,
	}
}

type OrganizationRepository interface {
	GetOrganization(ctx context.Context, operator AppUser) (Organization, error)

	FindOrganizationByName(ctx context.Context, operator SystemAdmin, name string) (Organization, error)

	AddOrganization(ctx context.Context, operator SystemAdmin, param *OrganizationAddParameter) (OrganizationID, error)

	// FindOrganizationByID(ctx context.Context, operator AppUser, id uint) (Organization, error)
	// FindOrganizationByName(ctx context.Context, operator SystemAdmin, name string) (Organization, error)
	// FindOrganization(ctx context.Context, operator AppUser) (Organization, error)
}
