package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
	appDM "github.com/kujilabo/cocotola-api/pkg_app/domain/mock"
	appS "github.com/kujilabo/cocotola-api/pkg_app/service"
	appSM "github.com/kujilabo/cocotola-api/pkg_app/service/mock"
	pluginSM "github.com/kujilabo/cocotola-api/pkg_plugin/common/service/mock"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/service"
)

var anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })

func Test_englishWordProblemProcessor_AddProblem_singleProblem_audioDisabled(t *testing.T) {
	ctx := context.Background()

	operator := new(appDM.StudentModel)
	synthesizer := new(pluginSM.Synthesizer)
	translationClient := new(pluginSM.TranslationClient)
	tatoebaClient := new(pluginSM.TatoebaClient)
	problemRepo := new(appSM.ProblemRepository)
	rf := new(appSM.RepositoryFactory)
	rf.On("NewProblemRepository", anythingOfContext, domain.EnglishWordProblemType).Return(problemRepo, nil)
	workbookModel := new(appDM.WorkbookModel)
	processor := service.NewEnglishWordProblemProcessor(synthesizer, translationClient, tatoebaClient, nil)

	// given
	// - workbook
	workbookModel.On("GetProperties").Return(map[string]string{
		"audioEnabled": "false",
	})
	// - param
	param := new(appSM.ProblemAddParameter)
	param.On("GetWorkbookID").Return(appD.WorkbookID(1))
	param.On("GetNumber").Return(2)
	param.On("GetProperties").Return(map[string]string{
		"pos":        "6",
		"text":       "pen",
		"translated": "ペン",
		"lang":       "ja",
	})
	// - problemRepo
	problemRepo.On("AddProblem", anythingOfContext, operator, mock.Anything).Return(appD.ProblemID(100), nil)
	// when
	problemIDs, err := processor.AddProblem(ctx, rf, operator, workbookModel, param)
	require.NoError(t, err)
	// given
	require.Equal(t, 1, len(problemIDs))
	require.Equal(t, 100, int(problemIDs[0]))
	paramCheck := mock.MatchedBy(func(p appS.ProblemAddParameter) bool {
		require.Equal(t, 1, int(p.GetWorkbookID()))
		require.Equal(t, 2, p.GetNumber())
		return true
	})
	problemRepo.AssertCalled(t, "AddProblem", anythingOfContext, operator, paramCheck)
	problemRepo.AssertNumberOfCalls(t, "AddProblem", 1)
}
