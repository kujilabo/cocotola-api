package domain

import lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"

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
	return m, lib.Validator.Struct(m)
}

func (m *organizationModel) GetName() string {
	return m.Name
}

func (m *organizationModel) String() string {
	return m.Name
}
