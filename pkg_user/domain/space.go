package domain

import lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"

type SpaceID uint
type SpaceTypeID int

type Space interface {
	Model
	GetOrganizationID() OrganizationID
	GetKey() string
	GetName() string
	GetDescription() string
}

type space struct {
	Model
	OrganizationID OrganizationID `validate:"required,gte=1"`
	SpaceType      int            `validate:"required,gte=1"`
	Key            string         `validate:"required"`
	Name           string         `validate:"required"`
	Description    string
}

func NewSpace(model Model, organizationID OrganizationID, spaceType int, key, name, description string) (Space, error) {
	m := &space{
		Model:          model,
		OrganizationID: organizationID,
		SpaceType:      spaceType,
		Key:            key,
		Name:           name,
		Description:    description,
	}

	return m, lib.Validator.Struct(m)
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
