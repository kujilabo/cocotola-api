package domain

import (
	"context"
	"errors"
	"io"
	"strconv"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	plugin "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	"golang.org/x/xerrors"
)

var (
	quotaSizeUnit  = app.QuotaUnitPersitance
	quotaSizeLimit = 5000

	quotaUpdateUnit  = app.QuotaUnitDay
	quotaUpdateLimit = 100
)

type englishWordProblemAddParemeter struct {
	Lang       app.Lang2      `validate:"required"`
	Text       string         `validate:"required"`
	Pos        plugin.WordPos `validate:"required"`
	Translated string
}

func toEnglishWordProblemAddParemeter(param app.ProblemAddParameter) (*englishWordProblemAddParemeter, error) {
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

func toEnglishWordProblemUpdateParemeter(param app.ProblemUpdateParameter) (*englishWordProblemUpdateParemeter, error) {
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
	app.ProblemAddProcessor
	app.ProblemUpdateProcessor
	app.ProblemRemoveProcessor
	app.ProblemImportProcessor
	app.ProblemQuotaProcessor
}

type englishWordProblemProcessor struct {
	synthesizer                     plugin.Synthesizer
	translator                      plugin.Translator
	newProblemAddParameterCSVReader func(workbookID app.WorkbookID, reader io.Reader) app.ProblemAddParameterIterator
}

func NewEnglishWordProblemProcessor(synthesizer plugin.Synthesizer, translator plugin.Translator, newProblemAddParameterCSVReader func(workbookID app.WorkbookID, reader io.Reader) app.ProblemAddParameterIterator) EnglishWordProblemProcessor {
	return &englishWordProblemProcessor{
		synthesizer:                     synthesizer,
		translator:                      translator,
		newProblemAddParameterCSVReader: newProblemAddParameterCSVReader,
	}
}

func (p *englishWordProblemProcessor) AddProblem(ctx context.Context, repo app.RepositoryFactory, operator app.Student, workbook app.Workbook, param app.ProblemAddParameter) (app.Added, app.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Debug("englishWordProblemProcessor.AddProblem, param: %+v", param)

	problemRepo, err := repo.NewProblemRepository(ctx, workbook.GetProblemType())
	if err != nil {
		return 0, 0, xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	extractedParam, err := toEnglishWordProblemAddParemeter(param)
	if err != nil {
		return 0, 0, xerrors.Errorf("failed to toNewEnglishWordProblemParemeter. param: %+v, err: %w", param, err)
	}

	audioID := app.AudioID(0)
	if workbook.GetProperties()["audioEnabled"] == "true" {
		audioIDtmp, err := p.findOrAddAudio(ctx, repo, extractedParam.Text)
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

func (p *englishWordProblemProcessor) addSingleProblem(ctx context.Context, operator app.Student, problemRepo app.ProblemRepository, param app.ProblemAddParameter, extractedParam *englishWordProblemAddParemeter, audioID app.AudioID) (app.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem, text: %s, audio ID: %d", extractedParam.Text, audioID)

	if extractedParam.Translated == "" {
		translation, err := p.translateWithPos(ctx, extractedParam.Text, extractedParam.Pos, app.Lang2EN, app.Lang2JA)
		if err != nil {
			if errors.Is(err, plugin.ErrTranslationNotFound) {
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
	newParam, err := app.NewProblemAddParameter(param.GetWorkbookID(), param.GetNumber(), properties)
	if err != nil {
		return 0, xerrors.Errorf("failed to NewParameter. err: %w", err)
	}

	problemID, err := problemRepo.AddProblem(ctx, operator, newParam)
	if err != nil {
		return 0, xerrors.Errorf("failed to problemRepo.AddProblem. param: %+v, err: %w", param, err)
	}

	return problemID, nil
}

func (p *englishWordProblemProcessor) addMultipleProblem(ctx context.Context, operator app.Student, problemRepo app.ProblemRepository, param app.ProblemAddParameter, extractedParam *englishWordProblemAddParemeter, audioID app.AudioID) (app.Added, app.ProblemID, error) {
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
		newParam, err := app.NewProblemAddParameter(param.GetWorkbookID(), param.GetNumber(), properties)
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
		newParam, err := app.NewProblemAddParameter(param.GetWorkbookID(), param.GetNumber(), properties)
		if err != nil {
			return 0, 0, xerrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
		}

		if _, err := problemRepo.AddProblem(ctx, operator, newParam); err != nil {
			return 0, 0, xerrors.Errorf("failed to problemRepo.AddProblem. param: %+v, err: %w", param, err)
		}
	}

	return app.Added(len(translated)), 0, nil
}

func (p *englishWordProblemProcessor) UpdateProblem(ctx context.Context, repo app.RepositoryFactory, operator app.Student, workbook app.Workbook, id app.ProblemSelectParameter2, param app.ProblemUpdateParameter) (app.Added, app.Updated, error) {
	logger := log.FromContext(ctx)
	logger.Debug("englishWordProblemProcessor.UpdateProblem, param: %+v", param)

	problemRepo, err := repo.NewProblemRepository(ctx, workbook.GetProblemType())
	if err != nil {
		return 0, 0, xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	extractedParam, err := toEnglishWordProblemUpdateParemeter(param)
	if err != nil {
		return 0, 0, xerrors.Errorf("failed to toNewEnglishWordProblemParemeter. param: %+v, err: %w", param, err)
	}

	audioID := app.AudioID(0)
	if workbook.GetProperties()["audioEnabled"] == "true" {
		audioIDtmp, err := p.findOrAddAudio(ctx, repo, extractedParam.Text)
		if err != nil {
			return 0, 0, xerrors.Errorf("failed to p.findOrAddAudio. err: %w", err)
		}

		if audioIDtmp == 0 {
			return 0, 0, xerrors.Errorf("audio ID is zero. text: %s", extractedParam.Text)
		}

		audioID = audioIDtmp
	}

	if err := p.updateSingleProblem(ctx, operator, problemRepo, id, param, extractedParam, audioID); err != nil {
		return 0, 0, xerrors.Errorf("failed to updateSingleProblem: extractedParam: %+v, err: %w", extractedParam, err)
	}

	return 1, 1, nil
}

func (p *englishWordProblemProcessor) updateSingleProblem(ctx context.Context, operator app.Student, problemRepo app.ProblemRepository, id app.ProblemSelectParameter2, param app.ProblemUpdateParameter, extractedParam *englishWordProblemUpdateParemeter, audioID app.AudioID) error {
	logger := log.FromContext(ctx)
	logger.Infof("updateSingleProblem, text: %s, audio ID: %d", extractedParam.Text, audioID)

	properties := map[string]string{
		"text":       extractedParam.Text,
		"translated": extractedParam.Translated,
		"audioId":    strconv.Itoa(int(audioID)),
	}
	paramToUpdate, err := app.NewProblemUpdateParameter(param.GetNumber(), properties)
	if err != nil {
		return xerrors.Errorf("failed to NewParameter. err: %w", err)
	}

	if err := problemRepo.UpdateProblem(ctx, operator, id, paramToUpdate); err != nil {
		return xerrors.Errorf("failed to problemRepo.UpdateProblem. param: %+v, err: %w", param, err)
	}

	return nil
}

func (p *englishWordProblemProcessor) RemoveProblem(ctx context.Context, repo app.RepositoryFactory, operator app.Student, id app.ProblemSelectParameter2) error {
	problemRepo, err := repo.NewProblemRepository(ctx, EnglishWordProblemType)
	if err != nil {
		return xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	if err := problemRepo.RemoveProblem(ctx, operator, id); err != nil {
		return err
	}

	return nil
}

func (p *englishWordProblemProcessor) CreateCSVReader(ctx context.Context, workbookID app.WorkbookID, reader io.Reader) (app.ProblemAddParameterIterator, error) {
	return p.newProblemAddParameterCSVReader(workbookID, reader), nil
}

func (p *englishWordProblemProcessor) findOrAddAudio(ctx context.Context, repo app.RepositoryFactory, text string) (app.AudioID, error) {
	audioRepo, err := repo.NewAudioRepository(ctx)
	if err != nil {
		return 0, err
	}

	{
		id, err := audioRepo.FindAudioIDByText(ctx, app.Lang5ENUS, text)
		if err != nil {
			if !errors.Is(err, app.ErrAudioNotFound) {
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

	return p.translator.DictionaryLookupWithPos(ctx, fromLang, toLang, text, pos)
}

func (p *englishWordProblemProcessor) translate(ctx context.Context, text string, fromLang, toLang app.Lang2) ([]plugin.Translation, error) {
	logger := log.FromContext(ctx)
	logger.Infof("translate. text: %s", text)

	return p.translator.DictionaryLookup(ctx, fromLang, toLang, text)
}

func (p *englishWordProblemProcessor) GetUnitForSizeQuota() app.QuotaUnit {
	return quotaSizeUnit
}

func (p *englishWordProblemProcessor) GetLimitForSizeQuota() int {
	return quotaSizeLimit
}

func (p *englishWordProblemProcessor) GetUnitForUpdateQuota() app.QuotaUnit {
	return quotaUpdateUnit
}

func (p *englishWordProblemProcessor) GetLimitForUpdateQuota() int {
	return quotaUpdateLimit
}
