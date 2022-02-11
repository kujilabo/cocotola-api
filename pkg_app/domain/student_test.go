package domain_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	domain_mock "github.com/kujilabo/cocotola-api/pkg_app/domain/mock"
	user_mock "github.com/kujilabo/cocotola-api/pkg_user/domain/mock"
)

func student_Init(t *testing.T, ctx context.Context) (*user_mock.SpaceRepositoryMock, *user_mock.RepositoryFactoryMock, *domain_mock.WorkbookRepositoryMock, *domain_mock.RepositoryFactoryMock) {

	workbookRepo := new(domain_mock.WorkbookRepositoryMock)
	rf := new(domain_mock.RepositoryFactoryMock)
	rf.On("NewWorkbookRepository", ctx).Return(workbookRepo, nil)
	userRf := new(user_mock.RepositoryFactoryMock)
	spaceRepo := new(user_mock.SpaceRepositoryMock)
	userRf.On("NewSpaceRepository").Return(spaceRepo)

	return spaceRepo, userRf, workbookRepo, rf
}

func Test_student_GetDefaultSpace(t *testing.T) {
	ctx := context.Background()
	userRf := new(user_mock.RepositoryFactoryMock)
	spaceRepo := new(user_mock.SpaceRepositoryMock)
	userRf.On("NewSpaceRepository").Return(spaceRepo)
	expected := new(user_mock.SpaceMock)
	spaceRepo.On("FindDefaultSpace", ctx, mock.Anything).Return(expected, nil)
	s, err := domain.NewStudent(nil, nil, userRf, nil)
	assert.NoError(t, err)
	// given
	expected.On("GetKey").Return("KEY")
	// when
	actual, err := s.GetDefaultSpace(ctx)
	assert.NoError(t, err)
	// then
	assert.Equal(t, "KEY", actual.GetKey())
	spaceRepo.AssertCalled(t, "FindDefaultSpace", ctx, mock.Anything)
	spaceRepo.AssertNumberOfCalls(t, "FindDefaultSpace", 1)
}

func Test_student_GetPersonalSpace(t *testing.T) {
	ctx := context.Background()
	userRf := new(user_mock.RepositoryFactoryMock)
	spaceRepo := new(user_mock.SpaceRepositoryMock)
	userRf.On("NewSpaceRepository").Return(spaceRepo)
	expected := new(user_mock.SpaceMock)
	spaceRepo.On("FindPersonalSpace", ctx, mock.Anything).Return(expected, nil)
	s, err := domain.NewStudent(nil, nil, userRf, nil)
	assert.NoError(t, err)
	// given
	expected.On("GetKey").Return("KEY")
	// when
	actual, err := s.GetPersonalSpace(ctx)
	assert.NoError(t, err)
	// then
	assert.Equal(t, "KEY", actual.GetKey())
	spaceRepo.AssertCalled(t, "FindPersonalSpace", ctx, mock.Anything)
	spaceRepo.AssertNumberOfCalls(t, "FindPersonalSpace", 1)
}

func Test_student_FindWorkbooksFromPersonalSpace(t *testing.T) {
	ctx := context.Background()
	expected := domain.WorkbookSearchResult{}
	workbookRepo := new(domain_mock.WorkbookRepositoryMock)
	workbookRepo.On("FindPersonalWorkbooks", ctx, mock.Anything, mock.Anything).Return(&expected, nil)
	rf := new(domain_mock.RepositoryFactoryMock)
	rf.On("NewWorkbookRepository", ctx).Return(workbookRepo, nil)
	userRf := new(user_mock.RepositoryFactoryMock)
	spaceRepo := new(user_mock.SpaceRepositoryMock)
	userRf.On("NewSpaceRepository").Return(spaceRepo)
	space := new(user_mock.SpaceMock)
	spaceRepo.On("FindPersonalSpace", ctx, mock.Anything).Return(space, nil)

	s, err := domain.NewStudent(nil, rf, userRf, nil)
	assert.NoError(t, err)
	// given
	expected.TotalCount = 100
	// when
	condition, err := domain.NewWorkbookSearchCondition(1, 100, nil)
	assert.NoError(t, err)
	actual, err := s.FindWorkbooksFromPersonalSpace(ctx, condition)
	assert.NoError(t, err)
	// then
	assert.Equal(t, expected.TotalCount, actual.TotalCount)
	spaceRepo.AssertCalled(t, "FindPersonalSpace", ctx, mock.Anything)
	spaceRepo.AssertNumberOfCalls(t, "FindPersonalSpace", 1)
}

func Test_student_FindWorkbookByID(t *testing.T) {
	ctx := context.Background()
	_, userRf, workbookRepo, rf := student_Init(t, ctx)

	expected := new(domain_mock.WorkbookMock)
	workbookRepo.On("FindWorkbookByID", ctx, mock.Anything, mock.Anything).Return(expected, nil)

	s, err := domain.NewStudent(nil, rf, userRf, nil)
	assert.NoError(t, err)
	// given
	expected.On("GetID").Return(uint(100))
	// when
	actual, err := s.FindWorkbookByID(ctx, domain.WorkbookID(100))
	assert.NoError(t, err)
	// then
	assert.Equal(t, uint(100), actual.GetID())
}
