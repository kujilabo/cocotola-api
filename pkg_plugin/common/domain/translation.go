//go:generate mockery --output mock --name Translation
package domain

import (
	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

type Translation interface {
	GetText() string
	GetPos() WordPos
	GetLang() app.Lang2
	GetTranslated() string
	GetProvider() string
}

type translation struct {
	Text       string `validate:"required"`
	Pos        WordPos
	Lang2      app.Lang2
	Translated string
	Provider   string
}

func NewTranslation(text string, pos WordPos, lang app.Lang2, translated, provider string) (Translation, error) {
	m := &translation{
		Text:       text,
		Pos:        pos,
		Lang2:      lang,
		Translated: translated,
		Provider:   provider,
	}

	return m, lib.Validator.Struct(m)
}

func (t *translation) GetText() string {
	return t.Text
}

func (t *translation) GetPos() WordPos {
	return t.Pos
}

func (t *translation) GetLang() app.Lang2 {
	return t.Lang2
}

func (t *translation) GetTranslated() string {
	return t.Translated
}

func (t *translation) GetProvider() string {
	return t.Provider
}
