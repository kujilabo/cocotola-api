package domain

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	plugin "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
)

type englishPhraseProblemAddParemeter struct {
	Text       string `validate:"required"`
	Lang       string `validate:"required"`
	Translated string `validate:"required"`
}

func toEnglishPhraseProblemAddParemeter(param app.ProblemAddParameter) (*englishPhraseProblemAddParemeter, error) {
	if _, ok := param.GetProperties()["lang"]; !ok {
		return nil, fmt.Errorf("lang is not defined. err: %w", lib.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["text"]; !ok {
		return nil, fmt.Errorf("text is not defined. err: %w", lib.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["translated"]; !ok {
		return nil, fmt.Errorf("translated is not defined. err: %w", lib.ErrInvalidArgument)
	}

	m := &englishPhraseProblemAddParemeter{
		Lang:       param.GetProperties()["lang"],
		Text:       param.GetProperties()["text"],
		Translated: param.GetProperties()["translated"],
	}

	return m, lib.Validator.Struct(m)
}

type EnglishPhraseProblemProcessor interface {
	app.ProblemAddProcessor
	app.ProblemRemoveProcessor
}

type englishPhraseProblemProcessor struct {
	synthesizer plugin.Synthesizer
	translator  plugin.Translator
}

func NewEnglishPhraseProblemProcessor(synthesizer plugin.Synthesizer, translator plugin.Translator) EnglishPhraseProblemProcessor {
	return &englishPhraseProblemProcessor{
		synthesizer: synthesizer,
		translator:  translator,
	}
}

func (p *englishPhraseProblemProcessor) AddProblem(ctx context.Context, repo app.RepositoryFactory, operator app.Student, workbook app.Workbook, param app.ProblemAddParameter) (app.Added, app.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem1")

	problemRepo, err := repo.NewProblemRepository(ctx, workbook.GetProblemType())
	if err != nil {
		return 0, 0, fmt.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	extractedParam, err := toEnglishPhraseProblemAddParemeter(param)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to toNewEnglishPhraseProblemParemeter. err: %w", err)
	}

	audioID, err := p.findOrAddAudio(ctx, repo, extractedParam.Text)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to p.findOrAddAudio. err: %w", err)
	}

	if audioID == 0 {
		return 0, 0, fmt.Errorf("audio ID is zero. text: %s", extractedParam.Text)
	}

	problemID, err := p.addSingleProblem(ctx, operator, problemRepo, param, extractedParam, audioID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to addSingleProblem: extractedParam: %+v, err: %w", extractedParam, err)
	}

	return 1, problemID, err
}

func (p *englishPhraseProblemProcessor) addSingleProblem(ctx context.Context, operator app.Student, problemRepo app.ProblemRepository, param app.ProblemAddParameter, extractedParam *englishPhraseProblemAddParemeter, audioID app.AudioID) (app.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem1")

	logger.Infof("text: %s, audio ID: %d", extractedParam.Text, audioID)

	properties := map[string]string{
		"text":       extractedParam.Text,
		"translated": extractedParam.Translated,
		"lang":       extractedParam.Lang,
		"audioId":    strconv.Itoa(int(audioID)),
	}
	newParam, err := app.NewProblemAddParameter(param.GetWorkbookID(), param.GetNumber(), properties)
	if err != nil {
		return 0, fmt.Errorf("failed to NewParameter. err: %w", err)
	}

	problemID, err := problemRepo.AddProblem(ctx, operator, newParam)
	if err != nil {
		return 0, fmt.Errorf("failed to problemRepo.AddProblem. err: %w", err)
	}

	return problemID, nil

}

func (p *englishPhraseProblemProcessor) RemoveProblem(ctx context.Context, repo app.RepositoryFactory, operator app.Student, id app.ProblemSelectParameter2) error {
	problemRepo, err := repo.NewProblemRepository(ctx, EnglishPhraseProblemType)
	if err != nil {
		return fmt.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	if err := problemRepo.RemoveProblem(ctx, operator, id); err != nil {
		return err
	}

	return nil
}

func (p *englishPhraseProblemProcessor) findOrAddAudio(ctx context.Context, repo app.RepositoryFactory, text string) (app.AudioID, error) {
	audioRepo, err := repo.NewAudioRepository(ctx)
	if err != nil {
		return 0, err
	}

	{
		id, err := audioRepo.FindAudioIDByText(ctx, app.Lang5ENUS, text)
		if err != nil {
			if !errors.Is(err, app.ErrAudioNotFound) {
				return 0, fmt.Errorf("failed to FindAudioID. err: %w", err)
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
