package domain

import (
	"context"
	"errors"
)

var ErrAppUserNotFound = errors.New("AppUser not found")
var ErrAppUserAlreadyExists = errors.New("AppUser already exists")

var ErrSystemOwnerNotFound = errors.New("SystemOwner not found")

type AppUserAddParameter struct {
	LoginID    string
	Username   string
	Roles      []string
	Properties map[string]string
}

func NewAppUserAddParameter(loginID, username string, roles []string, properties map[string]string) *AppUserAddParameter {
	return &AppUserAddParameter{
		LoginID:    loginID,
		Username:   username,
		Roles:      roles,
		Properties: properties,
	}
}

type AppUserRepository interface {
	FindSystemOwnerByOrganizationID(ctx context.Context, operator SystemAdmin, organizationID OrganizationID) (SystemOwner, error)

	FindSystemOwnerByOrganizationName(ctx context.Context, operator SystemAdmin, organizationName string) (SystemOwner, error)

	FindAppUserByID(ctx context.Context, operator AppUser, id AppUserID) (AppUser, error)

	FindAppUserByLoginID(ctx context.Context, operator AppUser, loginID string) (AppUser, error)

	FindOwnerByLoginID(ctx context.Context, operator SystemOwner, loginID string) (Owner, error)

	AddAppUser(ctx context.Context, operator Owner, param *AppUserAddParameter) (AppUserID, error)

	AddSystemOwner(ctx context.Context, operator SystemAdmin, organizationID OrganizationID) (AppUserID, error)

	AddFirstOwner(ctx context.Context, operator SystemOwner, param *FirstOwnerAddParameter) (AppUserID, error)
}
