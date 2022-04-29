package service_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola-api/src/user/service"
)

type OrganizationRepositoryMock struct {
	mock.Mock
}

func (m *OrganizationRepositoryMock) GetOrganization(ctx context.Context, operator domain.AppUserModel) (service.Organization, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(service.Organization), args.Error(1)
}

func (m *OrganizationRepositoryMock) FindOrganizationByName(ctx context.Context, operator domain.SystemAdminModel, name string) (service.Organization, error) {
	args := m.Called(ctx, operator, name)
	return args.Get(0).(service.Organization), args.Error(1)
}

func (m *OrganizationRepositoryMock) FindOrganizationByID(ctx context.Context, operator domain.SystemAdminModel, id domain.OrganizationID) (service.Organization, error) {
	args := m.Called(ctx, operator, id)
	return args.Get(0).(service.Organization), args.Error(1)
}

func (m *OrganizationRepositoryMock) AddOrganization(ctx context.Context, operator domain.SystemAdminModel, param service.OrganizationAddParameter) (domain.OrganizationID, error) {
	args := m.Called(ctx, operator, param)
	return args.Get(0).(domain.OrganizationID), args.Error(1)
}
