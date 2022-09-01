package service

import (
	"context"
	"errors"
	"io"
	"strconv"

	appD "github.com/kujilabo/cocotola-api/src/app/domain"
	appS "github.com/kujilabo/cocotola-api/src/app/service"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	liberrors "github.com/kujilabo/cocotola-api/src/lib/errors"
	"github.com/kujilabo/cocotola-api/src/lib/log"
	plugin "github.com/kujilabo/cocotola-api/src/plugin/common/domain"
	pluginS "github.com/kujilabo/cocotola-api/src/plugin/common/service"
	"github.com/kujilabo/cocotola-api/src/plugin/english/domain"
)

var (
	EnglishWordProblemQuotaSizeUnit            = appS.QuotaUnitPersitance
	EnglishWordProblemQuotaSizeLimit           = 5000
	EnglishWordProblemQuotaUpdateUnit          = appS.QuotaUnitDay
	EnglishWordProblemQuotaUpdateLimit         = 100
	EnglishWordProblemUpdatePropertyText       = "text"
	EnglishWordProblemUpdatePropertyTranslated = "translated"
	EnglishWordProblemUpdatePropertyAudioID    = "audioId"
	// EnglishWordProblemUpdatePropertyTatoebaSentenceNumber1 = "tatoebaSentenceNumber1"
	// EnglishWordProblemUpdatePropertyTatoebaSentenceNumber2 = "tatoebaSentenceNumber2"
	EnglishWordProblemUpdatePropertySentenceID1 = "sentenceId1"

	EnglishWordProblemAddPropertyAudioID    = "audioId"
	EnglishWordProblemAddPropertyLang2      = "lang2"
	EnglishWordProblemAddPropertyText       = "text"
	EnglishWordProblemAddPropertyTranslated = "translated"
	EnglishWordProblemAddPropertyPos        = "pos"
)

type EnglishWordProblemAddParemeter struct {
	Lang2      appD.Lang2     `validate:"required"`
	Text       string         `validate:"required"`
	Pos        plugin.WordPos `validate:"required"`
	Translated string
}

func (p *EnglishWordProblemAddParemeter) toProperties() map[string]string {
	return map[string]string{
		// EnglishWordProblemAddPropertyAudioID:    strconv.Itoa(int(uint(audioID))),
		EnglishWordProblemAddPropertyLang2:      p.Lang2.String(),
		EnglishWordProblemAddPropertyText:       p.Text,
		EnglishWordProblemAddPropertyTranslated: p.Translated,
		EnglishWordProblemAddPropertyPos:        strconv.Itoa(int(p.Pos)),
	}
}

