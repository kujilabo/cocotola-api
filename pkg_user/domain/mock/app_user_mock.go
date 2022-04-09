package domain_mock

import (
	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/stretchr/testify/mock"
)

type AppUserMock struct {
	mock.Mock
}

func (m *AppUserMock) GetID() uint {
	args := m.Called()
	return args.Get(0).(uint)
}

func (m *AppUserMock) GetOrganizationID() domain.OrganizationID {
	args := m.Called()
	return domain.OrganizationID(args.Get(0).(uint))
}

func (m *AppUserMock) GetLoginID() string {
	args := m.Called()
	return args.String(0)
}

func (m *AppUserMock) GetUsername() string {
	args := m.Called()
	return args.String(0)
}

func (m *AppUserMock) GetRoles() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *AppUserMock) GetProperties() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}
