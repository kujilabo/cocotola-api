package gateway

import (
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type repositoryFactory struct {
	db *gorm.DB
}

func NewRepositoryFactory(db *gorm.DB) (domain.RepositoryFactory, error) {
	return &repositoryFactory{
		db: db,
	}, nil
}

func (f *repositoryFactory) NewOrganizationRepository() domain.OrganizationRepository {
	return NewOrganizationRepository(f.db)
}

func (f *repositoryFactory) NewSpaceRepository() domain.SpaceRepository {
	return NewSpaceRepository(f.db)
}

func (f *repositoryFactory) NewAppUserRepository() domain.AppUserRepository {
	return NewAppUserRepository(f, f.db)
}

func (f *repositoryFactory) NewAppUserGroupRepository() domain.AppUserGroupRepository {
	return NewAppUserGroupRepository(f.db)
}

func (f *repositoryFactory) NewGroupUserRepository() domain.GroupUserRepository {
	return NewGroupUserRepository(f.db)
}

func (f *repositoryFactory) NewUserSpaceRepository() domain.UserSpaceRepository {
	return NewUserSpaceRepository(f, f.db)
}

func (f *repositoryFactory) NewRBACRepository() domain.RBACRepository {
	return NewRBACRepository(f.db)
}
