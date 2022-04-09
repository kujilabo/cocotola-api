package service

import (
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/domain"
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

	return m, lib.Validator.Struct(m)
}
