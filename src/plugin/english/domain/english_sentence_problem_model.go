//go:generate mockery --output mock --name EnglishSentenceProblemModel
package domain

import (
	"context"

	appD "github.com/kujilabo/cocotola-api/src/app/domain"
)

const EnglishSentenceProblemType = "english_sentence"

type EnglishSentenceProblemModel interface {
	appD.ProblemModel
	GetProvider() string
	GetAudioID() appD.AudioID
	GetText() string
	GetNote() string
}

type englishSentenceProblemModel struct {
	appD.ProblemModel
	Provider   string
	AudioID    appD.AudioID
	Text       string
	Lang2      appD.Lang2
	Translated string
	Note       string
}

func NewEnglishSentenceProblemModel(problemModel appD.ProblemModel, audioID appD.AudioID, provider string, text string, lang2 appD.Lang2, translated string, note string) (EnglishSentenceProblemModel, error) {
	return &englishSentenceProblemModel{
		ProblemModel: problemModel,
		AudioID:      audioID,
		Text:         text,
		Lang2:        lang2,
		Translated:   translated,
		Note:         note,
	}, nil
}

func (m *englishSentenceProblemModel) GetProvider() string {
	return m.Provider
}

func (m *englishSentenceProblemModel) GetAudioID() appD.AudioID {
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
		"lang2":      m.Lang2,
		"translated": m.Translated,
		"note":       m.Note,
	}
}
