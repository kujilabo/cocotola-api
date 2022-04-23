package service

import (
	"context"
	"errors"
	"strconv"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	appS "github.com/kujilabo/cocotola-api/pkg_app/service"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	pluginS "github.com/kujilabo/cocotola-api/pkg_plugin/common/service"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
	"golang.org/x/xerrors"
)

type englishPhraseProblemAddParemeter struct {
	Text       string `validate:"required"`
	Lang       string `validate:"required"`
	Translated string `validate:"required"`
}

func toEnglishPhraseProblemAddParemeter(param appS.ProblemAddParameter) (*englishPhraseProblemAddParemeter, error) {
	if _, ok := param.GetProperties()["lang"]; !ok {
		return nil, xerrors.Errorf("lang is not defined. err: %w", lib.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["text"]; !ok {
		return nil, xerrors.Errorf("text is not defined. err: %w", lib.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["translated"]; !ok {
		return nil, xerrors.Errorf("translated is not defined. err: %w", lib.ErrInvalidArgument)
	}

	m := &englishPhraseProblemAddParemeter{
		Lang:       param.GetProperties()["lang"],
		Text:       param.GetProperties()["text"],
		Translated: param.GetProperties()["translated"],
	}

	return m, lib.Validator.Struct(m)
}

type EnglishPhraseProblemProcessor interface {
	appS.ProblemAddProcessor
	appS.ProblemRemoveProcessor
}

type englishPhraseProblemProcessor struct {
	synthesizerClient pluginS.SynthesizerClient
	translatorClient  pluginS.TranslatorClient
}

func NewEnglishPhraseProblemProcessor(synthesizerClient pluginS.SynthesizerClient, translatorClient pluginS.TranslatorClient) EnglishPhraseProblemProcessor {
	return &englishPhraseProblemProcessor{
		synthesizerClient: synthesizerClient,
		translatorClient:  translatorClient,
	}
}

func (p *englishPhraseProblemProcessor) AddProblem(ctx context.Context, repo appS.RepositoryFactory, operator app.StudentModel, workbook app.WorkbookModel, param appS.ProblemAddParameter) ([]app.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem1")

	problemRepo, err := repo.NewProblemRepository(ctx, domain.EnglishPhraseProblemType)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	extractedParam, err := toEnglishPhraseProblemAddParemeter(param)
	if err != nil {
		return nil, xerrors.Errorf("failed to toNewEnglishPhraseProblemParemeter. err: %w", err)
	}

	audioID, err := p.findOrAddAudio(ctx, repo, extractedParam.Text)
	if err != nil {
		return nil, xerrors.Errorf("failed to p.findOrAddAudio. err: %w", err)
	}

	if audioID == 0 {
		return nil, xerrors.Errorf("audio ID is zero. text: %s", extractedParam.Text)
	}

	problemID, err := p.addSingleProblem(ctx, operator, problemRepo, param, extractedParam, audioID)
	if err != nil {
		return nil, xerrors.Errorf("failed to addSingleProblem: extractedParam: %+v, err: %w", extractedParam, err)
	}

	return []app.ProblemID{problemID}, err
}

func (p *englishPhraseProblemProcessor) addSingleProblem(ctx context.Context, operator app.StudentModel, problemRepo appS.ProblemRepository, param appS.ProblemAddParameter, extractedParam *englishPhraseProblemAddParemeter, audioID app.AudioID) (app.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem1")

	logger.Infof("text: %s, audio ID: %d", extractedParam.Text, audioID)

	properties := map[string]string{
		"text":       extractedParam.Text,
		"translated": extractedParam.Translated,
		"lang":       extractedParam.Lang,
		"audioId":    strconv.Itoa(int(audioID)),
	}
	newParam, err := appS.NewProblemAddParameter(param.GetWorkbookID(), param.GetNumber(), properties)
	if err != nil {
		return 0, xerrors.Errorf("failed to NewParameter. err: %w", err)
	}

	problemID, err := problemRepo.AddProblem(ctx, operator, newParam)
	if err != nil {
		return 0, xerrors.Errorf("failed to problemRepo.AddProblem. err: %w", err)
	}

	return problemID, nil

}

func (p *englishPhraseProblemProcessor) RemoveProblem(ctx context.Context, repo appS.RepositoryFactory, operator app.StudentModel, id appS.ProblemSelectParameter2) error {
	problemRepo, err := repo.NewProblemRepository(ctx, domain.EnglishPhraseProblemType)
	if err != nil {
		return xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	if err := problemRepo.RemoveProblem(ctx, operator, id); err != nil {
		return err
	}

	return nil
}

func (p *englishPhraseProblemProcessor) findOrAddAudio(ctx context.Context, repo appS.RepositoryFactory, text string) (app.AudioID, error) {
	audioRepo := repo.NewAudioRepository(ctx)

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

	audioContent, err := p.synthesizerClient.Synthesize(ctx, app.Lang5ENUS, text)
	if err != nil {
		return 0, err
	}

	id, err := audioRepo.AddAudio(ctx, app.Lang5ENUS, text, audioContent)
	if err != nil {
		return 0, err
	}

	return id, err
}
