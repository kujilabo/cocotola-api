package domain_mock

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type ModelMock struct {
	mock.Mock
}

func (m *ModelMock) GetID() uint {
	args := m.Called()
	return args.Get(0).(uint)
}
func (m *ModelMock) GetVersion() int {
	args := m.Called()
	return args.Int(0)
}
func (m *ModelMock) GetCreatedAt() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}
func (m *ModelMock) GetUpdatedAt() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}
func (m *ModelMock) GetCreatedBy() uint {
	args := m.Called()
	return args.Get(0).(uint)
}
func (m *ModelMock) GetUpdatedBy() uint {
	args := m.Called()
	return args.Get(0).(uint)
}
