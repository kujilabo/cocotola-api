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
	pluginD "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	pluginDM "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain/mock"
	pluginSM "github.com/kujilabo/cocotola-api/pkg_plugin/common/service/mock"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/service"
)

var anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })

func englishWordProblemProcessor_Init(t *testing.T) (
	synthesizer *pluginSM.Synthesizer,
	translationClient *pluginSM.TranslationClient,
	tatoebaClient *pluginSM.TatoebaClient,
	operator *appDM.StudentModel,
	workbookModel *appDM.WorkbookModel,
	rf *appSM.RepositoryFactory,
	problemRepo *appSM.ProblemRepository,
	englishWordProblemProcessor service.EnglishPhraseProblemProcessor) {

	synthesizer = new(pluginSM.Synthesizer)
	translationClient = new(pluginSM.TranslationClient)
	tatoebaClient = new(pluginSM.TatoebaClient)
	operator = new(appDM.StudentModel)
	problemRepo = new(appSM.ProblemRepository)
	rf = new(appSM.RepositoryFactory)
	rf.On("NewProblemRepository", anythingOfContext, domain.EnglishWordProblemType).Return(problemRepo, nil)
	workbookModel = new(appDM.WorkbookModel)
	englishWordProblemProcessor = service.NewEnglishWordProblemProcessor(synthesizer, translationClient, tatoebaClient, nil)
	return
}

func Test_englishWordProblemProcessor_AddProblem_singleProblem_audioDisabled(t *testing.T) {
	ctx := context.Background()
	_, _, _, operator, workbookModel, rf, problemRepo, processor := englishWordProblemProcessor_Init(t)

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
	// then
	require.Equal(t, 1, len(problemIDs))
	require.Equal(t, 100, int(problemIDs[0]))
	paramCheck := mock.MatchedBy(func(p appS.ProblemAddParameter) bool {
		require.Equal(t, 1, int(p.GetWorkbookID()))
		require.Equal(t, 2, p.GetNumber())
		require.Equal(t, "ペン", p.GetProperties()["translated"])
		require.Equal(t, "pen", p.GetProperties()["text"])
		require.Equal(t, "ja", p.GetProperties()["lang"])
		require.Equal(t, "0", p.GetProperties()["audioId"])
		require.Equal(t, "6", p.GetProperties()["pos"])
		require.Len(t, p.GetProperties(), 5)
		return true
	})
	problemRepo.AssertCalled(t, "AddProblem", anythingOfContext, operator, paramCheck)
	problemRepo.AssertNumberOfCalls(t, "AddProblem", 1)
}

func testNewTranslation(pos pluginD.WordPos, translated string) *pluginDM.Translation {
	translation := new(pluginDM.Translation)
	translation.On("GetPos").Return(pos)
	translation.On("GetTranslated").Return(translated)
	return translation
}

func Test_englishWordProblemProcessor_AddProblem_multipleProblem_audioDisabled(t *testing.T) {
	ctx := context.Background()
	_, translationClient, _, operator, workbookModel, rf, problemRepo, processor := englishWordProblemProcessor_Init(t)

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
		"pos":        "99",
		"text":       "book",
		"translated": "",
		"lang":       "ja",
	})
	// - problemRepo
	problemRepo.On("AddProblem", anythingOfContext, operator, mock.Anything).Return(appD.ProblemID(100), nil)
	// - translationClient
	t1 := testNewTranslation(pluginD.PosNoun, "本")
	t2 := testNewTranslation(pluginD.PosVerb, "予約する")
	translationClient.On("DictionaryLookup", anythingOfContext, appD.Lang2EN, appD.Lang2JA, "book").Return([]pluginD.Translation{t1, t2}, nil)
	// when
	problemIDs, err := processor.AddProblem(ctx, rf, operator, workbookModel, param)
	require.NoError(t, err)
	// then
	require.Equal(t, 2, len(problemIDs))
	require.Equal(t, 100, int(problemIDs[0]))
	// paramCheck := mock.MatchedBy(func(p appS.ProblemAddParameter) bool {
	// 	fmt.Println(p)
	// 	require.Equal(t, 1, int(p.GetWorkbookID()))
	// 	require.Equal(t, 2, p.GetNumber())
	// 	require.Equal(t, "本", p.GetProperties()["translated"])
	// 	require.Equal(t, "book", p.GetProperties()["text"])
	// 	require.Equal(t, "ja", p.GetProperties()["lang"])
	// 	require.Equal(t, "0", p.GetProperties()["audioId"])
	// 	require.Equal(t, "6", p.GetProperties()["pos"])
	// 	require.Len(t, p.GetProperties(), 5)
	// 	return true
	// })
	// paramCheck2 := mock.MatchedBy(func(p appS.ProblemAddParameter) bool {
	// 	fmt.Println(p)
	// 	require.Equal(t, 1, int(p.GetWorkbookID()))
	// 	require.Equal(t, 2, p.GetNumber())
	// 	require.Equal(t, "予約する", p.GetProperties()["translated"])
	// 	require.Equal(t, "book", p.GetProperties()["text"])
	// 	require.Equal(t, "ja", p.GetProperties()["lang"])
	// 	require.Equal(t, "0", p.GetProperties()["audioId"])
	// 	require.Equal(t, "6", p.GetProperties()["pos"])
	// 	require.Len(t, p.GetProperties(), 5)
	// 	return true
	// })
	// problemRepo.AssertCalled(t, "AddProblem", anythingOfContext, operator, paramCheck)
	// problemRepo.AssertCalled(t, "AddProblem", anythingOfContext, operator, paramCheck2)
	problemRepo.AssertNumberOfCalls(t, "AddProblem", 2)
}
