package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	domain_mock "github.com/kujilabo/cocotola-api/pkg_app/domain/mock"
	user_mock "github.com/kujilabo/cocotola-api/pkg_user/domain/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	problemType1 = "PROBLEM_TYPE_1"
	problemType2 = "PROBLEM_TYPE_2"
)

func student_Init(t *testing.T, ctx context.Context) (
	spaceRepo *user_mock.SpaceRepositoryMock,
	userRf *user_mock.RepositoryFactoryMock,
	workbookRepo *domain_mock.WorkbookRepositoryMock,
	userQuotaRepo *domain_mock.UserQuotaRepositoryMock,
	rf *domain_mock.RepositoryFactoryMock,
	problemQuotaProcessor *domain_mock.ProblemQuotaProcessorMock,
	pf *domain_mock.ProcessorFactoryMock) {

	workbookRepo = new(domain_mock.WorkbookRepositoryMock)
	userQuotaRepo = new(domain_mock.UserQuotaRepositoryMock)
	rf = new(domain_mock.RepositoryFactoryMock)
	rf.On("NewWorkbookRepository", ctx).Return(workbookRepo, nil)
	rf.On("NewUserQuotaRepository", ctx).Return(userQuotaRepo, nil)

	problemQuotaProcessor = new(domain_mock.ProblemQuotaProcessorMock)
	pf = new(domain_mock.ProcessorFactoryMock)
	pf.On("NewProblemQuotaProcessor", problemType1).Return(problemQuotaProcessor, nil)
	pf.On("NewProblemQuotaProcessor", problemType2).Return(problemQuotaProcessor, nil)

	userRf = new(user_mock.RepositoryFactoryMock)
	spaceRepo = new(user_mock.SpaceRepositoryMock)
	userRf.On("NewSpaceRepository").Return(spaceRepo)

	// return spaceRepo, userRf, workbookRepo, userQuotaRepo, rf, problemQuotaProcessor, pf
	return
}

func Test_student_GetDefaultSpace(t *testing.T) {
	ctx := context.Background()
	spaceRepo, userRf, _, _, _, _, _ := student_Init(t, ctx)

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
	spaceRepo, userRf, _, _, _, _, _ := student_Init(t, ctx)

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
	spaceRepo, userRf, workbookRepo, _, rf, _, _ := student_Init(t, ctx)

	space := new(user_mock.SpaceMock)
	space.On("GetID").Return(uint(100))
	spaceRepo.On("FindPersonalSpace", ctx, mock.Anything).Return(space, nil)

	s, err := domain.NewStudent(nil, rf, userRf, nil)
	assert.NoError(t, err)
	// given
	expected, err := domain.NewWorkbookSearchResult(123, nil)
	assert.NoError(t, err)
	workbookRepo.On("FindPersonalWorkbooks", ctx, mock.Anything, mock.Anything).Return(expected, nil)
	// when
	condition, err := domain.NewWorkbookSearchCondition(1, 100, nil)
	assert.NoError(t, err)
	actual, err := s.FindWorkbooksFromPersonalSpace(ctx, condition)
	assert.NoError(t, err)
	// then
	assert.Equal(t, 123, actual.GetTotalCount())
	spaceRepo.AssertCalled(t, "FindPersonalSpace", ctx, mock.Anything)
	spaceRepo.AssertNumberOfCalls(t, "FindPersonalSpace", 1)
}

func Test_student_FindWorkbookByID(t *testing.T) {
	ctx := context.Background()
	_, userRf, workbookRepo, _, rf, _, _ := student_Init(t, ctx)

	expected := new(domain_mock.WorkbookMock)
	workbookRepo.On("FindWorkbookByID", ctx, mock.Anything, mock.Anything).Return(expected, nil)

	s, err := domain.NewStudent(nil, rf, userRf, nil)
	assert.NoError(t, err)
	// given
	expected.On("GetID").Return(uint(123))
	// when
	actual, err := s.FindWorkbookByID(ctx, domain.WorkbookID(100))
	assert.NoError(t, err)
	// then
	assert.Equal(t, uint(123), actual.GetID())
}

func Test_student_CheckQuota(t *testing.T) {
	ctx := context.Background()

	type args struct {
		problemType string
		name        domain.QuotaName
	}
	tests := []struct {
		name              string
		isExceeded        bool
		problemTypeSuffix string
		quotaUnit         domain.QuotaUnit
		quotaLimit        int
		args              args
		err               error
	}{
		{
			name:              "QuotaNameSize,isNotExceeded",
			isExceeded:        false,
			problemTypeSuffix: "_size",
			quotaUnit:         domain.QuotaUnitPersitance,
			quotaLimit:        234,
			args: args{
				problemType: problemType1,
				name:        domain.QuotaNameSize,
			},
			err: nil,
		},
		{
			name:              "QuotaNameSize,isExceeded",
			isExceeded:        true,
			problemTypeSuffix: "_size",
			quotaUnit:         domain.QuotaUnitPersitance,
			quotaLimit:        234,
			args: args{
				problemType: problemType2,
				name:        domain.QuotaNameSize,
			},
			err: domain.ErrQuotaExceeded,
		},
		{
			name:              "QuotaNameUpdate,isNotExceeded",
			isExceeded:        false,
			problemTypeSuffix: "_update",
			quotaUnit:         domain.QuotaUnitDay,
			quotaLimit:        345,
			args: args{
				problemType: problemType1,
				name:        domain.QuotaNameUpdate,
			},
			err: nil,
		},
		{
			name:              "QuotaNameUpdate,isExceeded",
			isExceeded:        true,
			problemTypeSuffix: "_update",
			quotaUnit:         domain.QuotaUnitDay,
			quotaLimit:        345,
			args: args{
				problemType: problemType2,
				name:        domain.QuotaNameUpdate,
			},
			err: domain.ErrQuotaExceeded,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, userRf, _, userQuotaRepo, rf, problemQuotaProcessor, pf := student_Init(t, ctx)
			userQuotaRepo.On("IsExceeded", mock.Anything, mock.Anything, tt.args.problemType+tt.problemTypeSuffix, tt.quotaUnit, tt.quotaLimit).Return(tt.isExceeded, nil)
			problemQuotaProcessor.On("GetUnitForSizeQuota").Return(domain.QuotaUnitPersitance)
			problemQuotaProcessor.On("GetLimitForSizeQuota").Return(tt.quotaLimit)
			problemQuotaProcessor.On("GetUnitForUpdateQuota").Return(domain.QuotaUnitDay)
			problemQuotaProcessor.On("GetLimitForUpdateQuota").Return(tt.quotaLimit)

			s, err := domain.NewStudent(pf, rf, userRf, nil)
			assert.NoError(t, err)
			err = s.CheckQuota(ctx, tt.args.problemType, tt.args.name)
			if err == nil && tt.err != nil {
				t.Errorf("student.CheckQuota() error = %v, err %v", err, tt.err)
			} else if err != nil && tt.err == nil {
				t.Errorf("student.CheckQuota() error = %v, err %v", err, tt.err)
			} else if err != nil && tt.err != nil && !errors.Is(err, tt.err) {
				t.Errorf("student.CheckQuota() error = %v, err %v", err, tt.err)
			}
			userQuotaRepo.AssertCalled(t, "IsExceeded", mock.Anything, mock.Anything, tt.args.problemType+tt.problemTypeSuffix, tt.quotaUnit, tt.quotaLimit)
			userQuotaRepo.AssertNumberOfCalls(t, "IsExceeded", 1)
		})
	}
}
