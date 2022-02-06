package domain

import (
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type Guest interface {
	user.AppUser
}

type guest struct {
	user.AppUser
}

func NewGuest(appUser user.AppUser) (Guest, error) {
	m := &guest{
		AppUser: appUser,
	}

	return m, lib.Validator.Struct(m)
}
