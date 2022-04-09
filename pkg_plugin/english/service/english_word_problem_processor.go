package service

import (
	"context"
	"errors"
	"io"
	"strconv"

	"golang.org/x/xerrors"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	appS "github.com/kujilabo/cocotola-api/pkg_app/service"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	plugin "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	pluginS "github.com/kujilabo/cocotola-api/pkg_plugin/common/service"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
)

var (
	quotaSizeUnit  = appS.QuotaUnitPersitance
	quotaSizeLimit = 5000

	quotaUpdateUnit                            = appS.QuotaUnitDay
	quotaUpdateLimit                           = 100
	EnglishWordProblemUpdatePropertyText       = "text"
	EnglishWordProblemUpdatePropertyTranslated = "translated"
	EnglishWordProblemUpdatePropertyAudioID    = "audioId"
	// EnglishWordProblemUpdatePropertyTatoebaSentenceNumber1 = "tatoebaSentenceNumber1"
	// EnglishWordProblemUpdatePropertyTatoebaSentenceNumber2 = "tatoebaSentenceNumber2"
	EnglishWordProblemUpdatePropertySentenceID1 = "sentenceId1"
)

type englishWordProblemAddParemeter struct {
	Lang       app.Lang2      `validate:"required"`
	Text       string         `validate:"required"`
	Pos        plugin.WordPos `validate:"required"`
	Translated string
}

func toEnglishWordProblemAddParemeter(param appS.ProblemAddParameter) (*englishWordProblemAddParemeter, error) {
	posS := param.GetProperties()["pos"]
	pos, err := strconv.Atoi(posS)
	if err != nil {
		return nil, xerrors.Errorf("failed to cast to int. err: %w", err)
	}

	if _, ok := param.GetProperties()["text"]; !ok {
		return nil, xerrors.Errorf("text is not defined. err: %w", lib.ErrInvalidArgument)
	}

	translated := ""
	if _, ok := param.GetProperties()["translated"]; ok {
		translated = param.GetProperties()["translated"]
	}

	if _, ok := param.GetProperties()["lang"]; !ok {
		return nil, xerrors.Errorf("lang is not defined. err: %w", lib.ErrInvalidArgument)
	}

	lang2, err := app.NewLang2(param.GetProperties()["lang"])
	if err != nil {
		return nil, err
	}

	m := &englishWordProblemAddParemeter{
		Lang:       lang2,
		Text:       param.GetProperties()["text"],
		Pos:        plugin.WordPos(pos),
		Translated: translated,
	}

	return m, lib.Validator.Struct(m)
}

type englishWordProblemUpdateParemeter struct {
	Text       string `validate:"required"`
	Translated string
}

func toEnglishWordProblemUpdateParemeter(param appS.ProblemUpdateParameter) (*englishWordProblemUpdateParemeter, error) {
	if _, ok := param.GetProperties()["text"]; !ok {
		return nil, xerrors.Errorf("text is not defined. err: %w", lib.ErrInvalidArgument)
	}

	translated := ""
	if _, ok := param.GetProperties()["translated"]; ok {
		translated = param.GetProperties()["translated"]
	}

	if _, ok := param.GetProperties()["lang"]; !ok {
		return nil, xerrors.Errorf("lang is not defined. err: %w", lib.ErrInvalidArgument)
	}

	m := &englishWordProblemUpdateParemeter{
		Text:       param.GetProperties()["text"],
		Translated: translated,
	}
	return m, lib.Validator.Struct(m)
}

type EnglishWordProblemProcessor interface {
	appS.ProblemAddProcessor
	appS.ProblemUpdateProcessor
	appS.ProblemRemoveProcessor
	appS.ProblemImportProcessor
	appS.ProblemQuotaProcessor
}

type englishWordProblemProcessor struct {
	synthesizer                     pluginS.Synthesizer
	translationClient               pluginS.TranslationClient
	tatoebaClient                   pluginS.TatoebaClient
	newProblemAddParameterCSVReader func(workbookID app.WorkbookID, reader io.Reader) appS.ProblemAddParameterIterator
}

