package service

import (
	"context"
	"io"
	"strconv"

	appD "github.com/kujilabo/cocotola-api/src/app/domain"
	appS "github.com/kujilabo/cocotola-api/src/app/service"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	liberrors "github.com/kujilabo/cocotola-api/src/lib/errors"
	"github.com/kujilabo/cocotola-api/src/lib/log"
	pluginS "github.com/kujilabo/cocotola-api/src/plugin/common/service"
	"github.com/kujilabo/cocotola-api/src/plugin/english/domain"
)

var (
	EnglishSentenceProblemAddPropertyAudioID                = "audioId"
	EnglishSentenceProblemAddPropertyLang2                  = "lang2"
	EnglishSentenceProblemAddPropertyText                   = "text"
	EnglishSentenceProblemAddPropertyTranslated             = "translated"
	EnglishSentenceProblemAddPropertyProvider               = "provider"
	EnglishSentenceProblemAddPropertyTatoebaSentenceNumber1 = "tatoebaSentenceNumber1"
	EnglishSentenceProblemAddPropertyTatoebaSentenceNumber2 = "tatoebaSentenceNumber2"
	EnglishSentenceProblemAddPropertyTatoebaAuthor1         = "tatoebaAuthor1"
	EnglishSentenceProblemAddPropertyTatoebaAuthor2         = "tatoebaAuthor2"
)

type englishSentenceProblemAddParemeter struct {
	Lang2                  appD.Lang2 `validate:"required"`
	Text                   string     `validate:"required"`
	Translated             string
	Provider               string
	TatoebaSentenceNumber1 int
	TatoebaSentenceNumber2 int
	TatoebaAuthor1         string
	TatoebaAuthor2         string
}

func (p *englishSentenceProblemAddParemeter) toProperties(audioID appD.AudioID) map[string]string {
	return map[string]string{
		EnglishSentenceProblemAddPropertyAudioID:                strconv.Itoa(int(uint(audioID))),
		EnglishSentenceProblemAddPropertyLang2:                  p.Lang2.String(),
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
		return nil, liberrors.Errorf("text is not defined. err: %w", libD.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["translated"]; !ok {
		return nil, liberrors.Errorf("translated is not defined. err: %w", libD.ErrInvalidArgument)
	}
	if _, ok := param.GetProperties()["lang2"]; !ok {
		return nil, liberrors.Errorf("lang2 is not defined. err: %w", libD.ErrInvalidArgument)
	}

	lang2, err := appD.NewLang2(param.GetProperties()["lang2"])
	if err != nil {
		return nil, err
	}

	m := &englishSentenceProblemAddParemeter{
		Lang2:      lang2,
		Text:       param.GetProperties()["text"],
		Translated: param.GetProperties()["translated"],
	}

	return m, libD.Validator.Struct(m)
}

type EnglishSentenceProblemProcessor interface {
	appS.ProblemAddProcessor
	appS.ProblemRemoveProcessor
	appS.ProblemImportProcessor
}

type englishSentenceProblemProcessor struct {
	synthesizerClient               appS.SynthesizerClient
	translatorClient                pluginS.TranslatorClient
	newProblemAddParameterCSVReader func(workbookID appD.WorkbookID, reader io.Reader) appS.ProblemAddParameterIterator
}

func NewEnglishSentenceProblemProcessor(synthesizerClient appS.SynthesizerClient, translatorClient pluginS.TranslatorClient, newProblemAddParameterCSVReader func(workbookID appD.WorkbookID, reader io.Reader) appS.ProblemAddParameterIterator) EnglishSentenceProblemProcessor {
	return &englishSentenceProblemProcessor{
		synthesizerClient:               synthesizerClient,
		translatorClient:                translatorClient,
		newProblemAddParameterCSVReader: newProblemAddParameterCSVReader,
	}
}

func (p *englishSentenceProblemProcessor) AddProblem(ctx context.Context, repo appS.RepositoryFactory, operator appD.StudentModel, workbook appD.WorkbookModel, param appS.ProblemAddParameter) ([]appD.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem1")

	problemRepo, err := repo.NewProblemRepository(ctx, domain.EnglishSentenceProblemType)
	if err != nil {
		return nil, liberrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	extractedParam, err := toEnglishSentenceProblemAddParemeter(param)
	if err != nil {
		return nil, liberrors.Errorf("failed to toNewEnglishSentenceProblemParemeter. err: %w", err)
	}

	audio, err := p.synthesizerClient.Synthesize(ctx, appD.Lang2EN, extractedParam.Text)
	if err != nil {
		return nil, err
	}

	problemID, err := p.addSingleProblem(ctx, operator, problemRepo, param, extractedParam, appD.AudioID(audio.GetAudioModel().GetID()))
	if err != nil {
		return nil, liberrors.Errorf("failed to addSingleProblem: extractedParam: %+v, err: %w", extractedParam, err)
	}

	return []appD.ProblemID{problemID}, nil
}

func (p *englishSentenceProblemProcessor) addSingleProblem(ctx context.Context, operator appD.StudentModel, problemRepo appS.ProblemRepository, param appS.ProblemAddParameter, extractedParam *englishSentenceProblemAddParemeter, audioID appD.AudioID) (appD.ProblemID, error) {
	logger := log.FromContext(ctx)

	logger.Infof("text: %s, audio ID: %d", extractedParam.Text, audioID)

	properties := extractedParam.toProperties(audioID)
	newParam, err := appS.NewProblemAddParameter(param.GetWorkbookID(), param.GetNumber(), properties)
	if err != nil {
		return 0, liberrors.Errorf("failed to NewParameter. err: %w", err)
	}

	problemID, err := problemRepo.AddProblem(ctx, operator, newParam)
	if err != nil {
		return 0, liberrors.Errorf("failed to problemRepo.AddProblem. err: %w", err)
	}

	return problemID, nil

}
func (p *englishSentenceProblemProcessor) RemoveProblem(ctx context.Context, repo appS.RepositoryFactory, operator appD.StudentModel, id appS.ProblemSelectParameter2) error {
	problemRepo, err := repo.NewProblemRepository(ctx, domain.EnglishSentenceProblemType)
	if err != nil {
		return liberrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	if err := problemRepo.RemoveProblem(ctx, operator, id); err != nil {
		return liberrors.Errorf("failed to RemoveProblem. err: %w", err)
	}

	return nil
}

func (p *englishSentenceProblemProcessor) CreateCSVReader(ctx context.Context, workbookID appD.WorkbookID, reader io.Reader) (appS.ProblemAddParameterIterator, error) {
	return p.newProblemAddParameterCSVReader(workbookID, reader), nil
}
