// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/kujilabo/cocotola-api/src/app/domain"
	mock "github.com/stretchr/testify/mock"

	testing "testing"

	time "time"
)

// TatoebaSentence is an autogenerated mock type for the TatoebaSentence type
type TatoebaSentence struct {
	mock.Mock
}

// GetAuthor provides a mock function with given fields:
func (_m *TatoebaSentence) GetAuthor() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetLang2 provides a mock function with given fields:
func (_m *TatoebaSentence) GetLang2() domain.Lang2 {
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

// GetSentenceNumber provides a mock function with given fields:
func (_m *TatoebaSentence) GetSentenceNumber() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetText provides a mock function with given fields:
func (_m *TatoebaSentence) GetText() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetUpdatedAt provides a mock function with given fields:
func (_m *TatoebaSentence) GetUpdatedAt() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// NewTatoebaSentence creates a new instance of TatoebaSentence. It also registers a cleanup function to assert the mocks expectations.
func NewTatoebaSentence(t testing.TB) *TatoebaSentence {
	mock := &TatoebaSentence{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
