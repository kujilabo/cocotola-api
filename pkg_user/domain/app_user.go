package domain

import (
	"context"

	"github.com/go-playground/validator"
)

type AppUserID uint

type AppUser interface {
	GetID() uint
	GetOrganizationID() OrganizationID
	GetLoginID() string
	GetUsername() string
	GetRoles() []string
	GetProperties() map[string]string

	GetDefaultSpace(ctx context.Context) (Space, error)
	GetPersonalSpace(ctx context.Context) (Space, error)
}

type appUser struct {
	rf RepositoryFactory
	Model
	OrganizationID OrganizationID `validate:"required,gte=1"`
	LoginID        string         `validate:"required"`
	Username       string         `validate:"required"`
	Roles          []string
	Properties     map[string]string
}

func NewAppUser(rf RepositoryFactory, model Model, organizationID OrganizationID, loginID, username string, roles []string, properties map[string]string) (AppUser, error) {
	m := &appUser{
		rf:             rf,
		Model:          model,
		OrganizationID: organizationID,
		LoginID:        loginID,
		Username:       username,
		Roles:          roles,
		Properties:     properties,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (a *appUser) GetOrganizationID() OrganizationID {
	return a.OrganizationID
}

func (a *appUser) GetLoginID() string {
	return a.LoginID
}

func (a *appUser) GetUsername() string {
	return a.Username
}

func (a *appUser) GetRoles() []string {
	return a.Roles
}

func (a *appUser) GetProperties() map[string]string {
	return a.Properties
}

func (a *appUser) GetDefaultSpace(ctx context.Context) (Space, error) {
	return a.rf.NewSpaceRepository().FindDefaultSpace(ctx, a)
}

func (a *appUser) GetPersonalSpace(ctx context.Context) (Space, error) {
	return a.rf.NewSpaceRepository().FindPersonalSpace(ctx, a)
}
