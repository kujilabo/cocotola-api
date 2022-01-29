package domain

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/go-playground/validator"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	plugin "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
)

type englishSentenceProblemAddParemeter struct {
	Lang       app.Lang2 `validate:"required"`
	Text       string    `validate:"required"`
	Translated string
}

func toEnglishSentenceProblemAddParemeter(param app.ProblemAddParameter) (*englishSentenceProblemAddParemeter, error) {
	if _, ok := param.GetProperties()["text"]; !ok {
		return nil, fmt.Errorf("text is not defined. err: %w", lib.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["translated"]; !ok {
		return nil, fmt.Errorf("translated is not defined. err: %w", lib.ErrInvalidArgument)
	}
	if _, ok := param.GetProperties()["lang"]; !ok {
		return nil, fmt.Errorf("lang is not defined. err: %w", lib.ErrInvalidArgument)
	}

	lang2, err := app.NewLang2(param.GetProperties()["lang"])
	if err != nil {
		return nil, err
	}

	m := &englishSentenceProblemAddParemeter{
		Lang:       lang2,
		Text:       param.GetProperties()["text"],
		Translated: param.GetProperties()["translated"],
	}

	v := validator.New()
	return m, v.Struct(m)
}

type EnglishSentenceProblemProcessor interface {
	app.ProblemAddProcessor
	app.ProblemRemoveProcessor
	app.ProblemImportProcessor
}

type englishSentenceProblemProcessor struct {
	synthesizer                     plugin.Synthesizer
	translator                      plugin.Translator
	newProblemAddParameterCSVReader func(workbookID app.WorkbookID, reader io.Reader) app.ProblemAddParameterIterator
}

func NewEnglishSentenceProblemProcessor(synthesizer plugin.Synthesizer, translator plugin.Translator, newProblemAddParameterCSVReader func(workbookID app.WorkbookID, reader io.Reader) app.ProblemAddParameterIterator) EnglishSentenceProblemProcessor {
	return &englishSentenceProblemProcessor{
		synthesizer:                     synthesizer,
		translator:                      translator,
		newProblemAddParameterCSVReader: newProblemAddParameterCSVReader,
	}
}

func (p *englishSentenceProblemProcessor) AddProblem(ctx context.Context, repo app.RepositoryFactory, operator app.Student, workbook app.Workbook, param app.ProblemAddParameter) (app.Added, app.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem1")

	problemRepo, err := repo.NewProblemRepository(ctx, workbook.GetProblemType())
	if err != nil {
		return 0, 0, fmt.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	extractedParam, err := toEnglishSentenceProblemAddParemeter(param)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to toNewEnglishSentenceProblemParemeter. err: %w", err)
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

	return 1, problemID, nil
}

func (p *englishSentenceProblemProcessor) addSingleProblem(ctx context.Context, operator app.Student, problemRepo app.ProblemRepository, param app.ProblemAddParameter, extractedParam *englishSentenceProblemAddParemeter, audioID app.AudioID) (app.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem1")

	logger.Infof("text: %s, audio ID: %d", extractedParam.Text, audioID)

	properties := map[string]string{
		"lang":       extractedParam.Lang.String(),
		"text":       extractedParam.Text,
		"translated": extractedParam.Translated,
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
func (p *englishSentenceProblemProcessor) RemoveProblem(ctx context.Context, repo app.RepositoryFactory, operator app.Student, problemID app.ProblemID, version int) error {
	problemRepo, err := repo.NewProblemRepository(ctx, EnglishSentenceProblemType)
	if err != nil {
		return fmt.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	if err := problemRepo.RemoveProblem(ctx, operator, problemID, version); err != nil {
		return fmt.Errorf("failed to RemoveProblem. err: %w", err)
	}

	return nil
}

func (p *englishSentenceProblemProcessor) CreateCSVReader(ctx context.Context, workbookID app.WorkbookID, reader io.Reader) (app.ProblemAddParameterIterator, error) {
	return p.newProblemAddParameterCSVReader(workbookID, reader), nil
}

func (p *englishSentenceProblemProcessor) findOrAddAudio(ctx context.Context, repo app.RepositoryFactory, text string) (app.AudioID, error) {
	audioRepo, err := repo.NewAudioRepository(ctx)
	if err != nil {
		return 0, err
	}

	{
		id, err := audioRepo.FindAudioIDByText(ctx, app.Lang5ENUS, text)
		if err != nil {
			if !errors.Is(err, app.ErrAudioNotFound) {
				return 0, fmt.Errorf("failed to FindAudioIDByText. err: %w", err)
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
