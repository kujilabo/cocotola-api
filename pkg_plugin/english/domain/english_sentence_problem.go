package domain

import (
	"context"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

const EnglishSentenceProblemType = "english_sentence"

type EnglishSentenceProblem interface {
	app.Problem
	GetProvider() string
	GetAudioID() app.AudioID
	GetText() string
	GetNote() string
}

type englishSentenceProblem struct {
	app.Problem
	Provider   string
	AudioID    app.AudioID
	Text       string
	Lang       app.Lang2
	Translated string
	Note       string
}

func NewEnglishSentenceProblem(problem app.Problem, audioID app.AudioID, provider string, text string, lang app.Lang2, translated string, note string) (EnglishSentenceProblem, error) {
	return &englishSentenceProblem{
		Problem:    problem,
		AudioID:    audioID,
		Text:       text,
		Lang:       lang,
		Translated: translated,
		Note:       note,
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

func (m *englishSentenceProblem) GetNote() string {
	return m.Note
}

func (m *englishSentenceProblem) Properties(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"audioId":    uint(m.AudioID),
		"text":       m.Text,
		"lang":       m.Lang,
		"translated": m.Translated,
		"note":       m.Note,
	}
}
