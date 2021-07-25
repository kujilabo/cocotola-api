package domain

type AppUserGroupID uint

type AppUserGroup interface {
	Model
	OrganizationID() OrganizationID
	Key() string
	Name() string
	Description() string
}

type appUserGroup struct {
	Model
	organizationID OrganizationID
	key            string
	name           string
	description    string
}

// NewAppUserGroup returns a new AppUserGroup
func NewAppUserGroup(model Model, organizationID OrganizationID, key, name, description string) AppUserGroup {
	return &appUserGroup{
		Model:          model,
		organizationID: organizationID,
		key:            key,
		name:           name,
		description:    description,
	}
}

func (g *appUserGroup) OrganizationID() OrganizationID {
	return g.organizationID
}

func (g *appUserGroup) Key() string {
	return g.key
}

func (g *appUserGroup) Name() string {
	return g.name
}

func (g *appUserGroup) Description() string {
	return g.description
}
