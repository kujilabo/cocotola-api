// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	testing "testing"

	time "time"
)

// Problem is an autogenerated mock type for the Problem type
type Problem struct {
	mock.Mock
}

// GetCreatedAt provides a mock function with given fields:
func (_m *Problem) GetCreatedAt() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// GetCreatedBy provides a mock function with given fields:
func (_m *Problem) GetCreatedBy() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// GetID provides a mock function with given fields:
func (_m *Problem) GetID() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// GetNumber provides a mock function with given fields:
func (_m *Problem) GetNumber() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetProblemType provides a mock function with given fields:
func (_m *Problem) GetProblemType() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetProperties provides a mock function with given fields: ctx
func (_m *Problem) GetProperties(ctx context.Context) map[string]interface{} {
	ret := _m.Called(ctx)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(context.Context) map[string]interface{}); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// GetUpdatedAt provides a mock function with given fields:
func (_m *Problem) GetUpdatedAt() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// GetUpdatedBy provides a mock function with given fields:
func (_m *Problem) GetUpdatedBy() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// GetVersion provides a mock function with given fields:
func (_m *Problem) GetVersion() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// NewProblem creates a new instance of Problem. It also registers a cleanup function to assert the mocks expectations.
func NewProblem(t testing.TB) *Problem {
	mock := &Problem{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
