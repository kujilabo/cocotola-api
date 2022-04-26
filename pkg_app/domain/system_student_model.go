//go:generate mockery --output mock --name SystemStudentModel
package domain

import (
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type SystemStudentModel interface {
	userD.AppUserModel
}

type systemStudentModel struct {
	userD.AppUserModel
}

func NewSystemStudentModel(appUser userD.AppUserModel) (SystemStudentModel, error) {
	m := &systemStudentModel{
		AppUserModel: appUser,
	}

	return m, libD.Validator.Struct(m)
}
