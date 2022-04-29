package service

import (
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	"github.com/kujilabo/cocotola-api/src/user/domain"
)

type AppUserGroup interface {
	domain.AppUserGroupModel
}

type appUserGroup struct {
	domain.AppUserGroupModel
}

// NewAppUserGroup returns a new AppUserGroup
func NewAppUserGroup(appUserGroupModel domain.AppUserGroupModel) (AppUserGroup, error) {
	m := &appUserGroup{
		appUserGroupModel,
	}

	return m, libD.Validator.Struct(m)
}
