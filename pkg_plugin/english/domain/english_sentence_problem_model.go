package domain

import (
	"context"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

const EnglishSentenceProblemType = "english_sentence"

type EnglishSentenceProblemModel interface {
	app.ProblemModel
	GetProvider() string
	GetAudioID() app.AudioID
	GetText() string
	GetNote() string
}

type englishSentenceProblemModel struct {
	app.ProblemModel
	Provider   string
	AudioID    app.AudioID
	Text       string
	Lang       app.Lang2
	Translated string
	Note       string
}

func NewEnglishSentenceProblemModel(problemModel app.ProblemModel, audioID app.AudioID, provider string, text string, lang app.Lang2, translated string, note string) (EnglishSentenceProblemModel, error) {
	return &englishSentenceProblemModel{
		ProblemModel: problemModel,
		AudioID:      audioID,
		Text:         text,
		Lang:         lang,
		Translated:   translated,
		Note:         note,
	}, nil
}

func (m *englishSentenceProblemModel) GetProvider() string {
	return m.Provider
}

func (m *englishSentenceProblemModel) GetAudioID() app.AudioID {
	return m.AudioID
}

func (m *englishSentenceProblemModel) GetText() string {
	return m.Text
}

func (m *englishSentenceProblemModel) GetNote() string {
	return m.Note
}

func (m *englishSentenceProblemModel) Properties(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"audioId":    uint(m.AudioID),
		"text":       m.Text,
		"lang":       m.Lang,
		"translated": m.Translated,
		"note":       m.Note,
	}
}
