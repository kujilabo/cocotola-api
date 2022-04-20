package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
	mocks "github.com/kujilabo/cocotola-api/pkg_app/service/mock"
	user_mock "github.com/kujilabo/cocotola-api/pkg_user/domain/mock"
	userSM "github.com/kujilabo/cocotola-api/pkg_user/service/mock"
)

const (
	problemType1 = "PROBLEM_TYPE_1"
	problemType2 = "PROBLEM_TYPE_2"
)

func student_Init(t *testing.T, ctx context.Context) (
	spaceRepo *userSM.SpaceRepositoryMock,
	userRf *userSM.RepositoryFactoryMock,
	workbookRepo *mocks.WorkbookRepository,
	userQuotaRepo *mocks.UserQuotaRepository,
	rf *mocks.RepositoryFactory,
	problemQuotaProcessor *mocks.ProblemQuotaProcessor,
	pf *mocks.ProcessorFactory) {

	workbookRepo = new(mocks.WorkbookRepository)
	userQuotaRepo = new(mocks.UserQuotaRepository)
	rf = new(mocks.RepositoryFactory)
	rf.On("NewWorkbookRepository", ctx).Return(workbookRepo, nil)
	rf.On("NewUserQuotaRepository", ctx).Return(userQuotaRepo, nil)

	problemQuotaProcessor = new(mocks.ProblemQuotaProcessor)
	pf = new(mocks.ProcessorFactory)
	pf.On("NewProblemQuotaProcessor", problemType1).Return(problemQuotaProcessor, nil)
	pf.On("NewProblemQuotaProcessor", problemType2).Return(problemQuotaProcessor, nil)

	userRf = new(userSM.RepositoryFactoryMock)
	spaceRepo = new(userSM.SpaceRepositoryMock)
	userRf.On("NewSpaceRepository").Return(spaceRepo)

	// return spaceRepo, userRf, workbookRepo, userQuotaRepo, rf, problemQuotaProcessor, pf
	return
}

func Test_student_GetDefaultSpace(t *testing.T) {
	ctx := context.Background()
	spaceRepo, userRf, _, _, _, _, _ := student_Init(t, ctx)

	userRf.On("NewSpaceRepository").Return(spaceRepo)
	expected := new(user_mock.SpaceModel)
	spaceRepo.On("FindDefaultSpace", ctx, mock.Anything).Return(expected, nil)
	studentModel, err := domain.NewStudentModel(nil)
	require.NoError(t, err)
	student, err := service.NewStudent(nil, nil, userRf, studentModel)
	require.NoError(t, err)
	// given
	expected.On("GetKey").Return("KEY")
	// when
	actual, err := student.GetDefaultSpace(ctx)
	require.NoError(t, err)
	// then
	require.Equal(t, "KEY", actual.GetKey())
	spaceRepo.AssertCalled(t, "FindDefaultSpace", ctx, mock.Anything)
	spaceRepo.AssertNumberOfCalls(t, "FindDefaultSpace", 1)
}

func Test_student_GetPersonalSpace(t *testing.T) {
	ctx := context.Background()
	spaceRepo, userRf, _, _, _, _, _ := student_Init(t, ctx)

	expected := new(user_mock.SpaceModel)
	spaceRepo.On("FindPersonalSpace", ctx, mock.Anything).Return(expected, nil)
	studentModel, err := domain.NewStudentModel(nil)
	require.NoError(t, err)
	student, err := service.NewStudent(nil, nil, userRf, studentModel)
	require.NoError(t, err)
	// given
	expected.On("GetKey").Return("KEY")
	// when
	actual, err := student.GetPersonalSpace(ctx)
	require.NoError(t, err)
	// then
	require.Equal(t, "KEY", actual.GetKey())
	spaceRepo.AssertCalled(t, "FindPersonalSpace", ctx, mock.Anything)
	spaceRepo.AssertNumberOfCalls(t, "FindPersonalSpace", 1)
}

