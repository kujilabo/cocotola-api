// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/kujilabo/cocotola-api/pkg_app/domain"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// ProblemAddParameter is an autogenerated mock type for the ProblemAddParameter type
type ProblemAddParameter struct {
	mock.Mock
}

// GetNumber provides a mock function with given fields:
func (_m *ProblemAddParameter) GetNumber() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetProperties provides a mock function with given fields:
func (_m *ProblemAddParameter) GetProperties() map[string]string {
	ret := _m.Called()

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func() map[string]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	return r0
}

// GetWorkbookID provides a mock function with given fields:
func (_m *ProblemAddParameter) GetWorkbookID() domain.WorkbookID {
	ret := _m.Called()

	var r0 domain.WorkbookID
	if rf, ok := ret.Get(0).(func() domain.WorkbookID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.WorkbookID)
	}

	return r0
}

// NewProblemAddParameter creates a new instance of ProblemAddParameter. It also registers a cleanup function to assert the mocks expectations.
func NewProblemAddParameter(t testing.TB) *ProblemAddParameter {
	mock := &ProblemAddParameter{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}