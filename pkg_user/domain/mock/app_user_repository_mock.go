package domain_mock

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/stretchr/testify/mock"
)

type AppUserRepositoryMock struct {
	mock.Mock
}

func (m *AppUserRepositoryMock) FindSystemOwnerByOrganizationID(ctx context.Context, operator domain.SystemAdmin, organizationID domain.OrganizationID) (domain.SystemOwner, error) {
	args := m.Called(ctx, operator, organizationID)
	return args.Get(0).(domain.SystemOwner), args.Error(1)
}
func (m *AppUserRepositoryMock) FindSystemOwnerByOrganizationName(ctx context.Context, operator domain.SystemAdmin, organizationName string) (domain.SystemOwner, error) {
	args := m.Called(ctx, operator, organizationName)
	return args.Get(0).(domain.SystemOwner), args.Error(1)
}
func (m *AppUserRepositoryMock) FindAppUserByID(ctx context.Context, operator domain.AppUser, id domain.AppUserID) (domain.AppUser, error) {
	args := m.Called(ctx, operator, id)
	return args.Get(0).(domain.AppUser), args.Error(1)
}
func (m *AppUserRepositoryMock) FindAppUserByLoginID(ctx context.Context, operator domain.AppUser, loginID string) (domain.AppUser, error) {
	args := m.Called(ctx, operator, loginID)
	return args.Get(0).(domain.AppUser), args.Error(1)
}
func (m *AppUserRepositoryMock) FindOwnerByLoginID(ctx context.Context, operator domain.SystemOwner, loginID string) (domain.Owner, error) {
	args := m.Called(ctx, operator, loginID)
	return args.Get(0).(domain.Owner), args.Error(1)
}
func (m *AppUserRepositoryMock) AddAppUser(ctx context.Context, operator domain.Owner, param domain.AppUserAddParameter) (domain.AppUserID, error) {
	args := m.Called(ctx, operator, param)
	return args.Get(0).(domain.AppUserID), args.Error(1)
}
func (m *AppUserRepositoryMock) AddSystemOwner(ctx context.Context, operator domain.SystemAdmin, organizationID domain.OrganizationID) (domain.AppUserID, error) {
	args := m.Called(ctx, operator, organizationID)
	return args.Get(0).(domain.AppUserID), args.Error(1)
}
func (m *AppUserRepositoryMock) AddFirstOwner(ctx context.Context, operator domain.SystemOwner, param domain.FirstOwnerAddParameter) (domain.AppUserID, error) {
	args := m.Called(ctx, operator, param)
	return args.Get(0).(domain.AppUserID), args.Error(1)
}
