package service_mock

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

const (
	two = 2
)

type WorkbookMock struct {
	mock.Mock
}

// workbook model mock
func (m *WorkbookMock) GetID() uint {
	args := m.Called()
	return args.Get(0).(uint)
}
func (m *WorkbookMock) GetVersion() int {
	args := m.Called()
	return args.Int(0)
}
func (m *WorkbookMock) GetCreatedAt() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}
func (m *WorkbookMock) GetUpdatedAt() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}
func (m *WorkbookMock) GetCreatedBy() uint {
	args := m.Called()
	return args.Get(0).(uint)
}
func (m *WorkbookMock) GetUpdatedBy() uint {
	args := m.Called()
	return args.Get(0).(uint)
}

func (m *WorkbookMock) GetSpaceID() user.SpaceID {
	args := m.Called()
	return args.Get(0).(user.SpaceID)
}

func (m *WorkbookMock) GetOwnerID() user.AppUserID {
	args := m.Called()
	return args.Get(0).(user.AppUserID)
}

func (m *WorkbookMock) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *WorkbookMock) GetProblemType() string {
	args := m.Called()
	return args.String(0)
}

func (m *WorkbookMock) GetQuestionText() string {
	args := m.Called()
	return args.String(0)
}

func (m *WorkbookMock) GetProperties() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m *WorkbookMock) HasPrivilege(privilege user.RBACAction) bool {
	args := m.Called(privilege)
	return args.Bool(0)
}

// workbook mock

func (m *WorkbookMock) FindProblems(ctx context.Context, operator domain.StudentModel, param service.ProblemSearchCondition) (service.ProblemSearchResult, error) {
	args := m.Called(ctx, operator, param)
	return args.Get(0).(service.ProblemSearchResult), args.Error(1)
}

func (m *WorkbookMock) FindAllProblems(ctx context.Context, operator domain.StudentModel) (service.ProblemSearchResult, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(service.ProblemSearchResult), args.Error(1)
}

func (m *WorkbookMock) FindProblemsByProblemIDs(ctx context.Context, operator domain.StudentModel, param service.ProblemIDsCondition) (service.ProblemSearchResult, error) {
	args := m.Called(ctx, operator, param)
	return args.Get(0).(service.ProblemSearchResult), args.Error(1)
}

func (m *WorkbookMock) FindProblemIDs(ctx context.Context, operator domain.StudentModel) ([]domain.ProblemID, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).([]domain.ProblemID), args.Error(1)
}

func (m *WorkbookMock) FindProblemByID(ctx context.Context, operator domain.StudentModel, problemID domain.ProblemID) (service.Problem, error) {
	args := m.Called(ctx, operator, problemID)
	return args.Get(0).(service.Problem), args.Error(1)
}

func (m *WorkbookMock) AddProblem(ctx context.Context, operator domain.StudentModel, param service.ProblemAddParameter) (service.Added, domain.ProblemID, error) {
	args := m.Called(ctx, operator, param)
	return args.Get(0).(service.Added), args.Get(1).(domain.ProblemID), args.Error(two)
}

func (m *WorkbookMock) UpdateProblem(ctx context.Context, operator domain.StudentModel, id service.ProblemSelectParameter2, param service.ProblemUpdateParameter) (service.Added, service.Updated, error) {
	args := m.Called(ctx, operator, id, param)
	return args.Get(0).(service.Added), args.Get(0).(service.Updated), args.Error(two)
}

func (m *WorkbookMock) RemoveProblem(ctx context.Context, operator domain.StudentModel, id service.ProblemSelectParameter2) error {
	args := m.Called(ctx, operator, id)
	return args.Error(0)
}

func (m *WorkbookMock) UpdateWorkbook(ctx context.Context, operator domain.StudentModel, version int, parameter service.WorkbookUpdateParameter) error {
	args := m.Called(ctx, operator, version, parameter)
	return args.Error(0)
}

func (m *WorkbookMock) RemoveWorkbook(ctx context.Context, operator domain.StudentModel, version int) error {
	args := m.Called(ctx, operator, version)
	return args.Error(0)
}
