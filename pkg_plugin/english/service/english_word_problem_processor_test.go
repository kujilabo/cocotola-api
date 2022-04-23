package service_test

import (
	"context"
	"testing"

	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
	appDM "github.com/kujilabo/cocotola-api/pkg_app/domain/mock"
	appS "github.com/kujilabo/cocotola-api/pkg_app/service"
	appSM "github.com/kujilabo/cocotola-api/pkg_app/service/mock"
	pluginD "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	pluginDM "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain/mock"
	pluginSM "github.com/kujilabo/cocotola-api/pkg_plugin/common/service/mock"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })

func englishWordProblemProcessor_Init(t *testing.T) (
	synthesizerClient *pluginSM.SynthesizerClient,
	translatorClient *pluginSM.TranslationClient,
	tatoebaClient *pluginSM.TatoebaClient,
	operator *appDM.StudentModel,
	workbookModel *appDM.WorkbookModel,
	rf *appSM.RepositoryFactory,
	problemRepo *appSM.ProblemRepository,
	englishWordProblemProcessor service.EnglishWordProblemProcessor) {

	synthesizerClient = new(pluginSM.SynthesizerClient)
	translatorClient = new(pluginSM.TranslationClient)
	tatoebaClient = new(pluginSM.TatoebaClient)
	operator = new(appDM.StudentModel)
	problemRepo = new(appSM.ProblemRepository)
	rf = new(appSM.RepositoryFactory)
	rf.On("NewProblemRepository", anythingOfContext, domain.EnglishWordProblemType).Return(problemRepo, nil)
	workbookModel = new(appDM.WorkbookModel)
	englishWordProblemProcessor = service.NewEnglishWordProblemProcessor(synthesizerClient, translatorClient, tatoebaClient, nil)
	return
}

func testNewTranslation(pos pluginD.WordPos, translated string) *pluginDM.Translation {
	translation := new(pluginDM.Translation)
	translation.On("GetPos").Return(pos)
	translation.On("GetTranslated").Return(translated)
	return translation
}

func Test_englishWordProblemProcessor_AddProblem_singleProblem_audioDisabled(t *testing.T) {
	ctx := context.Background()
	_, _, _, operator, workbookModel, rf, problemRepo, processor := englishWordProblemProcessor_Init(t)

	// given
	// - workbook
	workbookModel.On("GetProperties").Return(map[string]string{
		"audioEnabled": "false",
	})
	// - problemRepo
	problemRepo.On("AddProblem", anythingOfContext, operator, mock.Anything).Return(appD.ProblemID(100), nil)
	// when
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

func Test_englishWordProblemProcessor_AddProblem_multipleProblem_audioDisabled(t *testing.T) {
	ctx := context.Background()
	_, translatorClient, _, operator, workbookModel, rf, problemRepo, processor := englishWordProblemProcessor_Init(t)

	// given
	// - workbook
	workbookModel.On("GetProperties").Return(map[string]string{
		"audioEnabled": "false",
	})
	// - problemRepo
	problemRepo.On("AddProblem", anythingOfContext, operator, mock.Anything).Return(appD.ProblemID(100), nil)
	// - translatorClient
	translatorClient.On("DictionaryLookup", anythingOfContext, appD.Lang2EN, appD.Lang2JA, "book").Return([]pluginD.Translation{
		testNewTranslation(pluginD.PosNoun, "本"),
		testNewTranslation(pluginD.PosVerb, "予約する"),
	}, nil)
	// when
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
	problemIDs, err := processor.AddProblem(ctx, rf, operator, workbookModel, param)
	require.NoError(t, err)
	// then
	require.Equal(t, 2, len(problemIDs))
	require.Equal(t, 100, int(problemIDs[0]))
	problemRepo.AssertNumberOfCalls(t, "AddProblem", 2)
	{
		param := (problemRepo.Calls[0].Arguments[2]).(appS.ProblemAddParameter)
		assert.Equal(t, 1, int(param.GetWorkbookID()))
		assert.Equal(t, 2, param.GetNumber())
		assert.Equal(t, "本", param.GetProperties()["translated"])
		assert.Equal(t, "book", param.GetProperties()["text"])
		assert.Equal(t, "ja", param.GetProperties()["lang"])
		assert.Equal(t, "0", param.GetProperties()["audioId"])
		assert.Equal(t, "6", param.GetProperties()["pos"])
	}
	{
		param := (problemRepo.Calls[1].Arguments[2]).(appS.ProblemAddParameter)
		assert.Equal(t, 1, int(param.GetWorkbookID()))
		assert.Equal(t, 2, param.GetNumber())
		assert.Equal(t, "予約する", param.GetProperties()["translated"])
		assert.Equal(t, "book", param.GetProperties()["text"])
		assert.Equal(t, "ja", param.GetProperties()["lang"])
		assert.Equal(t, "0", param.GetProperties()["audioId"])
		assert.Equal(t, "9", param.GetProperties()["pos"])
	}
}

func Test_englishWordProblemProcessor_UpdateProblem(t *testing.T) {
	ctx := context.Background()
	_, _, _, operator, workbookModel, rf, problemRepo, processor := englishWordProblemProcessor_Init(t)

	// given
	// - workbook
	workbookModel.On("GetProperties").Return(map[string]string{
		"audioEnabled": "false",
	})
	// - problemRepo
	problemRepo.On("UpdateProblem", anythingOfContext, operator, mock.Anything, mock.Anything).Return(nil)
	// when
	// - param
	paramSelect := new(appSM.ProblemSelectParameter2)
	paramSelect.On("GetProblem")

	param := new(appSM.ProblemUpdateParameter)
	param.On("GetWorkbookID").Return(appD.WorkbookID(1))
	param.On("GetNumber").Return(2)
	param.On("GetProperties").Return(map[string]string{
		"pos":        "6",
		"text":       "pen",
		"translated": "ペン",
		"lang":       "ja",
	})
	added, updated, err := processor.UpdateProblem(ctx, rf, operator, workbookModel, paramSelect, param)
	require.NoError(t, err)
	// then
	require.Equal(t, 1, int(added))
	require.Equal(t, 1, int(updated))
	problemRepo.AssertNumberOfCalls(t, "UpdateProblem", 1)
	{
		param := (problemRepo.Calls[0].Arguments[3]).(appS.ProblemUpdateParameter)
		assert.Equal(t, 2, param.GetNumber())
		assert.Equal(t, "ペン", param.GetProperties()["translated"])
		assert.Equal(t, "pen", param.GetProperties()["text"])
		assert.Equal(t, "0", param.GetProperties()["audioId"])
		assert.Equal(t, "0", param.GetProperties()["sentenceId1"])
		assert.Len(t, param.GetProperties(), 4)
	}
}
