package domain

import (
	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

// type TranslationID uint

type Translation interface {
	// GetID() TranslationID
	// GetVersion() int
	// GetCreatedAt() time.Time
	// GetUpdatedAt() time.Time
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

func NewTranslation(
	// id TranslationID,
	// version int, createdAt time.Time, updatedAt time.Time,
	text string, pos WordPos, lang app.Lang2, translated, provider string) (Translation, error) {
	m := &translation{
		// ID:         id,
		// Version:    version,
		// CreatedAt:  createdAt,
		// UpdatedAt:  updatedAt,
		Text:       text,
		Pos:        pos,
		Lang2:      lang,
		Translated: translated,
		Provider:   provider,
	}

	return m, lib.Validator.Struct(m)
}

// func (t *translation) GetID() TranslationID {
// 	return t.ID
// }

// func (t *translation) GetVersion() int {
// 	return t.Version
// }

// func (t *translation) GetCreatedAt() time.Time {
// 	return t.CreatedAt
// }

// func (t *translation) GetUpdatedAt() time.Time {
// 	return t.UpdatedAt
// }

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
