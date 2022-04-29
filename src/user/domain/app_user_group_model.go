package domain

import libD "github.com/kujilabo/cocotola-api/src/lib/domain"

type AppUserGroupID uint

type AppUserGroupModel interface {
	Model
	GetOrganizationID() OrganizationID
	GetKey() string
	GetName() string
	GetDescription() string
}

type appUserGroupModel struct {
	Model
	OrganizationID OrganizationID
	Key            string `validate:"required"`
	Name           string `validate:"required"`
	Description    string
}

// NewAppUserGroup returns a new AppUserGroup
func NewAppUserGroup(model Model, organizationID OrganizationID, key, name, description string) (AppUserGroupModel, error) {
	m := &appUserGroupModel{
		Model:          model,
		OrganizationID: organizationID,
		Key:            key,
		Name:           name,
		Description:    description,
	}

	return m, libD.Validator.Struct(m)
}

func (g *appUserGroupModel) GetOrganizationID() OrganizationID {
	return g.OrganizationID
}

func (g *appUserGroupModel) GetKey() string {
	return g.Key
}

func (g *appUserGroupModel) GetName() string {
	return g.Name
}

func (g *appUserGroupModel) GetDescription() string {
	return g.Description
}
