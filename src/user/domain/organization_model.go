package domain

import libD "github.com/kujilabo/cocotola-api/src/lib/domain"

type OrganizationID uint

type OrganizationModel interface {
	Model
	GetName() string
}

type organizationModel struct {
	Model
	Name string `validate:"required"`
}

func NewOrganizationModel(model Model, name string) (OrganizationModel, error) {
	m := &organizationModel{
		Model: model,
		Name:  name,
	}
	return m, libD.Validator.Struct(m)
}

func (m *organizationModel) GetName() string {
	return m.Name
}

func (m *organizationModel) String() string {
	return m.Name
}
