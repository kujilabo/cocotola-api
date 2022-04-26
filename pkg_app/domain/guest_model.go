package domain

import (
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type GuestModel interface {
	userD.AppUserModel
}

type guestModel struct {
	userD.AppUserModel
}

func NewGuestModel(appUser userD.AppUserModel) (GuestModel, error) {
	m := &guestModel{
		AppUserModel: appUser,
	}

	return m, libD.Validator.Struct(m)
}
