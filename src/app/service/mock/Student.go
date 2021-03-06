// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kujilabo/cocotola-api/src/app/domain"
	mock "github.com/stretchr/testify/mock"

	service "github.com/kujilabo/cocotola-api/src/app/service"

	testing "testing"

	userdomain "github.com/kujilabo/cocotola-api/src/user/domain"

	userservice "github.com/kujilabo/cocotola-api/src/user/service"
)

// Student is an autogenerated mock type for the Student type
type Student struct {
	mock.Mock
}

// AddWorkbookToPersonalSpace provides a mock function with given fields: ctx, parameter
func (_m *Student) AddWorkbookToPersonalSpace(ctx context.Context, parameter service.WorkbookAddParameter) (domain.WorkbookID, error) {
	ret := _m.Called(ctx, parameter)

	var r0 domain.WorkbookID
	if rf, ok := ret.Get(0).(func(context.Context, service.WorkbookAddParameter) domain.WorkbookID); ok {
		r0 = rf(ctx, parameter)
	} else {
		r0 = ret.Get(0).(domain.WorkbookID)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, service.WorkbookAddParameter) error); ok {
		r1 = rf(ctx, parameter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckQuota provides a mock function with given fields: ctx, problemType, name
func (_m *Student) CheckQuota(ctx context.Context, problemType string, name service.QuotaName) error {
	ret := _m.Called(ctx, problemType, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, service.QuotaName) error); ok {
		r0 = rf(ctx, problemType, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DecrementQuotaUsage provides a mock function with given fields: ctx, problemType, name, value
func (_m *Student) DecrementQuotaUsage(ctx context.Context, problemType string, name service.QuotaName, value int) error {
	ret := _m.Called(ctx, problemType, name, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, service.QuotaName, int) error); ok {
		r0 = rf(ctx, problemType, name, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindRecordbook provides a mock function with given fields: ctx, workbookID, studyType
func (_m *Student) FindRecordbook(ctx context.Context, workbookID domain.WorkbookID, studyType string) (service.Recordbook, error) {
	ret := _m.Called(ctx, workbookID, studyType)

	var r0 service.Recordbook
	if rf, ok := ret.Get(0).(func(context.Context, domain.WorkbookID, string) service.Recordbook); ok {
		r0 = rf(ctx, workbookID, studyType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.Recordbook)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.WorkbookID, string) error); ok {
		r1 = rf(ctx, workbookID, studyType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindRecordbookSummary provides a mock function with given fields: ctx, workbookID
func (_m *Student) FindRecordbookSummary(ctx context.Context, workbookID domain.WorkbookID) (service.RecordbookSummary, error) {
	ret := _m.Called(ctx, workbookID)

	var r0 service.RecordbookSummary
	if rf, ok := ret.Get(0).(func(context.Context, domain.WorkbookID) service.RecordbookSummary); ok {
		r0 = rf(ctx, workbookID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.RecordbookSummary)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.WorkbookID) error); ok {
		r1 = rf(ctx, workbookID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindWorkbookByID provides a mock function with given fields: ctx, id
func (_m *Student) FindWorkbookByID(ctx context.Context, id domain.WorkbookID) (service.Workbook, error) {
	ret := _m.Called(ctx, id)

	var r0 service.Workbook
	if rf, ok := ret.Get(0).(func(context.Context, domain.WorkbookID) service.Workbook); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.Workbook)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.WorkbookID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindWorkbookByName provides a mock function with given fields: ctx, name
func (_m *Student) FindWorkbookByName(ctx context.Context, name string) (service.Workbook, error) {
	ret := _m.Called(ctx, name)

	var r0 service.Workbook
	if rf, ok := ret.Get(0).(func(context.Context, string) service.Workbook); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.Workbook)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindWorkbooksFromPersonalSpace provides a mock function with given fields: ctx, condition
func (_m *Student) FindWorkbooksFromPersonalSpace(ctx context.Context, condition service.WorkbookSearchCondition) (service.WorkbookSearchResult, error) {
	ret := _m.Called(ctx, condition)

	var r0 service.WorkbookSearchResult
	if rf, ok := ret.Get(0).(func(context.Context, service.WorkbookSearchCondition) service.WorkbookSearchResult); ok {
		r0 = rf(ctx, condition)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.WorkbookSearchResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, service.WorkbookSearchCondition) error); ok {
		r1 = rf(ctx, condition)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDefaultSpace provides a mock function with given fields: ctx
func (_m *Student) GetDefaultSpace(ctx context.Context) (userservice.Space, error) {
	ret := _m.Called(ctx)

	var r0 userservice.Space
	if rf, ok := ret.Get(0).(func(context.Context) userservice.Space); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(userservice.Space)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetID provides a mock function with given fields:
func (_m *Student) GetID() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// GetLoginID provides a mock function with given fields:
func (_m *Student) GetLoginID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetOrganizationID provides a mock function with given fields:
func (_m *Student) GetOrganizationID() userdomain.OrganizationID {
	ret := _m.Called()

	var r0 userdomain.OrganizationID
	if rf, ok := ret.Get(0).(func() userdomain.OrganizationID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(userdomain.OrganizationID)
	}

	return r0
}

// GetPersonalSpace provides a mock function with given fields: ctx
func (_m *Student) GetPersonalSpace(ctx context.Context) (userservice.Space, error) {
	ret := _m.Called(ctx)

	var r0 userservice.Space
	if rf, ok := ret.Get(0).(func(context.Context) userservice.Space); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(userservice.Space)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProperties provides a mock function with given fields:
func (_m *Student) GetProperties() map[string]string {
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

// GetRoles provides a mock function with given fields:
func (_m *Student) GetRoles() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// GetUsername provides a mock function with given fields:
func (_m *Student) GetUsername() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IncrementQuotaUsage provides a mock function with given fields: ctx, problemType, name, value
func (_m *Student) IncrementQuotaUsage(ctx context.Context, problemType string, name service.QuotaName, value int) error {
	ret := _m.Called(ctx, problemType, name, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, service.QuotaName, int) error); ok {
		r0 = rf(ctx, problemType, name, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveWorkbook provides a mock function with given fields: ctx, id, version
func (_m *Student) RemoveWorkbook(ctx context.Context, id domain.WorkbookID, version int) error {
	ret := _m.Called(ctx, id, version)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.WorkbookID, int) error); ok {
		r0 = rf(ctx, id, version)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateWorkbook provides a mock function with given fields: ctx, workbookID, version, parameter
func (_m *Student) UpdateWorkbook(ctx context.Context, workbookID domain.WorkbookID, version int, parameter service.WorkbookUpdateParameter) error {
	ret := _m.Called(ctx, workbookID, version, parameter)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.WorkbookID, int, service.WorkbookUpdateParameter) error); ok {
		r0 = rf(ctx, workbookID, version, parameter)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStudent creates a new instance of Student. It also registers a cleanup function to assert the mocks expectations.
func NewStudent(t testing.TB) *Student {
	mock := &Student{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
