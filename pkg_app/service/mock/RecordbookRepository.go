// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kujilabo/cocotola-api/pkg_app/domain"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// RecordbookRepository is an autogenerated mock type for the RecordbookRepository type
type RecordbookRepository struct {
	mock.Mock
}

// CountMemorizedProblem provides a mock function with given fields: ctx, operator, workbookID
func (_m *RecordbookRepository) CountMemorizedProblem(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID) (map[string]int, error) {
	ret := _m.Called(ctx, operator, workbookID)

	var r0 map[string]int
	if rf, ok := ret.Get(0).(func(context.Context, domain.StudentModel, domain.WorkbookID) map[string]int); ok {
		r0 = rf(ctx, operator, workbookID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.StudentModel, domain.WorkbookID) error); ok {
		r1 = rf(ctx, operator, workbookID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindStudyRecords provides a mock function with given fields: ctx, operator, workbookID, studyType
func (_m *RecordbookRepository) FindStudyRecords(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyType string) (map[domain.ProblemID]domain.StudyRecord, error) {
	ret := _m.Called(ctx, operator, workbookID, studyType)

	var r0 map[domain.ProblemID]domain.StudyRecord
	if rf, ok := ret.Get(0).(func(context.Context, domain.StudentModel, domain.WorkbookID, string) map[domain.ProblemID]domain.StudyRecord); ok {
		r0 = rf(ctx, operator, workbookID, studyType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[domain.ProblemID]domain.StudyRecord)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.StudentModel, domain.WorkbookID, string) error); ok {
		r1 = rf(ctx, operator, workbookID, studyType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetResult provides a mock function with given fields: ctx, operator, workbookID, studyType, problemType, problemID, studyResult, memorized
func (_m *RecordbookRepository) SetResult(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyType string, problemType string, problemID domain.ProblemID, studyResult bool, memorized bool) error {
	ret := _m.Called(ctx, operator, workbookID, studyType, problemType, problemID, studyResult, memorized)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.StudentModel, domain.WorkbookID, string, string, domain.ProblemID, bool, bool) error); ok {
		r0 = rf(ctx, operator, workbookID, studyType, problemType, problemID, studyResult, memorized)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRecordbookRepository creates a new instance of RecordbookRepository. It also registers a cleanup function to assert the mocks expectations.
func NewRecordbookRepository(t testing.TB) *RecordbookRepository {
	mock := &RecordbookRepository{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
