//go:generate mockery --output mock --name StudentModel
package domain

import (
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type StudentModel interface {
	userD.AppUserModel
}

type studentModel struct {
	userD.AppUserModel
}

func NewStudentModel(appUserModel userD.AppUserModel) (StudentModel, error) {
	m := &studentModel{
		AppUserModel: appUserModel,
	}

	return m, libD.Validator.Struct(m)
}
