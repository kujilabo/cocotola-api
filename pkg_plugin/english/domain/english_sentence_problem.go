package domain

import (
	"context"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

type EnglishSentenceProblem interface {
	app.Problem
	GetProvider() string
	GetAudioID() app.AudioID
	GetText() string
}

type englishSentenceProblem struct {
	app.Problem
	Provider   string
	AudioID    app.AudioID
	Text       string
	Lang       string
	Translated string
}

func NewEnglishSentenceProblem(problem app.Problem, audioID app.AudioID, provider string, text, lang, translated string) (EnglishSentenceProblem, error) {
	return &englishSentenceProblem{
		Problem:    problem,
		AudioID:    audioID,
		Text:       text,
		Lang:       lang,
		Translated: translated,
	}, nil
}

func (m *englishSentenceProblem) GetProvider() string {
	return m.Provider
}

func (m *englishSentenceProblem) GetAudioID() app.AudioID {
	return m.AudioID
}

func (m *englishSentenceProblem) GetText() string {
	return m.Text
}

func (m *englishSentenceProblem) Properties(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"audioId":    uint(m.AudioID),
		"text":       m.Text,
		"lang":       m.Lang,
		"translated": m.Translated,
	}
}
