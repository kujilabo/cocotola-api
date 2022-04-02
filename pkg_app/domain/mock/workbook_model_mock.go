package domain_mock

import (
	"time"

	"github.com/stretchr/testify/mock"

	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type WorkbookModelMock struct {
	mock.Mock
}

func (m *WorkbookModelMock) GetID() uint {
	args := m.Called()
	return args.Get(0).(uint)
}
func (m *WorkbookModelMock) GetVersion() int {
	args := m.Called()
	return args.Int(0)
}
func (m *WorkbookModelMock) GetCreatedAt() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}
func (m *WorkbookModelMock) GetUpdatedAt() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}
func (m *WorkbookModelMock) GetCreatedBy() uint {
	args := m.Called()
	return args.Get(0).(uint)
}
func (m *WorkbookModelMock) GetUpdatedBy() uint {
	args := m.Called()
	return args.Get(0).(uint)
}

func (m *WorkbookModelMock) GetSpaceID() user.SpaceID {
	args := m.Called()
	return args.Get(0).(user.SpaceID)
}

func (m *WorkbookModelMock) GetOwnerID() user.AppUserID {
	args := m.Called()
	return args.Get(0).(user.AppUserID)
}

func (m *WorkbookModelMock) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *WorkbookModelMock) GetProblemType() string {
	args := m.Called()
	return args.String(0)
}

func (m *WorkbookModelMock) GetQuestionText() string {
	args := m.Called()
	return args.String(0)
}

func (m *WorkbookModelMock) GetProperties() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m *WorkbookModelMock) HasPrivilege(privilege user.RBACAction) bool {
	args := m.Called(privilege)
	return args.Bool(0)
}
