package domain_mock

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/stretchr/testify/mock"
)

type WorkbookRepositoryMock struct {
	mock.Mock
}

func (m *WorkbookRepositoryMock) FindPersonalWorkbooks(ctx context.Context, operator domain.Student, param domain.WorkbookSearchCondition) (domain.WorkbookSearchResult, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.WorkbookSearchResult), args.Error(1)
}
func (m *WorkbookRepositoryMock) FindWorkbookByID(ctx context.Context, operator domain.Student, id domain.WorkbookID) (domain.Workbook, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.Workbook), args.Error(1)
}
func (m *WorkbookRepositoryMock) FindWorkbookByName(ctx context.Context, operator domain.Student, spaceID user.SpaceID, name string) (domain.Workbook, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.Workbook), args.Error(1)
}
func (m *WorkbookRepositoryMock) AddWorkbook(ctx context.Context, operator domain.Student, spaceID user.SpaceID, param domain.WorkbookAddParameter) (domain.WorkbookID, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.WorkbookID), args.Error(1)
}
func (m *WorkbookRepositoryMock) UpdateWorkbook(ctx context.Context, operator domain.Student, workbookID domain.WorkbookID, version int, param domain.WorkbookUpdateParameter) error {
	args := m.Called(ctx)
	return args.Error(0)
}
func (m *WorkbookRepositoryMock) RemoveWorkbook(ctx context.Context, operator domain.Student, workbookID domain.WorkbookID, version int) error {
	args := m.Called(ctx)
	return args.Error(0)
}
