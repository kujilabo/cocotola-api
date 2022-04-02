package domain

import (
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type GuestModel interface {
	user.AppUserModel
}

type guestModel struct {
	user.AppUserModel
}

func NewGuestModel(appUser user.AppUserModel) (GuestModel, error) {
	m := &guestModel{
		AppUserModel: appUser,
	}

	return m, lib.Validator.Struct(m)
}
