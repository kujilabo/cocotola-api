package gateway

import (
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/src/user/service"
)

type repositoryFactory struct {
	db *gorm.DB
}

func NewRepositoryFactory(db *gorm.DB) (service.RepositoryFactory, error) {
	return &repositoryFactory{
		db: db,
	}, nil
}

func (f *repositoryFactory) NewOrganizationRepository() service.OrganizationRepository {
	return NewOrganizationRepository(f.db)
}

func (f *repositoryFactory) NewSpaceRepository() service.SpaceRepository {
	return NewSpaceRepository(f.db)
}

func (f *repositoryFactory) NewAppUserRepository() service.AppUserRepository {
	return NewAppUserRepository(f, f.db)
}

func (f *repositoryFactory) NewAppUserGroupRepository() service.AppUserGroupRepository {
	return NewAppUserGroupRepository(f.db)
}

func (f *repositoryFactory) NewGroupUserRepository() service.GroupUserRepository {
	return NewGroupUserRepository(f.db)
}

func (f *repositoryFactory) NewUserSpaceRepository() service.UserSpaceRepository {
	return NewUserSpaceRepository(f, f.db)
}

func (f *repositoryFactory) NewRBACRepository() service.RBACRepository {
	return NewRBACRepository(f.db)
}
