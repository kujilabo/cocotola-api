package service_mock

import (
	"github.com/casbin/casbin/v2"
	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type RBACRepositoryyMock struct {
	mock.Mock
}

func (m *RBACRepositoryyMock) Init() error {
	args := m.Called()
	return args.Error(0)
}

func (m *RBACRepositoryyMock) AddNamedPolicy(subject domain.RBACRole, object domain.RBACObject, action domain.RBACAction) error {
	args := m.Called(subject, action)
	return args.Error(0)
}

func (m *RBACRepositoryyMock) AddNamedGroupingPolicy(subject domain.RBACUser, object domain.RBACRole) error {
	args := m.Called(subject, object)
	return args.Error(0)
}

func (m *RBACRepositoryyMock) NewEnforcerWithRolesAndUsers(roles []domain.RBACRole, users []domain.RBACUser) (*casbin.Enforcer, error) {
	args := m.Called(roles, users)
	return args.Get(0).(*casbin.Enforcer), args.Error(1)
}
