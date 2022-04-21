//go:generate mockery --output mock --name AppUserModel
package domain

import (
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

type AppUserID uint

type AppUserModel interface {
	GetID() uint
	GetOrganizationID() OrganizationID
	GetLoginID() string
	GetUsername() string
	GetRoles() []string
	GetProperties() map[string]string
}

type appUserModel struct {
	Model
	OrganizationID OrganizationID `validate:"required,gte=1"`
	LoginID        string         `validate:"required"`
	Username       string         `validate:"required"`
	Roles          []string
	Properties     map[string]string
}

func NewAppUserModel(model Model, organizationID OrganizationID, loginID, username string, roles []string, properties map[string]string) (AppUserModel, error) {
	m := &appUserModel{
		Model:          model,
		OrganizationID: organizationID,
		LoginID:        loginID,
		Username:       username,
		Roles:          roles,
		Properties:     properties,
	}

	return m, lib.Validator.Struct(m)
}

func (a *appUserModel) GetOrganizationID() OrganizationID {
	return a.OrganizationID
}

func (a *appUserModel) GetLoginID() string {
	return a.LoginID
}

func (a *appUserModel) GetUsername() string {
	return a.Username
}

func (a *appUserModel) GetRoles() []string {
	return a.Roles
}

func (a *appUserModel) GetProperties() map[string]string {
	return a.Properties
}
