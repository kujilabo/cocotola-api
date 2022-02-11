package domain_mock

import (
	"time"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/stretchr/testify/mock"
)

type SpaceMock struct {
	mock.Mock
}

func (m *SpaceMock) GetID() uint {
	args := m.Called()
	return args.Get(0).(uint)
}
func (m *SpaceMock) GetVersion() int {
	args := m.Called()
	return args.Int(0)
}
func (m *SpaceMock) GetCreatedAt() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}
func (m *SpaceMock) GetUpdatedAt() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}
func (m *SpaceMock) GetCreatedBy() uint {
	args := m.Called()
	return args.Get(0).(uint)
}
func (m *SpaceMock) GetUpdatedBy() uint {
	args := m.Called()
	return args.Get(0).(uint)
}

func (m *SpaceMock) GetOrganizationID() domain.OrganizationID {
	args := m.Called()
	return args.Get(0).(domain.OrganizationID)
}

func (m *SpaceMock) GetKey() string {
	args := m.Called()
	return args.String(0)
}

func (m *SpaceMock) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *SpaceMock) GetDescription() string {
	args := m.Called()
	return args.String(0)
}
