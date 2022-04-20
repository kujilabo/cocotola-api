//go:generate mockery --output mock --name SystemStudentModel
package domain

import (
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type SystemStudentModel interface {
	user.AppUserModel
}

type systemStudentModel struct {
	user.AppUserModel
}

func NewSystemStudentModel(appUser user.AppUserModel) (SystemStudentModel, error) {
	m := &systemStudentModel{
		AppUserModel: appUser,
	}

	return m, lib.Validator.Struct(m)
}
