package domain

import (
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
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

	return m, libD.Validator.Struct(m)
}
