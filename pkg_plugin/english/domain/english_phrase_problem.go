package domain

import (
	"context"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

const EnglishPhraseProblemType = "english_phrase"

type EnglishPhraseProblem interface {
	app.Problem
	GetAudioID() app.AudioID
	GetText() string
	GetLang() app.Lang2
	GetTranslated() string
}

type englishPhraseProblem struct {
	app.Problem
	AudioID    app.AudioID
	Text       string
	Lang       app.Lang2
	Translated string
}

func NewEnglishPhraseProblem(problem app.Problem, audioID app.AudioID, text string, lang app.Lang2, translated string) (EnglishPhraseProblem, error) {
	m := &englishPhraseProblem{
		Problem:    problem,
		AudioID:    audioID,
		Text:       text,
		Lang:       lang,
		Translated: translated,
	}

	return m, lib.Validator.Struct(m)
}

func (m *englishPhraseProblem) GetAudioID() app.AudioID {
	return m.AudioID
}

func (m *englishPhraseProblem) GetText() string {
	return m.Text
}

func (m *englishPhraseProblem) GetLang() app.Lang2 {
	return m.Lang
}
func (m *englishPhraseProblem) GetTranslated() string {
	return m.Translated
}

func (m *englishPhraseProblem) Properties(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"audioId":    uint(m.AudioID),
		"text":       m.Text,
		"lang":       m.Lang,
		"translated": m.Translated,
	}
}
