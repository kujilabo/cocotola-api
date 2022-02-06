package domain

import lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"

type OrganizationID uint

type Organization interface {
	Model
	GetName() string
}

type organization struct {
	Model
	Name string `validate:"required"`
}

func NewOrganization(model Model, name string) (Organization, error) {
	m := &organization{
		Model: model,
		Name:  name,
	}
	return m, lib.Validator.Struct(m)
}

func (m *organization) GetName() string {
	return m.Name
}

func (m *organization) String() string {
	return m.Name
}
