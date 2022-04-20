package service

import (
	"context"
	"errors"
	"io"
	"strconv"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	appS "github.com/kujilabo/cocotola-api/pkg_app/service"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	pluginS "github.com/kujilabo/cocotola-api/pkg_plugin/common/service"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
	"golang.org/x/xerrors"
)

var (
	EnglishSentenceProblemAddPropertyAudioID                = "audioId"
	EnglishSentenceProblemAddPropertyLang                   = "lang"
	EnglishSentenceProblemAddPropertyText                   = "text"
	EnglishSentenceProblemAddPropertyTranslated             = "translated"
	EnglishSentenceProblemAddPropertyProvider               = "provider"
	EnglishSentenceProblemAddPropertyTatoebaSentenceNumber1 = "tatoebaSentenceNumber1"
	EnglishSentenceProblemAddPropertyTatoebaSentenceNumber2 = "tatoebaSentenceNumber2"
	EnglishSentenceProblemAddPropertyTatoebaAuthor1         = "tatoebaAuthor1"
	EnglishSentenceProblemAddPropertyTatoebaAuthor2         = "tatoebaAuthor2"
)

type englishSentenceProblemAddParemeter struct {
	Lang                   app.Lang2 `validate:"required"`
	Text                   string    `validate:"required"`
	Translated             string
	Provider               string
	TatoebaSentenceNumber1 int
	TatoebaSentenceNumber2 int
	TatoebaAuthor1         string
	TatoebaAuthor2         string
}

func (p *englishSentenceProblemAddParemeter) toProperties(audioID app.AudioID) map[string]string {
	return map[string]string{
		EnglishSentenceProblemAddPropertyAudioID:                strconv.Itoa(int(uint(audioID))),
		EnglishSentenceProblemAddPropertyLang:                   p.Lang.String(),
		EnglishSentenceProblemAddPropertyText:                   p.Text,
		EnglishSentenceProblemAddPropertyTranslated:             p.Translated,
		EnglishSentenceProblemAddPropertyProvider:               p.Provider,
		EnglishSentenceProblemAddPropertyTatoebaSentenceNumber1: strconv.Itoa(p.TatoebaSentenceNumber1),
		EnglishSentenceProblemAddPropertyTatoebaSentenceNumber2: strconv.Itoa(p.TatoebaSentenceNumber2),
		EnglishSentenceProblemAddPropertyTatoebaAuthor1:         p.TatoebaAuthor1,
		EnglishSentenceProblemAddPropertyTatoebaAuthor2:         p.TatoebaAuthor2,
	}
}

func toEnglishSentenceProblemAddParemeter(param appS.ProblemAddParameter) (*englishSentenceProblemAddParemeter, error) {
	if _, ok := param.GetProperties()["text"]; !ok {
		return nil, xerrors.Errorf("text is not defined. err: %w", lib.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["translated"]; !ok {
		return nil, xerrors.Errorf("translated is not defined. err: %w", lib.ErrInvalidArgument)
	}
	if _, ok := param.GetProperties()["lang"]; !ok {
		return nil, xerrors.Errorf("lang is not defined. err: %w", lib.ErrInvalidArgument)
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

	return m, lib.Validator.Struct(m)
}

type EnglishSentenceProblemProcessor interface {
	appS.ProblemAddProcessor
	appS.ProblemRemoveProcessor
	appS.ProblemImportProcessor
}

type englishSentenceProblemProcessor struct {
	synthesizer                     pluginS.Synthesizer
	translationClient               pluginS.TranslationClient
	newProblemAddParameterCSVReader func(workbookID app.WorkbookID, reader io.Reader) appS.ProblemAddParameterIterator
}

func NewEnglishSentenceProblemProcessor(synthesizer pluginS.Synthesizer, translationClient pluginS.TranslationClient, newProblemAddParameterCSVReader func(workbookID app.WorkbookID, reader io.Reader) appS.ProblemAddParameterIterator) EnglishSentenceProblemProcessor {
	return &englishSentenceProblemProcessor{
		synthesizer:                     synthesizer,
		translationClient:               translationClient,
		newProblemAddParameterCSVReader: newProblemAddParameterCSVReader,
	}
}

func (p *englishSentenceProblemProcessor) AddProblem(ctx context.Context, repo appS.RepositoryFactory, operator app.StudentModel, workbook app.WorkbookModel, param appS.ProblemAddParameter) ([]app.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem1")

	problemRepo, err := repo.NewProblemRepository(ctx, domain.EnglishSentenceProblemType)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	extractedParam, err := toEnglishSentenceProblemAddParemeter(param)
	if err != nil {
		return nil, xerrors.Errorf("failed to toNewEnglishSentenceProblemParemeter. err: %w", err)
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

	return []app.ProblemID{problemID}, nil
}

func (p *englishSentenceProblemProcessor) addSingleProblem(ctx context.Context, operator app.StudentModel, problemRepo appS.ProblemRepository, param appS.ProblemAddParameter, extractedParam *englishSentenceProblemAddParemeter, audioID app.AudioID) (app.ProblemID, error) {
	logger := log.FromContext(ctx)

	logger.Infof("text: %s, audio ID: %d", extractedParam.Text, audioID)

	properties := extractedParam.toProperties(audioID)
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
func (p *englishSentenceProblemProcessor) RemoveProblem(ctx context.Context, repo appS.RepositoryFactory, operator app.StudentModel, id appS.ProblemSelectParameter2) error {
	problemRepo, err := repo.NewProblemRepository(ctx, domain.EnglishSentenceProblemType)
	if err != nil {
		return xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	if err := problemRepo.RemoveProblem(ctx, operator, id); err != nil {
		return xerrors.Errorf("failed to RemoveProblem. err: %w", err)
	}

	return nil
}

func (p *englishSentenceProblemProcessor) CreateCSVReader(ctx context.Context, workbookID app.WorkbookID, reader io.Reader) (appS.ProblemAddParameterIterator, error) {
	return p.newProblemAddParameterCSVReader(workbookID, reader), nil
}

func (p *englishSentenceProblemProcessor) findOrAddAudio(ctx context.Context, repo appS.RepositoryFactory, text string) (app.AudioID, error) {
	audioRepo := repo.NewAudioRepository(ctx)

	{
		id, err := audioRepo.FindAudioIDByText(ctx, app.Lang5ENUS, text)
		if err != nil {
			if !errors.Is(err, appS.ErrAudioNotFound) {
				return 0, xerrors.Errorf("failed to FindAudioIDByText. err: %w", err)
			}
		} else {
			return id, nil
		}
	}

	audioContent, err := p.synthesizer.Synthesize(ctx, app.Lang5ENUS, text)
	if err != nil {
		return 0, err
	}

	id, err := audioRepo.AddAudio(ctx, app.Lang5ENUS, text, audioContent)
	if err != nil {
		return 0, err
	}

	return id, err
}
