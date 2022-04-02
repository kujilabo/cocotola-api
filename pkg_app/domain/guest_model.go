package domain

import (
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type Guest interface {
	user.AppUserModel
}

type guest struct {
	user.AppUserModel
}

func NewGuest(appUser user.AppUserModel) (Guest, error) {
	m := &guest{
		AppUserModel: appUser,
	}

	return m, lib.Validator.Struct(m)
}
