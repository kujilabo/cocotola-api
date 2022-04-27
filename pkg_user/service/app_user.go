package service

import (
	"context"

	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

// type AppUserID uint

type AppUser interface {
	domain.AppUserModel

	GetDefaultSpace(ctx context.Context) (Space, error)
	GetPersonalSpace(ctx context.Context) (Space, error)
}

type appUser struct {
	rf RepositoryFactory
	domain.AppUserModel
}

func NewAppUser(rf RepositoryFactory, appUserModel domain.AppUserModel) (AppUser, error) {
	m := &appUser{
		rf:           rf,
		AppUserModel: appUserModel,
	}

	return m, libD.Validator.Struct(m)
}

func (a *appUser) GetDefaultSpace(ctx context.Context) (Space, error) {
	return a.rf.NewSpaceRepository().FindDefaultSpace(ctx, a)
}

func (a *appUser) GetPersonalSpace(ctx context.Context) (Space, error) {
	return a.rf.NewSpaceRepository().FindPersonalSpace(ctx, a)
}
