package service

import (
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type Space interface {
	domain.SpaceModel
}

type space struct {
	domain.SpaceModel
}

func NewSpace(spaceModel domain.SpaceModel) (Space, error) {
	m := &space{
		spaceModel,
	}

	return m, lib.Validator.Struct(m)
}