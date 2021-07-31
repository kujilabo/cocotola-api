package domain

import (
	"context"
	"strconv"

	"github.com/go-playground/validator"
	"golang.org/x/xerrors"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	plugin "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
)

type englishWordProblemAddParemeter struct {
	Text       string         `validate:"required"`
	Pos        plugin.WordPos `validate:"required"`
	Translated string
}

func toEnglishWordProblemAddParemeter(param app.ProblemAddParameter) (*englishWordProblemAddParemeter, error) {
	posS := param.GetProperties()["pos"]
	pos, err := strconv.Atoi(posS)
	if err != nil {
		return nil, xerrors.Errorf("faield to cast to int. err: %w", err)
	}

	m := &englishWordProblemAddParemeter{
		Text:       param.GetProperties()["text"],
		Pos:        plugin.WordPos(pos),
		Translated: param.GetProperties()["translated"],
	}

	v := validator.New()
	return m, v.Struct(m)
}

type EnglishWordProblemProcessor interface {
	app.ProblemAddProcessor
	app.ProblemRemoveProcessor
}

type englishWordProblemProcessor struct {
	synthesizer plugin.Synthesizer
	translator  plugin.Translator
}

func NewEnglishWordProblemProcessor(synthesizer plugin.Synthesizer, translator plugin.Translator) EnglishWordProblemProcessor {
	return &englishWordProblemProcessor{
		synthesizer: synthesizer,
		translator:  translator,
	}
}

func (p *englishWordProblemProcessor) AddProblem(ctx context.Context, repo app.RepositoryFactory, operator app.Student, param app.ProblemAddParameter) (app.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem1")

	problemRepo, err := repo.NewProblemRepository(ctx, param.GetProblemType())
	if err != nil {
		return 0, xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	extractedParam, err := toEnglishWordProblemAddParemeter(param)
	if err != nil {
		return 0, xerrors.Errorf("failed to toNewEnglishWordProblemParemeter. err: %w", err)
	}

	audioID, err := p.findOrAddAudio(ctx, repo, extractedParam.Text)
	if err != nil {
		return 0, xerrors.Errorf("failed to p.findOrAddAudio. err: %w", err)
	}

	if audioID == 0 {
		return 0, xerrors.Errorf("audio ID is zero. text: %s", extractedParam.Text)
	}

	logger.Infof("text: %s, audio ID: %d", extractedParam.Text, audioID)

	if extractedParam.Translated == "" && extractedParam.Pos == plugin.PosOther {
		translated, err := p.translate(ctx, extractedParam.Text, app.Lang2EN, app.Lang2JA)
		if err != nil {
			logger.Errorf("translate err: %v", err)
			properties := map[string]string{
				"text":       extractedParam.Text,
				"translated": extractedParam.Translated,
				"pos":        strconv.Itoa(int(extractedParam.Pos)),
				"audioId":    strconv.Itoa(int(audioID)),
			}
			newParam, err := app.NewProblemAddParameter(param.GetWorkbookID(), param.GetNumber(), param.GetProblemType(), properties)
			if err != nil {
				return 0, xerrors.Errorf("failed to NewParameter. err: %w", err)
			}

			problemID, err := problemRepo.AddProblem(ctx, operator, newParam)
			if err != nil {
				return 0, xerrors.Errorf("failed to problemRepo.AddProblem. err: %w", err)
			}

			return problemID, nil
		}

		for _, t := range translated {
			extractedParam.Pos = t.Pos
			extractedParam.Translated = t.Target
			properties := map[string]string{
				"text":       extractedParam.Text,
				"translated": extractedParam.Translated,
				"pos":        strconv.Itoa(int(extractedParam.Pos)),
				"audioId":    strconv.Itoa(int(audioID)),
			}
			newParam, err := app.NewProblemAddParameter(param.GetWorkbookID(), param.GetNumber(), param.GetProblemType(), properties)
			if err != nil {
				return 0, xerrors.Errorf("failed to NewParameter. err: %w", err)
			}

			if _, err := problemRepo.AddProblem(ctx, operator, newParam); err != nil {
				return 0, xerrors.Errorf("failed to problemRepo.AddProblem. err: %w", err)
			}
		}

		return 0, nil
	}

	if extractedParam.Translated == "" {
		translated, err := p.translateWithPos(ctx, extractedParam.Text, extractedParam.Pos, app.Lang2EN, app.Lang2JA)
		if err != nil {
			logger.Errorf("translate err: %v", err)
		} else {
			extractedParam.Translated = translated
		}
	}

	properties := map[string]string{
		"text":       extractedParam.Text,
		"translated": extractedParam.Translated,
		"pos":        strconv.Itoa(int(extractedParam.Pos)),
		"audioId":    strconv.Itoa(int(audioID)),
	}
	newParam, err := app.NewProblemAddParameter(param.GetWorkbookID(), param.GetNumber(), param.GetProblemType(), properties)
	if err != nil {
		return 0, xerrors.Errorf("failed to NewParameter. err: %w", err)
	}

	problemID, err := problemRepo.AddProblem(ctx, operator, newParam)
	if err != nil {
		return 0, xerrors.Errorf("failed to problemRepo.AddProblem. err: %w", err)
	}

	return problemID, nil
}

func (p *englishWordProblemProcessor) RemoveProblem(ctx context.Context, repo app.RepositoryFactory, operator app.Student, problemID app.ProblemID, version int) error {
	problemRepo, err := repo.NewProblemRepository(ctx, EnglishWordProblemType)
	if err != nil {
		return xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	if err := problemRepo.RemoveProblem(ctx, operator, problemID, version); err != nil {
		return err
	}

	return nil
}

func (p *englishWordProblemProcessor) findOrAddAudio(ctx context.Context, repo app.RepositoryFactory, text string) (app.AudioID, error) {
	audioRepo, err := repo.NewAudioRepository(ctx)
	if err != nil {
		return 0, err
	}

	{
		id, err := audioRepo.FindAudioID(ctx, app.Lang5ENUS, text)
		if err != nil {
			if !xerrors.Is(err, app.ErrAudioNotFound) {
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

func (p *englishWordProblemProcessor) translateWithPos(ctx context.Context, text string, pos plugin.WordPos, fromLang, toLang app.Lang2) (string, error) {
	logger := log.FromContext(ctx)
	logger.Infof("translateWithPos. text: %s", text)

	result, err := p.translator.DictionaryLookup(ctx, text, fromLang, toLang)
	if err != nil {
		return "", err
	}

	logger.Infof("translate:%v", result)
	var translated string
	var confidence = 0.0
	for _, r := range result {
		if r.Pos == pos && r.Confidence > confidence {
			confidence = r.Confidence
			translated = r.Target
		}
	}

	return translated, nil
}

func (p *englishWordProblemProcessor) translate(ctx context.Context, text string, fromLang, toLang app.Lang2) ([]plugin.TranslationResult, error) {
	logger := log.FromContext(ctx)
	logger.Infof("translate. text: %s", text)

	result, err := p.translator.DictionaryLookup(ctx, text, fromLang, toLang)
	if err != nil {
		return nil, err
	}

	logger.Infof("translate:%v", result)

	posList := make(map[plugin.WordPos]plugin.TranslationResult)
	for _, r := range result {
		if _, ok := posList[r.Pos]; !ok {
			posList[r.Pos] = r
		} else if r.Confidence > posList[r.Pos].Confidence {
			posList[r.Pos] = r
		}
	}

	translated := make([]plugin.TranslationResult, 0)
	for _, v := range posList {
		translated = append(translated, v)
	}

	return translated, nil
}
