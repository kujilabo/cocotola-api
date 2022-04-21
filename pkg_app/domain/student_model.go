//go:generate mockery --output mock --name StudentModel
package domain

import (
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type StudentModel interface {
	user.AppUserModel
}

type studentModel struct {
	user.AppUserModel
}

func NewStudentModel(appUserModel user.AppUserModel) (StudentModel, error) {
	m := &studentModel{
		AppUserModel: appUserModel,
	}

	return m, lib.Validator.Struct(m)
}
