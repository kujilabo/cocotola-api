package service

import (
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	"github.com/kujilabo/cocotola-api/src/user/domain"
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

	return m, libD.Validator.Struct(m)
}