func Test_student_FindWorkbooksFromPersonalSpace(t *testing.T) {
	ctx := context.Background()
	spaceRepo, userRf, workbookRepo, _, rf, _, _ := student_Init(t, ctx)

	space := new(user_mock.SpaceModel)
	space.On("GetID").Return(uint(100))
	spaceRepo.On("FindPersonalSpace", ctx, mock.Anything).Return(space, nil)

	studentModel, err := domain.NewStudentModel(nil)
	require.NoError(t, err)
	student, err := service.NewStudent(nil, rf, userRf, studentModel)
	require.NoError(t, err)
	// given
	expected, err := service.NewWorkbookSearchResult(123, nil)
	require.NoError(t, err)
	workbookRepo.On("FindPersonalWorkbooks", ctx, mock.Anything, mock.Anything).Return(expected, nil)
	// when
	condition, err := service.NewWorkbookSearchCondition(1, 100, nil)
	require.NoError(t, err)
	actual, err := student.FindWorkbooksFromPersonalSpace(ctx, condition)
	require.NoError(t, err)
	// then
	require.Equal(t, 123, actual.GetTotalCount())
	spaceRepo.AssertCalled(t, "FindPersonalSpace", ctx, mock.Anything)
	spaceRepo.AssertNumberOfCalls(t, "FindPersonalSpace", 1)
}

func Test_student_FindWorkbookByID(t *testing.T) {
	ctx := context.Background()
	_, userRf, workbookRepo, _, rf, _, _ := student_Init(t, ctx)

	expected := new(mocks.Workbook)
	workbookRepo.On("FindWorkbookByID", ctx, mock.Anything, mock.Anything).Return(expected, nil)

	studentModel, err := domain.NewStudentModel(nil)
	require.NoError(t, err)
	student, err := service.NewStudent(nil, rf, userRf, studentModel)
	require.NoError(t, err)
	// given
	expected.On("GetID").Return(uint(123))
	// when
	actual, err := student.FindWorkbookByID(ctx, domain.WorkbookID(100))
	require.NoError(t, err)
	// then
	require.Equal(t, uint(123), actual.GetID())
}

func Test_student_CheckQuota(t *testing.T) {
	ctx := context.Background()

	type args struct {
		problemType string
		name        service.QuotaName
	}
	tests := []struct {
		name              string
		isExceeded        bool
		problemTypeSuffix string
		quotaUnit         service.QuotaUnit
		quotaLimit        int
		args              args
		err               error
	}{
		{
			name:              "QuotaNameSize,isNotExceeded",
			isExceeded:        false,
			problemTypeSuffix: "_size",
			quotaUnit:         service.QuotaUnitPersitance,
			quotaLimit:        234,
			args: args{
				problemType: problemType1,
				name:        service.QuotaNameSize,
			},
			err: nil,
		},
		{
			name:              "QuotaNameSize,isExceeded",
			isExceeded:        true,
			problemTypeSuffix: "_size",
			quotaUnit:         service.QuotaUnitPersitance,
			quotaLimit:        234,
			args: args{
				problemType: problemType2,
				name:        service.QuotaNameSize,
			},
			err: service.ErrQuotaExceeded,
		},
		{
			name:              "QuotaNameUpdate,isNotExceeded",
			isExceeded:        false,
			problemTypeSuffix: "_update",
			quotaUnit:         service.QuotaUnitDay,
			quotaLimit:        345,
			args: args{
				problemType: problemType1,
				name:        service.QuotaNameUpdate,
			},
			err: nil,
		},
		{
			name:              "QuotaNameUpdate,isExceeded",
			isExceeded:        true,
			problemTypeSuffix: "_update",
			quotaUnit:         service.QuotaUnitDay,
			quotaLimit:        345,
			args: args{
				problemType: problemType2,
				name:        service.QuotaNameUpdate,
			},
			err: service.ErrQuotaExceeded,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, userRf, _, userQuotaRepo, rf, problemQuotaProcessor, pf := student_Init(t, ctx)
			userQuotaRepo.On("IsExceeded", mock.Anything, mock.Anything, tt.args.problemType+tt.problemTypeSuffix, tt.quotaUnit, tt.quotaLimit).Return(tt.isExceeded, nil)
			problemQuotaProcessor.On("GetUnitForSizeQuota").Return(service.QuotaUnitPersitance)
			problemQuotaProcessor.On("GetLimitForSizeQuota").Return(tt.quotaLimit)
			problemQuotaProcessor.On("GetUnitForUpdateQuota").Return(service.QuotaUnitDay)
			problemQuotaProcessor.On("GetLimitForUpdateQuota").Return(tt.quotaLimit)

			s, err := service.NewStudent(pf, rf, userRf, nil)
			require.NoError(t, err)
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
