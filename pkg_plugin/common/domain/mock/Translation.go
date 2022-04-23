// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/kujilabo/cocotola-api/pkg_app/domain"
	commondomain "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Translation is an autogenerated mock type for the Translation type
type Translation struct {
	mock.Mock
}

// GetLang provides a mock function with given fields:
func (_m *Translation) GetLang() domain.Lang2 {
	ret := _m.Called()

	var r0 domain.Lang2
	if rf, ok := ret.Get(0).(func() domain.Lang2); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Lang2)
		}
	}

	return r0
}

// GetPos provides a mock function with given fields:
func (_m *Translation) GetPos() commondomain.WordPos {
	ret := _m.Called()

	var r0 commondomain.WordPos
	if rf, ok := ret.Get(0).(func() commondomain.WordPos); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(commondomain.WordPos)
	}

	return r0
}

// GetProvider provides a mock function with given fields:
func (_m *Translation) GetProvider() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetText provides a mock function with given fields:
func (_m *Translation) GetText() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetTranslated provides a mock function with given fields:
func (_m *Translation) GetTranslated() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewTranslation creates a new instance of Translation. It also registers a cleanup function to assert the mocks expectations.
func NewTranslation(t testing.TB) *Translation {
	mock := &Translation{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}