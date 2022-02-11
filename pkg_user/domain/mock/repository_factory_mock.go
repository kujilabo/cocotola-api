package domain_mock

import (
	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/stretchr/testify/mock"
)

type RepositoryFactoryMock struct {
	mock.Mock
}

func (m *RepositoryFactoryMock) NewOrganizationRepository() domain.OrganizationRepository {
	args := m.Called()
	return args.Get(0).(domain.OrganizationRepository)
}
func (m *RepositoryFactoryMock) NewSpaceRepository() domain.SpaceRepository {
	args := m.Called()
	return args.Get(0).(domain.SpaceRepository)
}
func (m *RepositoryFactoryMock) NewAppUserRepository() domain.AppUserRepository {
	args := m.Called()
	return args.Get(0).(domain.AppUserRepository)
}
func (m *RepositoryFactoryMock) NewAppUserGroupRepository() domain.AppUserGroupRepository {
	args := m.Called()
	return args.Get(0).(domain.AppUserGroupRepository)
}
func (m *RepositoryFactoryMock) NewGroupUserRepository() domain.GroupUserRepository {
	args := m.Called()
	return args.Get(0).(domain.GroupUserRepository)
}
func (m *RepositoryFactoryMock) NewUserSpaceRepository() domain.UserSpaceRepository {
	args := m.Called()
	return args.Get(0).(domain.UserSpaceRepository)
}
func (m *RepositoryFactoryMock) NewRBACRepository() domain.RBACRepository {
	args := m.Called()
	return args.Get(0).(domain.RBACRepository)
}
