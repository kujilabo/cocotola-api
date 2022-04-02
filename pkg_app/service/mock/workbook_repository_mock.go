package service_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type WorkbookRepositoryMock struct {
	mock.Mock
}

func (m *WorkbookRepositoryMock) FindPersonalWorkbooks(ctx context.Context, operator domain.StudentModel, param service.WorkbookSearchCondition) (service.WorkbookSearchResult, error) {
	args := m.Called(ctx, operator, param)
	return args.Get(0).(service.WorkbookSearchResult), args.Error(1)
}
func (m *WorkbookRepositoryMock) FindWorkbookByID(ctx context.Context, operator domain.StudentModel, id domain.WorkbookID) (service.Workbook, error) {
	args := m.Called(ctx, operator, id)
	return args.Get(0).(service.Workbook), args.Error(1)
}
func (m *WorkbookRepositoryMock) FindWorkbookByName(ctx context.Context, operator user.AppUserModel, spaceID user.SpaceID, name string) (service.Workbook, error) {
	args := m.Called(ctx, operator, spaceID, name)
	return args.Get(0).(service.Workbook), args.Error(1)
}
func (m *WorkbookRepositoryMock) AddWorkbook(ctx context.Context, operator user.AppUserModel, spaceID user.SpaceID, param service.WorkbookAddParameter) (domain.WorkbookID, error) {
	args := m.Called(ctx, operator, spaceID, param)
	return args.Get(0).(domain.WorkbookID), args.Error(1)
}
func (m *WorkbookRepositoryMock) UpdateWorkbook(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, version int, param service.WorkbookUpdateParameter) error {
	args := m.Called(ctx, operator, workbookID, version, param)
	return args.Error(0)
}
func (m *WorkbookRepositoryMock) RemoveWorkbook(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, version int) error {
	args := m.Called(ctx, operator, workbookID, version)
	return args.Error(0)
}
