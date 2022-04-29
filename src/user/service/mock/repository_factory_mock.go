package service_mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/src/user/service"
)

type RepositoryFactoryMock struct {
	mock.Mock
}

func (m *RepositoryFactoryMock) NewOrganizationRepository() service.OrganizationRepository {
	args := m.Called()
	return args.Get(0).(service.OrganizationRepository)
}
func (m *RepositoryFactoryMock) NewSpaceRepository() service.SpaceRepository {
	args := m.Called()
	return args.Get(0).(service.SpaceRepository)
}
func (m *RepositoryFactoryMock) NewAppUserRepository() service.AppUserRepository {
	args := m.Called()
	return args.Get(0).(service.AppUserRepository)
}
func (m *RepositoryFactoryMock) NewAppUserGroupRepository() service.AppUserGroupRepository {
	args := m.Called()
	return args.Get(0).(service.AppUserGroupRepository)
}
func (m *RepositoryFactoryMock) NewGroupUserRepository() service.GroupUserRepository {
	args := m.Called()
	return args.Get(0).(service.GroupUserRepository)
}
func (m *RepositoryFactoryMock) NewUserSpaceRepository() service.UserSpaceRepository {
	args := m.Called()
	return args.Get(0).(service.UserSpaceRepository)
}
func (m *RepositoryFactoryMock) NewRBACRepository() service.RBACRepository {
	args := m.Called()
	return args.Get(0).(service.RBACRepository)
}