func NewEnglishWordProblemProcessor(synthesizer pluginS.Synthesizer, translationClient pluginS.TranslationClient, tatoebaClient pluginS.TatoebaClient, newProblemAddParameterCSVReader func(workbookID app.WorkbookID, reader io.Reader) appS.ProblemAddParameterIterator) EnglishWordProblemProcessor {
	return &englishWordProblemProcessor{
		synthesizer:                     synthesizer,
		translationClient:               translationClient,
		tatoebaClient:                   tatoebaClient,
		newProblemAddParameterCSVReader: newProblemAddParameterCSVReader,
	}
}

func (p *englishWordProblemProcessor) AddProblem(ctx context.Context, rf appS.RepositoryFactory, operator app.StudentModel, workbook app.WorkbookModel, param appS.ProblemAddParameter) (appS.Added, app.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Debug("englishWordProblemProcessor.AddProblem, param: %+v", param)

	problemRepo, err := rf.NewProblemRepository(ctx, domain.EnglishWordProblemType)
	if err != nil {
		return 0, 0, xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	extractedParam, err := toEnglishWordProblemAddParemeter(param)
	if err != nil {
		return 0, 0, xerrors.Errorf("failed to toNewEnglishWordProblemParemeter. param: %+v, err: %w", param, err)
	}

	audioID := app.AudioID(0)
	if workbook.GetProperties()["audioEnabled"] == "true" {
		audioIDtmp, err := p.findOrAddAudio(ctx, rf, extractedParam.Text)
		if err != nil {
			return 0, 0, xerrors.Errorf("failed to p.findOrAddAudio. err: %w", err)
		}

		if audioIDtmp == 0 {
			return 0, 0, xerrors.Errorf("audio ID is zero. text: %s", extractedParam.Text)
		}

		audioID = audioIDtmp
	}

	if extractedParam.Translated == "" && extractedParam.Pos == plugin.PosOther {
		count, problemID, err := p.addMultipleProblem(ctx, operator, problemRepo, param, extractedParam, audioID)
		if err != nil {
			return 0, 0, xerrors.Errorf("failed to addMultipleProblem: err: %w", err)
		}

		return count, problemID, nil
	}

	problemID, err := p.addSingleProblem(ctx, operator, problemRepo, param, extractedParam, audioID)
	if err != nil {
		return 0, 0, xerrors.Errorf("failed to addSingleProblem: extractedParam: %+v, err: %w", extractedParam, err)
	}

	return 1, problemID, nil
}

func (p *englishWordProblemProcessor) addSingleProblem(ctx context.Context, operator app.StudentModel, problemRepo appS.ProblemRepository, param appS.ProblemAddParameter, extractedParam *englishWordProblemAddParemeter, audioID app.AudioID) (app.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem, text: %s, audio ID: %d", extractedParam.Text, audioID)

	if extractedParam.Translated == "" {
		translation, err := p.translateWithPos(ctx, extractedParam.Text, extractedParam.Pos, app.Lang2EN, app.Lang2JA)
		if err != nil {
			if errors.Is(err, pluginS.ErrTranslationNotFound) {
				extractedParam.Translated = ""
			} else {
				logger.Errorf("translate err: %v", err)
			}
		} else {
			extractedParam.Translated = translation.GetTranslated()
		}
	}

	properties := map[string]string{
		"lang":       extractedParam.Lang.String(),
		"text":       extractedParam.Text,
		"translated": extractedParam.Translated,
		"pos":        strconv.Itoa(int(extractedParam.Pos)),
		"audioId":    strconv.Itoa(int(audioID)),
	}
	newParam, err := appS.NewProblemAddParameter(param.GetWorkbookID(), param.GetNumber(), properties)
	if err != nil {
		return 0, xerrors.Errorf("failed to NewParameter. err: %w", err)
	}

	problemID, err := problemRepo.AddProblem(ctx, operator, newParam)
	if err != nil {
		return 0, xerrors.Errorf("failed to problemRepo.AddProblem. param: %+v, err: %w", param, err)
	}

	return problemID, nil
}

func (p *englishWordProblemProcessor) addMultipleProblem(ctx context.Context, operator app.StudentModel, problemRepo appS.ProblemRepository, param appS.ProblemAddParameter, extractedParam *englishWordProblemAddParemeter, audioID app.AudioID) (appS.Added, app.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Debugf("addMultipleProblem. text: %s, audio ID: %d", extractedParam.Text, audioID)

	translated, err := p.translate(ctx, extractedParam.Text, app.Lang2EN, app.Lang2JA)
	if err != nil {
		logger.Errorf("translate err: %v", err)
		properties := map[string]string{
			"text":       extractedParam.Text,
			"translated": extractedParam.Translated,
			"pos":        strconv.Itoa(int(extractedParam.Pos)),
			"audioId":    strconv.Itoa(int(audioID)),
			"lang":       app.Lang2JA.String(),
		}
		newParam, err := appS.NewProblemAddParameter(param.GetWorkbookID(), param.GetNumber(), properties)
		if err != nil {
			return 0, 0, xerrors.Errorf("failed to NewParameter. err: %w", err)
		}

		problemID, err := problemRepo.AddProblem(ctx, operator, newParam)
		if err != nil {
			return 0, 0, xerrors.Errorf("failed to problemRepo.AddProblem. param: %+v, err: %w", param, err)
		}

		return 1, problemID, nil
	}

	for _, t := range translated {
		properties := map[string]string{
			"text":       extractedParam.Text,
			"translated": t.GetTranslated(),
			"pos":        strconv.Itoa(int(t.GetPos())),
			"audioId":    strconv.Itoa(int(audioID)),
			"lang":       app.Lang2JA.String(),
		}
		newParam, err := appS.NewProblemAddParameter(param.GetWorkbookID(), param.GetNumber(), properties)
		if err != nil {
			return 0, 0, xerrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
		}

		if _, err := problemRepo.AddProblem(ctx, operator, newParam); err != nil {
			return 0, 0, xerrors.Errorf("failed to problemRepo.AddProblem. param: %+v, err: %w", param, err)
		}
	}

	return appS.Added(len(translated)), 0, nil
}

func (p *englishWordProblemProcessor) UpdateProblem(ctx context.Context, rf appS.RepositoryFactory, operator app.StudentModel, workbook app.WorkbookModel, id appS.ProblemSelectParameter2, param appS.ProblemUpdateParameter) (appS.Added, appS.Updated, error) {
	logger := log.FromContext(ctx)
	logger.Debug("englishWordProblemProcessor.UpdateProblem, param: %+v", param)

	problemRepo, err := rf.NewProblemRepository(ctx, domain.EnglishWordProblemType)
	if err != nil {
		return 0, 0, xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	extractedParam, err := toEnglishWordProblemUpdateParemeter(param)
	if err != nil {
		return 0, 0, xerrors.Errorf("failed to toNewEnglishWordProblemParemeter. param: %+v, err: %w", param, err)
	}

	audioID := app.AudioID(0)
	if workbook.GetProperties()["audioEnabled"] == "true" {
		audioIDtmp, err := p.findOrAddAudio(ctx, rf, extractedParam.Text)
		if err != nil {
			return 0, 0, xerrors.Errorf("failed to p.findOrAddAudio. err: %w", err)
		}

		if audioIDtmp == 0 {
			return 0, 0, xerrors.Errorf("audio ID is zero. text: %s", extractedParam.Text)
		}

		audioID = audioIDtmp
	}

	sentenceID := app.ProblemID(0)
	sentenceProvider := param.GetProperties()["sentenceProvider"]
	tatoebaSentenceNumberFromS := param.GetProperties()["tatoebaSentenceNumber1"]
	tatoebaSentenceNumberToS := param.GetProperties()["tatoebaSentenceNumber2"]
	if sentenceProvider == "tatoeba" {
		tatoebaSentenceNumberFrom, err := strconv.Atoi(tatoebaSentenceNumberFromS)
		if err != nil {
			return 0, 0, xerrors.Errorf("failed to Atoi. value: %s, err: %w", tatoebaSentenceNumberFromS, err)
		}
		tatoebaSentenceNumberTo, err := strconv.Atoi(tatoebaSentenceNumberToS)
		if err != nil {
			return 0, 0, xerrors.Errorf("failed to Atoi. value: %s, err: %w", tatoebaSentenceNumberToS, err)
		}

		sentenceIDtmp, err := p.findOrAddSentenceFromTatoeba(ctx, rf, operator, tatoebaSentenceNumberFrom, tatoebaSentenceNumberTo)
		if err != nil {
			return 0, 0, xerrors.Errorf("failed to findOrAddSentenceFromTatoeba. err: %w", err)
		}
		sentenceID = sentenceIDtmp
	}

	if err := p.updateSingleProblem(ctx, operator, problemRepo, id, param, extractedParam, audioID, sentenceID); err != nil {
		return 0, 0, xerrors.Errorf("failed to updateSingleProblem: extractedParam: %+v, err: %w", extractedParam, err)
	}

	return 1, 1, nil
}

func (p *englishWordProblemProcessor) updateSingleProblem(ctx context.Context, operator app.StudentModel, problemRepo appS.ProblemRepository, id appS.ProblemSelectParameter2, param appS.ProblemUpdateParameter, extractedParam *englishWordProblemUpdateParemeter, audioID app.AudioID, sentenceID1 app.ProblemID) error {
	logger := log.FromContext(ctx)
	logger.Infof("updateSingleProblem, text: %s, audio ID: %d", extractedParam.Text, audioID)

	properties := map[string]string{
		EnglishWordProblemUpdatePropertyText:        extractedParam.Text,
		EnglishWordProblemUpdatePropertyTranslated:  extractedParam.Translated,
		EnglishWordProblemUpdatePropertyAudioID:     strconv.Itoa(int(audioID)),
		EnglishWordProblemUpdatePropertySentenceID1: strconv.Itoa(int(sentenceID1)),
	}
	paramToUpdate, err := appS.NewProblemUpdateParameter(param.GetNumber(), properties)
	if err != nil {
		return xerrors.Errorf("failed to NewParameter. err: %w", err)
	}

	if err := problemRepo.UpdateProblem(ctx, operator, id, paramToUpdate); err != nil {
		return xerrors.Errorf("failed to problemRepo.UpdateProblem. param: %+v, err: %w", param, err)
	}

	return nil
}

func (p *englishWordProblemProcessor) RemoveProblem(ctx context.Context, rf appS.RepositoryFactory, operator app.StudentModel, id appS.ProblemSelectParameter2) error {
	problemRepo, err := rf.NewProblemRepository(ctx, domain.EnglishWordProblemType)
	if err != nil {
		return xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	if err := problemRepo.RemoveProblem(ctx, operator, id); err != nil {
		return err
	}

	return nil
}

func (p *englishWordProblemProcessor) CreateCSVReader(ctx context.Context, workbookID app.WorkbookID, reader io.Reader) (appS.ProblemAddParameterIterator, error) {
	return p.newProblemAddParameterCSVReader(workbookID, reader), nil
}

func (p *englishWordProblemProcessor) findOrAddAudio(ctx context.Context, rf appS.RepositoryFactory, text string) (app.AudioID, error) {
	audioRepo, err := rf.NewAudioRepository(ctx)
	if err != nil {
		return 0, err
	}

	{
		id, err := audioRepo.FindAudioIDByText(ctx, app.Lang5ENUS, text)
		if err != nil {
			if !errors.Is(err, appS.ErrAudioNotFound) {
				return 0, xerrors.Errorf("failed to FindAudioID. err: %w", err)
			}
		} else {
			return id, nil
		}
	}

	audioContent, err := p.synthesizer.Synthesize(app.Lang5ENUS, text)
	if err != nil {
		return 0, err
	}

	id, err := audioRepo.AddAudio(ctx, app.Lang5ENUS, text, audioContent)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (p *englishWordProblemProcessor) translateWithPos(ctx context.Context, text string, pos plugin.WordPos, fromLang, toLang app.Lang2) (plugin.Translation, error) {
	logger := log.FromContext(ctx)
	logger.Infof("translateWithPos. text: %s", text)

	return p.translationClient.DictionaryLookupWithPos(ctx, fromLang, toLang, text, pos)
}

func (p *englishWordProblemProcessor) translate(ctx context.Context, text string, fromLang, toLang app.Lang2) ([]plugin.Translation, error) {
	logger := log.FromContext(ctx)
	logger.Infof("translate. text: %s", text)

	return p.translationClient.DictionaryLookup(ctx, fromLang, toLang, text)
}

func (p *englishWordProblemProcessor) GetUnitForSizeQuota() appS.QuotaUnit {
	return quotaSizeUnit
}

func (p *englishWordProblemProcessor) GetLimitForSizeQuota() int {
	return quotaSizeLimit
}

func (p *englishWordProblemProcessor) GetUnitForUpdateQuota() appS.QuotaUnit {
	return quotaUpdateUnit
}

func (p *englishWordProblemProcessor) GetLimitForUpdateQuota() int {
	return quotaUpdateLimit
}

func (p *englishWordProblemProcessor) findOrAddSentenceFromTatoeba(ctx context.Context, rf appS.RepositoryFactory, operator app.StudentModel, tatoebaSentenceNumberFrom, tatoebaSentenceNumberTo int) (app.ProblemID, error) {
	systemSpaceID := appS.GetSystemSpaceID()
	workbookRepo, err := rf.NewWorkbookRepository(ctx)
	if err != nil {
		return 0, xerrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

	tatoebaWorkbook, err := workbookRepo.FindWorkbookByName(ctx, operator, systemSpaceID, appS.TatoebaWorkbookName)
	if err != nil {
		return 0, xerrors.Errorf("failed to FindWorkbookByName. name: %s, err: %w", appS.TatoebaWorkbookName, err)
	}

	tatoebaSentenceFrom, err := p.tatoebaClient.FindSentenceBySentenceNumber(ctx, tatoebaSentenceNumberFrom)
	if err != nil {
		return 0, xerrors.Errorf("failed to FindTatoebaSentenceBySentenceNumber. err: %w", err)

	}

	tatoebaSentenceTo, err := p.tatoebaClient.FindSentenceBySentenceNumber(ctx, tatoebaSentenceNumberTo)
	if err != nil {
		return 0, xerrors.Errorf("failed to FindTatoebaSentenceBySentenceNumber. err: %w", err)
	}

	sentenceProblemRepo, err := rf.NewProblemRepository(ctx, domain.EnglishSentenceProblemType)
	if err != nil {
		return 0, xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}
	condition := map[string]interface{}{
		"workbookId": tatoebaWorkbook.GetID(),
		"text":       tatoebaSentenceFrom.GetText(),
		"translated": tatoebaSentenceTo.GetText(),
	}
	problems, err := sentenceProblemRepo.FindProblemsByCustomCondition(ctx, operator, condition)
	if err != nil {
		return 0, xerrors.Errorf("failed to FindProblemsByCustomCondition. err: %w", err)
	}

	if len(problems) != 0 {
		return app.ProblemID(problems[0].GetID()), nil
	}

	sentenceAddParam := englishSentenceProblemAddParemeter{
		Lang:                   app.Lang2JA,
		Text:                   tatoebaSentenceFrom.GetText(),
		Translated:             tatoebaSentenceTo.GetText(),
		Provider:               "tatoeba",
		TatoebaSentenceNumber1: tatoebaSentenceFrom.GetSentenceNumber(),
		TatoebaSentenceNumber2: tatoebaSentenceTo.GetSentenceNumber(),
		TatoebaAuthor1:         tatoebaSentenceFrom.GetAuthor(),
		TatoebaAuthor2:         tatoebaSentenceTo.GetAuthor(),
	}
	sentenceAddProperties := sentenceAddParam.toProperties(0)

	param, err := appS.NewProblemAddParameter(app.WorkbookID(tatoebaWorkbook.GetID()), 1, sentenceAddProperties)
	if err != nil {
		return 0, xerrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
	}
	id, err := sentenceProblemRepo.AddProblem(ctx, operator, param)
	if err != nil {
		return 0, xerrors.Errorf("failed to AddProblem. err: %w", err)
	}

	return id, err
}
