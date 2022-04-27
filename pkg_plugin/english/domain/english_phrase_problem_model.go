//go:generate mockery --output mock --name EnglishPhraseProblemModel
package domain

import (
	"context"

	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

const EnglishPhraseProblemType = "english_phrase"

type EnglishPhraseProblemModel interface {
	appD.ProblemModel
	GetAudioID() appD.AudioID
	GetText() string
	GetLang2() appD.Lang2
	GetTranslated() string
}

type englishPhraseProblemModel struct {
	appD.ProblemModel
	AudioID    appD.AudioID
	Text       string
	Lang2      appD.Lang2
	Translated string
}

func NewEnglishPhraseProblemModel(problemModel appD.ProblemModel, audioID appD.AudioID, text string, lang2 appD.Lang2, translated string) (EnglishPhraseProblemModel, error) {
	m := &englishPhraseProblemModel{
		ProblemModel: problemModel,
		AudioID:      audioID,
		Text:         text,
		Lang2:        lang2,
		Translated:   translated,
	}

	return m, libD.Validator.Struct(m)
}

func (m *englishPhraseProblemModel) GetAudioID() appD.AudioID {
	return m.AudioID
}

func (m *englishPhraseProblemModel) GetText() string {
	return m.Text
}

func (m *englishPhraseProblemModel) GetLang2() appD.Lang2 {
	return m.Lang2
}

func (m *englishPhraseProblemModel) GetTranslated() string {
	return m.Translated
}

func (m *englishPhraseProblemModel) Properties(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"audioId":    uint(m.AudioID),
		"text":       m.Text,
		"lang2":      m.Lang2,
		"translated": m.Translated,
	}
}
