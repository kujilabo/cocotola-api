package domain

import (
	"github.com/go-playground/validator"

	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type Guest interface {
	user.AppUser
}

type guest struct {
	user.AppUser
}

func NewGuest(appUser user.AppUser) (Guest, error) {
	m := &guest{
		AppUser: appUser,
	}

	v := validator.New()
	return m, v.Struct(m)
}
