package domain_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type OrganizationRepositoryMock struct {
	mock.Mock
}

func (m *OrganizationRepositoryMock) GetOrganization(ctx context.Context, operator domain.AppUser) (domain.Organization, error) {
	args := m.Called(ctx, operator)
	return args.Get(0).(domain.Organization), args.Error(1)
}

func (m *OrganizationRepositoryMock) FindOrganizationByName(ctx context.Context, operator domain.SystemAdmin, name string) (domain.Organization, error) {
	args := m.Called(ctx, operator, name)
	return args.Get(0).(domain.Organization), args.Error(1)
}

func (m *OrganizationRepositoryMock) FindOrganizationByID(ctx context.Context, operator domain.SystemAdmin, id domain.OrganizationID) (domain.Organization, error) {
	args := m.Called(ctx, operator, id)
	return args.Get(0).(domain.Organization), args.Error(1)
}

func (m *OrganizationRepositoryMock) AddOrganization(ctx context.Context, operator domain.SystemAdmin, param domain.OrganizationAddParameter) (domain.OrganizationID, error) {
	args := m.Called(ctx, operator, param)
	return args.Get(0).(domain.OrganizationID), args.Error(1)
}