func NewEnglishWordProblemAddParemeter(param appS.ProblemAddParameter) (*EnglishWordProblemAddParemeter, error) {
	posS := param.GetProperties()["pos"]
	pos, err := strconv.Atoi(posS)
	if err != nil {
		return nil, liberrors.Errorf("failed to cast to int. err: %w", libD.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["text"]; !ok {
		return nil, liberrors.Errorf("text is not defined. err: %w", libD.ErrInvalidArgument)
	}

	translated := ""
	if _, ok := param.GetProperties()["translated"]; ok {
		translated = param.GetProperties()["translated"]
	}

	if _, ok := param.GetProperties()["lang2"]; !ok {
		return nil, liberrors.Errorf("lang2 is not defined. err: %w", libD.ErrInvalidArgument)
	}

	lang2, err := appD.NewLang2(param.GetProperties()["lang2"])
	if err != nil {
		return nil, liberrors.Errorf("lang2 format is invalid. err: %w", err)
	}

	m := &EnglishWordProblemAddParemeter{
		Lang2:      lang2,
		Text:       param.GetProperties()["text"],
		Pos:        plugin.WordPos(pos),
		Translated: translated,
	}
	return m, libD.Validator.Struct(m)
}

type EnglishWordProblemUpdateParemeter struct {
	Lang2                     appD.Lang2 `validate:"required"`
	Text                      string     `validate:"required"`
	Translated                string
	SentenceProvider          string
	TatoebaSentenceNumberFrom int
	TatoebaSentenceNumberTo   int
	// sentenceProvider := param.GetProperties()["sentenceProvider"]
	// tatoebaSentenceNumberFromS := param.GetProperties()["tatoebaSentenceNumber1"]
	// tatoebaSentenceNumberToS := param.GetProperties()["tatoebaSentenceNumber2"]
}

func NewEnglishWordProblemUpdateParemeter(param appS.ProblemUpdateParameter) (*EnglishWordProblemUpdateParemeter, error) {
	if _, ok := param.GetProperties()["text"]; !ok {
		return nil, liberrors.Errorf("text is not defined. err: %w", libD.ErrInvalidArgument)
	}

	translated := ""
	if _, ok := param.GetProperties()["translated"]; ok {
		translated = param.GetProperties()["translated"]
	}

	if _, ok := param.GetProperties()["lang2"]; !ok {
		return nil, liberrors.Errorf("lang2 is not defined. err: %w", libD.ErrInvalidArgument)
	}

	lang2, err := appD.NewLang2(param.GetProperties()["lang2"])
	if err != nil {
		return nil, liberrors.Errorf("lang2 format is invalid. err: %w", err)
	}

	tatoebaSentenceNumberFrom := 0
	tatoebaSentenceNumberTo := 0

	sentenceProvider := param.GetProperties()["sentenceProvider"]
	if sentenceProvider == "tatoeba" {
		tatoebaSentenceNumberFromS := param.GetProperties()["tatoebaSentenceNumber1"]

		from, err := strconv.Atoi(tatoebaSentenceNumberFromS)
		if err != nil {
			return nil, liberrors.Errorf("failed to Atoi. value: %s, err: %w", tatoebaSentenceNumberFromS, err)
		}

		tatoebaSentenceNumberToS := param.GetProperties()["tatoebaSentenceNumber2"]
		to, err := strconv.Atoi(tatoebaSentenceNumberToS)
		if err != nil {
			return nil, liberrors.Errorf("failed to Atoi. value: %s, err: %w", tatoebaSentenceNumberToS, err)
		}
		tatoebaSentenceNumberFrom = from
		tatoebaSentenceNumberTo = to
	}

	m := &EnglishWordProblemUpdateParemeter{
		Lang2:                     lang2,
		Text:                      param.GetProperties()["text"],
		Translated:                translated,
		SentenceProvider:          sentenceProvider,
		TatoebaSentenceNumberFrom: tatoebaSentenceNumberFrom,
		TatoebaSentenceNumberTo:   tatoebaSentenceNumberTo,
	}
	return m, libD.Validator.Struct(m)
}

type EnglishWordProblemProcessor interface {
	appS.ProblemAddProcessor
	appS.ProblemUpdateProcessor
	appS.ProblemRemoveProcessor
	appS.ProblemImportProcessor
	appS.ProblemQuotaProcessor
}

type englishWordProblemProcessor struct {
	synthesizerClient               appS.SynthesizerClient
	translatorClient                pluginS.TranslatorClient
	tatoebaClient                   pluginS.TatoebaClient
	newProblemAddParameterCSVReader func(workbookID appD.WorkbookID, reader io.Reader) appS.ProblemAddParameterIterator
}

func NewEnglishWordProblemProcessor(synthesizerClient appS.SynthesizerClient, translatorClient pluginS.TranslatorClient, tatoebaClient pluginS.TatoebaClient, newProblemAddParameterCSVReader func(workbookID appD.WorkbookID, reader io.Reader) appS.ProblemAddParameterIterator) EnglishWordProblemProcessor {
	return &englishWordProblemProcessor{
		synthesizerClient:               synthesizerClient,
		translatorClient:                translatorClient,
		tatoebaClient:                   tatoebaClient,
		newProblemAddParameterCSVReader: newProblemAddParameterCSVReader,
	}
}

func (p *englishWordProblemProcessor) AddProblem(ctx context.Context, rf appS.RepositoryFactory, operator appD.StudentModel, workbook appD.WorkbookModel, param appS.ProblemAddParameter) ([]appD.ProblemID, error) {
	ctx, span := tracer.Start(ctx, "englishWordProblemProcessor.AddProblem")
	defer span.End()

	logger := log.FromContext(ctx)
	logger.Debug("englishWordProblemProcessor.AddProblem, param: %+v", param)

	extractedParam, err := NewEnglishWordProblemAddParemeter(param)
	if err != nil {
		return nil, liberrors.Errorf("failed to toNewEnglishWordProblemParemeter. param: %+v, err: %w", param, err)
	}

	audioID := appD.AudioID(0)
	if workbook.GetProperties()["audioEnabled"] == "true" {
		audio, err := p.synthesizerClient.Synthesize(ctx, appD.Lang2EN, extractedParam.Text)
		if err != nil {
			return nil, err
		}

		audioID = appD.AudioID(audio.GetAudioModel().GetID())
	}

	logger.Debug("audioID: %d", audioID)

	var converter ToEnglishWordProblemAddParameter
	if extractedParam.Translated == "" && extractedParam.Pos == plugin.PosOther {
		converter = NewToMultipleEnglishWordProblemAddParameter(p.translatorClient, param.GetWorkbookID(), param.GetNumber(), extractedParam, audioID)
	} else {
		converter = NewToSingleEnglishWordProblemAddParameter(p.translatorClient, param.GetWorkbookID(), param.GetNumber(), extractedParam, audioID)
	}

	toAddParams, err := converter.Run(ctx)
	if err != nil {
		if errors.Is(err, pluginS.ErrTranslationNotFound) {
			message := "Translation not found"
			return nil, appD.NewPluginError("client", message, []string{message}, err)
		}
		return nil, err
	}

	problemRepo, err := rf.NewProblemRepository(ctx, domain.EnglishWordProblemType)
	if err != nil {
		return nil, liberrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	idsOfAddedProblem := make([]appD.ProblemID, len(toAddParams))
	for i, toAddParam := range toAddParams {
		problemID, err := problemRepo.AddProblem(ctx, operator, toAddParam)
		if err != nil {
			return nil, liberrors.Errorf("failed to problemRepo.AddProblem. param: %+v, err: %w", param, err)
		}
		idsOfAddedProblem[i] = problemID
	}

	return idsOfAddedProblem, nil
}

func (p *englishWordProblemProcessor) UpdateProblem(ctx context.Context, rf appS.RepositoryFactory, operator appD.StudentModel, workbook appD.WorkbookModel, id appS.ProblemSelectParameter2, param appS.ProblemUpdateParameter) (appS.Added, appS.Updated, error) {
	logger := log.FromContext(ctx)
	logger.Debugf("englishWordProblemProcessor.UpdateProblem, param: %+v", param)

	extractedParam, err := NewEnglishWordProblemUpdateParemeter(param)
	if err != nil {
		logger.Warnf("err: %+v", err)
		message := "Invalid parameter"
		return 0, 0, liberrors.Errorf("failed to toNewEnglishWordProblemParemeter. param: %+v, err: %w", param, appD.NewPluginError(appD.ErrorType(appD.ErrorTypeClient), message, []string{message, err.Error()}, err))
	}

	audioID := appD.AudioID(0)
	if workbook.GetProperties()["audioEnabled"] == "true" {
		audio, err := p.synthesizerClient.Synthesize(ctx, appD.Lang2EN, extractedParam.Text)
		if err != nil {
			return 0, 0, err
		}

		audioID = appD.AudioID(audio.GetAudioModel().GetID())
	}

	sentenceID := appD.ProblemID(0)
	if extractedParam.SentenceProvider == "tatoeba" {
		sentenceIDtmp, err := p.findOrAddSentenceFromTatoeba(ctx, rf, operator, extractedParam.TatoebaSentenceNumberFrom, extractedParam.TatoebaSentenceNumberTo, extractedParam.Lang2)
		if err != nil {
			return 0, 0, liberrors.Errorf("failed to findOrAddSentenceFromTatoeba. err: %w", err)
		}
		sentenceID = sentenceIDtmp
	}

	converter := NewToSingleEnglishWordProblemUpdateParameter(p.translatorClient, param.GetNumber(), extractedParam, audioID, sentenceID)
	toUpdateParams, err := converter.Run(ctx)
	if err != nil {
		return 0, 0, err
	}

	problemRepo, err := rf.NewProblemRepository(ctx, domain.EnglishWordProblemType)
	if err != nil {
		return 0, 0, liberrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	for _, toUpdateParam := range toUpdateParams {
		if err := problemRepo.UpdateProblem(ctx, operator, id, toUpdateParam); err != nil {
			return 0, 0, liberrors.Errorf("failed to problemRepo.UpdateProblem. param: %+v, err: %w", param, err)
		}
	}

	return 1, 1, nil
}

func (p *englishWordProblemProcessor) RemoveProblem(ctx context.Context, rf appS.RepositoryFactory, operator appD.StudentModel, id appS.ProblemSelectParameter2) error {
	problemRepo, err := rf.NewProblemRepository(ctx, domain.EnglishWordProblemType)
	if err != nil {
		return liberrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	if err := problemRepo.RemoveProblem(ctx, operator, id); err != nil {
		return err
	}

	return nil
}

func (p *englishWordProblemProcessor) CreateCSVReader(ctx context.Context, workbookID appD.WorkbookID, reader io.Reader) (appS.ProblemAddParameterIterator, error) {
	return p.newProblemAddParameterCSVReader(workbookID, reader), nil
}

func (p *englishWordProblemProcessor) GetUnitForSizeQuota() appS.QuotaUnit {
	return EnglishWordProblemQuotaSizeUnit
}

func (p *englishWordProblemProcessor) GetLimitForSizeQuota() int {
	return EnglishWordProblemQuotaSizeLimit
}

func (p *englishWordProblemProcessor) GetUnitForUpdateQuota() appS.QuotaUnit {
	return EnglishWordProblemQuotaUpdateUnit
}

func (p *englishWordProblemProcessor) GetLimitForUpdateQuota() int {
	return EnglishWordProblemQuotaUpdateLimit
}

func (p *englishWordProblemProcessor) findOrAddSentenceFromTatoeba(ctx context.Context, rf appS.RepositoryFactory, operator appD.StudentModel, tatoebaSentenceNumberFrom, tatoebaSentenceNumberTo int, lang2 appD.Lang2) (appD.ProblemID, error) {
	systemSpaceID := appS.GetSystemSpaceID()
	workbookRepo, err := rf.NewWorkbookRepository(ctx)
	if err != nil {
		return 0, liberrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

	tatoebaWorkbook, err := workbookRepo.FindWorkbookByName(ctx, operator, systemSpaceID, appS.TatoebaWorkbookName)
	if err != nil {
		return 0, liberrors.Errorf("failed to FindWorkbookByName. name: %s, err: %w", appS.TatoebaWorkbookName, err)
	}

	tatoebaSentenceFrom, err := p.tatoebaClient.FindSentenceBySentenceNumber(ctx, tatoebaSentenceNumberFrom)
	if err != nil {
		return 0, liberrors.Errorf("failed to FindTatoebaSentenceBySentenceNumber. err: %w", err)
	}

	tatoebaSentenceTo, err := p.tatoebaClient.FindSentenceBySentenceNumber(ctx, tatoebaSentenceNumberTo)
	if err != nil {
		return 0, liberrors.Errorf("failed to FindTatoebaSentenceBySentenceNumber. err: %w", err)
	}

	sentenceProblemRepo, err := rf.NewProblemRepository(ctx, domain.EnglishSentenceProblemType)
	if err != nil {
		return 0, liberrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}
	condition := map[string]interface{}{
		"workbookId": tatoebaWorkbook.GetID(),
		"text":       tatoebaSentenceFrom.GetText(),
		"translated": tatoebaSentenceTo.GetText(),
	}
	problems, err := sentenceProblemRepo.FindProblemsByCustomCondition(ctx, operator, condition)
	if err != nil {
		return 0, liberrors.Errorf("failed to FindProblemsByCustomCondition. err: %w", err)
	}

	if len(problems) != 0 {
		return appD.ProblemID(problems[0].GetID()), nil
	}

	sentenceAddParam := englishSentenceProblemAddParemeter{
		Lang2:                  lang2,
		Text:                   tatoebaSentenceFrom.GetText(),
		Translated:             tatoebaSentenceTo.GetText(),
		Provider:               "tatoeba",
		TatoebaSentenceNumber1: tatoebaSentenceFrom.GetSentenceNumber(),
		TatoebaSentenceNumber2: tatoebaSentenceTo.GetSentenceNumber(),
		TatoebaAuthor1:         tatoebaSentenceFrom.GetAuthor(),
		TatoebaAuthor2:         tatoebaSentenceTo.GetAuthor(),
	}
	sentenceAddProperties := sentenceAddParam.toProperties(0)

	param, err := appS.NewProblemAddParameter(appD.WorkbookID(tatoebaWorkbook.GetID()), 1, sentenceAddProperties)
	if err != nil {
		return 0, liberrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
	}
	id, err := sentenceProblemRepo.AddProblem(ctx, operator, param)
	if err != nil {
		return 0, liberrors.Errorf("failed to AddProblem. err: %w", err)
	}

	return id, err
}
