package domain

type SpaceID uint

type Space interface {
	Model
	GetOrganizationID() OrganizationID
	GetKey() string
	GetName() string
	GetDescription() string
}

type space struct {
	Model
	OrganizationID OrganizationID
	SpaceType      int
	Key            string `validate:"required"`
	Name           string `validate:"required"`
	Description    string `validate:"required"`
}

func NewSpace(model Model, organizationID OrganizationID, spaceType int, key, name, description string) Space {
	return &space{
		Model:          model,
		OrganizationID: organizationID,
		SpaceType:      spaceType,
		Key:            key,
		Name:           name,
		Description:    description,
	}
}

func (m *space) GetOrganizationID() OrganizationID {
	return m.OrganizationID
}

func (m *space) GetKey() string {
	return m.Key
}

func (m *space) GetName() string {
	return m.Name
}

func (m *space) GetDescription() string {
	return m.Description
}

func (m *space) String() string {
	return m.Name
}
