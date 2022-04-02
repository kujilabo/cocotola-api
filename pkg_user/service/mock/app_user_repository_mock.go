package service_mock

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/service"
	"github.com/stretchr/testify/mock"
)

type AppUserRepositoryMock struct {
	mock.Mock
}

func (m *AppUserRepositoryMock) FindSystemOwnerByOrganizationID(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (service.SystemOwner, error) {
	args := m.Called(ctx, operator, organizationID)
	return args.Get(0).(service.SystemOwner), args.Error(1)
}
func (m *AppUserRepositoryMock) FindSystemOwnerByOrganizationName(ctx context.Context, operator domain.SystemAdminModel, organizationName string) (service.SystemOwner, error) {
	args := m.Called(ctx, operator, organizationName)
	return args.Get(0).(service.SystemOwner), args.Error(1)
}
func (m *AppUserRepositoryMock) FindAppUserByID(ctx context.Context, operator domain.AppUserModel, id domain.AppUserID) (service.AppUser, error) {
	args := m.Called(ctx, operator, id)
	return args.Get(0).(service.AppUser), args.Error(1)
}
func (m *AppUserRepositoryMock) FindAppUserByLoginID(ctx context.Context, operator domain.AppUserModel, loginID string) (service.AppUser, error) {
	args := m.Called(ctx, operator, loginID)
	return args.Get(0).(service.AppUser), args.Error(1)
}
func (m *AppUserRepositoryMock) FindOwnerByLoginID(ctx context.Context, operator domain.SystemOwnerModel, loginID string) (service.Owner, error) {
	args := m.Called(ctx, operator, loginID)
	return args.Get(0).(service.Owner), args.Error(1)
}
func (m *AppUserRepositoryMock) AddAppUser(ctx context.Context, operator domain.OwnerModel, param service.AppUserAddParameter) (domain.AppUserID, error) {
	args := m.Called(ctx, operator, param)
	return args.Get(0).(domain.AppUserID), args.Error(1)
}
func (m *AppUserRepositoryMock) AddSystemOwner(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (domain.AppUserID, error) {
	args := m.Called(ctx, operator, organizationID)
	return args.Get(0).(domain.AppUserID), args.Error(1)
}
func (m *AppUserRepositoryMock) AddFirstOwner(ctx context.Context, operator domain.SystemOwnerModel, param service.FirstOwnerAddParameter) (domain.AppUserID, error) {
	args := m.Called(ctx, operator, param)
	return args.Get(0).(domain.AppUserID), args.Error(1)
}
