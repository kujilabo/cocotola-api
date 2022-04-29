package service

import (
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	"github.com/kujilabo/cocotola-api/src/user/domain"
)

type Organization interface {
	domain.OrganizationModel
}

type organization struct {
	domain.OrganizationModel
}

func NewOrganization(organizationModel domain.OrganizationModel) (Organization, error) {
	m := &organization{
		organizationModel,
	}
	return m, libD.Validator.Struct(m)
}
