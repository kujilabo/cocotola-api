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
	GetLang() appD.Lang2
	GetTranslated() string
}

type englishPhraseProblemModel struct {
	appD.ProblemModel
	AudioID    appD.AudioID
	Text       string
	Lang       appD.Lang2
	Translated string
}

func NewEnglishPhraseProblemModel(problemModel appD.ProblemModel, audioID appD.AudioID, text string, lang appD.Lang2, translated string) (EnglishPhraseProblemModel, error) {
	m := &englishPhraseProblemModel{
		ProblemModel: problemModel,
		AudioID:      audioID,
		Text:         text,
		Lang:         lang,
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

func (m *englishPhraseProblemModel) GetLang() appD.Lang2 {
	return m.Lang
}
func (m *englishPhraseProblemModel) GetTranslated() string {
	return m.Translated
}

func (m *englishPhraseProblemModel) Properties(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"audioId":    uint(m.AudioID),
		"text":       m.Text,
		"lang":       m.Lang,
		"translated": m.Translated,
	}
}
