package domain

type AppUserGroupID uint

type AppUserGroup interface {
	Model
	GetOrganizationID() OrganizationID
	GetKey() string
	GetName() string
	GetDescription() string
}

type appUserGroup struct {
	Model
	OrganizationID OrganizationID
	Key            string `validate:"required"`
	Name           string `validate:"required"`
	Description    string
}

// NewAppUserGroup returns a new AppUserGroup
func NewAppUserGroup(model Model, organizationID OrganizationID, key, name, description string) AppUserGroup {
	return &appUserGroup{
		Model:          model,
		OrganizationID: organizationID,
		Key:            key,
		Name:           name,
		Description:    description,
	}
}

func (g *appUserGroup) GetOrganizationID() OrganizationID {
	return g.OrganizationID
}

func (g *appUserGroup) GetKey() string {
	return g.Key
}

func (g *appUserGroup) GetName() string {
	return g.Name
}

func (g *appUserGroup) GetDescription() string {
	return g.Description
}
