package service

import (
	"context"
	"errors"

	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

var ErrOrganizationNotFound = errors.New("organization not found")
var ErrOrganizationAlreadyExists = errors.New("organization already exists")

type FirstOwnerAddParameter interface {
	GetLoginID() string
	GetUsername() string
	GetPassword() string
}

type firstOwnerAddParameter struct {
	LoginID  string `validate:"required"`
	Username string `validate:"required"`
	Password string
}

func NewFirstOwnerAddParameter(loginID, username, password string) (FirstOwnerAddParameter, error) {
	m := &firstOwnerAddParameter{
		LoginID:  loginID,
		Username: username,
		Password: password,
	}
	return m, libD.Validator.Struct(m)
}

func (p *firstOwnerAddParameter) GetLoginID() string {
	return p.LoginID
}
func (p *firstOwnerAddParameter) GetUsername() string {
	return p.Username
}
func (p *firstOwnerAddParameter) GetPassword() string {
	return p.Password
}

type OrganizationAddParameter interface {
	GetName() string
	GetFirstOwner() FirstOwnerAddParameter
}

type organizationAddParameter struct {
	Name       string `validate:"required"`
	FirstOwner FirstOwnerAddParameter
}

func NewOrganizationAddParameter(name string, firstOwner FirstOwnerAddParameter) (OrganizationAddParameter, error) {
	m := &organizationAddParameter{
		Name:       name,
		FirstOwner: firstOwner,
	}
	return m, libD.Validator.Struct(m)
}

func (p *organizationAddParameter) GetName() string {
	return p.Name
}
func (p *organizationAddParameter) GetFirstOwner() FirstOwnerAddParameter {
	return p.FirstOwner
}

type OrganizationRepository interface {
	GetOrganization(ctx context.Context, operator domain.AppUserModel) (Organization, error)

	FindOrganizationByName(ctx context.Context, operator domain.SystemAdminModel, name string) (Organization, error)

	FindOrganizationByID(ctx context.Context, operator domain.SystemAdminModel, id domain.OrganizationID) (Organization, error)

	AddOrganization(ctx context.Context, operator domain.SystemAdminModel, param OrganizationAddParameter) (domain.OrganizationID, error)

	// FindOrganizationByName(ctx context.Context, operator SystemAdmin, name string) (Organization, error)
	// FindOrganization(ctx context.Context, operator AppUser) (Organization, error)
}
