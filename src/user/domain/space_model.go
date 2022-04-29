//go:generate mockery --output mock --name SpaceModel
package domain

import libD "github.com/kujilabo/cocotola-api/src/lib/domain"

type SpaceID uint
type SpaceTypeID int

type SpaceModel interface {
	Model
	GetOrganizationID() OrganizationID
	GetKey() string
	GetName() string
	GetDescription() string
}

type spaceModel struct {
	Model
	OrganizationID OrganizationID `validate:"required,gte=1"`
	SpaceType      int            `validate:"required,gte=1"`
	Key            string         `validate:"required"`
	Name           string         `validate:"required"`
	Description    string
}

func NewSpaceModel(model Model, organizationID OrganizationID, spaceType int, key, name, description string) (SpaceModel, error) {
	m := &spaceModel{
		Model:          model,
		OrganizationID: organizationID,
		SpaceType:      spaceType,
		Key:            key,
		Name:           name,
		Description:    description,
	}

	return m, libD.Validator.Struct(m)
}

func (m *spaceModel) GetOrganizationID() OrganizationID {
	return m.OrganizationID
}

func (m *spaceModel) GetKey() string {
	return m.Key
}

func (m *spaceModel) GetName() string {
	return m.Name
}

func (m *spaceModel) GetDescription() string {
	return m.Description
}

func (m *spaceModel) String() string {
	return m.Name
}
