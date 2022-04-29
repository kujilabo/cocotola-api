//go:generate mockery --output mock --name Translation
package domain

import (
	appD "github.com/kujilabo/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
)

type Translation interface {
	GetText() string
	GetPos() WordPos
	GetLang2() appD.Lang2
	GetTranslated() string
	GetProvider() string
}

type translation struct {
	Text       string `validate:"required"`
	Pos        WordPos
	Lang2      appD.Lang2
	Translated string
	Provider   string
}

func NewTranslation(text string, pos WordPos, lang2 appD.Lang2, translated, provider string) (Translation, error) {
	m := &translation{
		Text:       text,
		Pos:        pos,
		Lang2:      lang2,
		Translated: translated,
		Provider:   provider,
	}

	return m, libD.Validator.Struct(m)
}

func (t *translation) GetText() string {
	return t.Text
}

func (t *translation) GetPos() WordPos {
	return t.Pos
}

func (t *translation) GetLang2() appD.Lang2 {
	return t.Lang2
}

func (t *translation) GetTranslated() string {
	return t.Translated
}

func (t *translation) GetProvider() string {
	return t.Provider
}
