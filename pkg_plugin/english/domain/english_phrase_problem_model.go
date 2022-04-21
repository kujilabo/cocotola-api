//go:generate mockery --output mock --name EnglishPhraseProblemModel
package domain

import (
	"context"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

const EnglishPhraseProblemType = "english_phrase"

type EnglishPhraseProblemModel interface {
	app.ProblemModel
	GetAudioID() app.AudioID
	GetText() string
	GetLang() app.Lang2
	GetTranslated() string
}

type englishPhraseProblemModel struct {
	app.ProblemModel
	AudioID    app.AudioID
	Text       string
	Lang       app.Lang2
	Translated string
}

func NewEnglishPhraseProblemModel(problemModel app.ProblemModel, audioID app.AudioID, text string, lang app.Lang2, translated string) (EnglishPhraseProblemModel, error) {
	m := &englishPhraseProblemModel{
		ProblemModel: problemModel,
		AudioID:      audioID,
		Text:         text,
		Lang:         lang,
		Translated:   translated,
	}

	return m, lib.Validator.Struct(m)
}

func (m *englishPhraseProblemModel) GetAudioID() app.AudioID {
	return m.AudioID
}

func (m *englishPhraseProblemModel) GetText() string {
	return m.Text
}

func (m *englishPhraseProblemModel) GetLang() app.Lang2 {
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
