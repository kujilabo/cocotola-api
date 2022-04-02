package domain

import (
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

const SystemOwnerID = 2

type SystemOwnerModel interface {
	AppUserModel
}

type systemOwnerModel struct {
	AppUserModel
}

func NewSystemOwnerModel(appUser AppUserModel) (SystemOwnerModel, error) {
	m := &systemOwnerModel{
		AppUserModel: appUser,
	}

	return m, lib.Validator.Struct(m)
}
