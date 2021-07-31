package domain

import (
	"context"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

type EnglishPhraseProblem interface {
	app.Problem
	GetAudioID() app.AudioID
	GetText() string
}

type englishPhraseProblem struct {
	app.Problem
	AudioID    app.AudioID
	Text       string
	Lang       string
	Translated string
}

func NewEnglishPhraseProblem(problem app.Problem, audioID app.AudioID, text, lang, translated string) (EnglishPhraseProblem, error) {
	return &englishPhraseProblem{
		Problem:    problem,
		AudioID:    audioID,
		Text:       text,
		Lang:       lang,
		Translated: translated,
	}, nil
}

func (m *englishPhraseProblem) GetAudioID() app.AudioID {
	return m.AudioID
}

func (m *englishPhraseProblem) GetText() string {
	return m.Text
}

func (m *englishPhraseProblem) Properties(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"audioId":    uint(m.AudioID),
		"text":       m.Text,
		"lang":       m.Lang,
		"translated": m.Translated,
	}
}
