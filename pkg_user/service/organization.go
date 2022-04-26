package service

import (
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/domain"
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
