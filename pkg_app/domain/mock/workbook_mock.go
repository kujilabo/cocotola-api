package domain_mock

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

const (
	two = 2
)

type WorkbookMock struct {
	mock.Mock
}

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

func (m *WorkbookMock) FindProblems(ctx context.Context, operator domain.Student, param domain.ProblemSearchCondition) (domain.ProblemSearchResult, error) {
	args := m.Called(ctx, operator, param)
	return args.Get(0).(domain.ProblemSearchResult), args.Error(1)
}

func (m *WorkbookMock) FindAllProblems(ctx context.Context, operator domain.Student) (domain.ProblemSearchResult, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(domain.ProblemSearchResult), args.Error(1)
}

func (m *WorkbookMock) FindProblemsByProblemIDs(ctx context.Context, operator domain.Student, param domain.ProblemIDsCondition) (domain.ProblemSearchResult, error) {
	args := m.Called(ctx, operator, param)
	return args.Get(0).(domain.ProblemSearchResult), args.Error(1)
}

func (m *WorkbookMock) FindProblemIDs(ctx context.Context, operator domain.Student) ([]domain.ProblemID, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).([]domain.ProblemID), args.Error(1)
}

func (m *WorkbookMock) FindProblemByID(ctx context.Context, operator domain.Student, problemID domain.ProblemID) (domain.Problem, error) {
	args := m.Called(ctx, operator, problemID)
	return args.Get(0).(domain.Problem), args.Error(1)
}

func (m *WorkbookMock) AddProblem(ctx context.Context, operator domain.Student, param domain.ProblemAddParameter) (domain.Added, domain.ProblemID, error) {
	args := m.Called(ctx, operator, param)
	return args.Get(0).(domain.Added), args.Get(1).(domain.ProblemID), args.Error(two)
}

func (m *WorkbookMock) UpdateProblem(ctx context.Context, operator domain.Student, id domain.ProblemSelectParameter2, param domain.ProblemUpdateParameter) (domain.Added, domain.Updated, error) {
	args := m.Called(ctx, operator, id, param)
	return args.Get(0).(domain.Added), args.Get(0).(domain.Updated), args.Error(two)
}

func (m *WorkbookMock) RemoveProblem(ctx context.Context, operator domain.Student, id domain.ProblemSelectParameter2) error {
	args := m.Called(ctx, operator, id)
	return args.Error(0)
}

func (m *WorkbookMock) UpdateWorkbook(ctx context.Context, operator domain.Student, version int, parameter domain.WorkbookUpdateParameter) error {
	args := m.Called(ctx, operator, version, parameter)
	return args.Error(0)
}

func (m *WorkbookMock) RemoveWorkbook(ctx context.Context, operator domain.Student, version int) error {
	args := m.Called(ctx, operator, version)
	return args.Error(0)
}
