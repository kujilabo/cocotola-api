package domain

import (
	"context"
	"errors"

	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

var ErrAppUserNotFound = errors.New("AppUser not found")
var ErrAppUserAlreadyExists = errors.New("AppUser already exists")

var ErrSystemOwnerNotFound = errors.New("SystemOwner not found")

type AppUserAddParameter interface {
	GetLoginID() string
	GetUsername() string
	GetRoles() []string
	GetProperties() map[string]string
}

type appUserAddParameter struct {
	LoginID    string
	Username   string
	Roles      []string
	Properties map[string]string
}

func NewAppUserAddParameter(loginID, username string, roles []string, properties map[string]string) (AppUserAddParameter, error) {
	m := &appUserAddParameter{
		LoginID:    loginID,
		Username:   username,
		Roles:      roles,
		Properties: properties,
	}
	return m, lib.Validator.Struct(m)
}

func (p *appUserAddParameter) GetLoginID() string {
	return p.LoginID
}
func (p *appUserAddParameter) GetUsername() string {
	return p.Username
}
func (p *appUserAddParameter) GetRoles() []string {
	return p.Roles
}
func (p *appUserAddParameter) GetProperties() map[string]string {
	return p.Properties
}

type AppUserRepository interface {
	FindSystemOwnerByOrganizationID(ctx context.Context, operator SystemAdmin, organizationID OrganizationID) (SystemOwner, error)

	FindSystemOwnerByOrganizationName(ctx context.Context, operator SystemAdmin, organizationName string) (SystemOwner, error)

	FindAppUserByID(ctx context.Context, operator AppUser, id AppUserID) (AppUser, error)

	FindAppUserByLoginID(ctx context.Context, operator AppUser, loginID string) (AppUser, error)

	FindOwnerByLoginID(ctx context.Context, operator SystemOwner, loginID string) (Owner, error)

	AddAppUser(ctx context.Context, operator Owner, param AppUserAddParameter) (AppUserID, error)

	AddSystemOwner(ctx context.Context, operator SystemAdmin, organizationID OrganizationID) (AppUserID, error)

	AddFirstOwner(ctx context.Context, operator SystemOwner, param FirstOwnerAddParameter) (AppUserID, error)
}
