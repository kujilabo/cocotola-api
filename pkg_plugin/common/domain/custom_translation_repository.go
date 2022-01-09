package domain

import (
	"context"

	"github.com/go-playground/validator"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

// var ErrCustomTranslationNotFound = errors.New("azure translation not found")
// var ErrCustomTranslationAlreadyExists = errors.New("azure translation already exists")

type TranslationAddParameter interface {
	GetText() string
	GetPos() WordPos
	GetLang() app.Lang2
	GetTranslated() string
}

type translationAddParameter struct {
	Text       string `validate:"required"`
	Pos        WordPos
	Lang2      app.Lang2
	Translated string
}

func NewTransaltionAddParameter(text string, pos WordPos, lang app.Lang2, translated string) (TranslationAddParameter, error) {
	m := &translationAddParameter{
		Text:       text,
		Pos:        pos,
		Lang2:      lang,
		Translated: translated,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (p *translationAddParameter) GetText() string {
	return p.Text
}

func (p *translationAddParameter) GetPos() WordPos {
	return p.Pos
}

func (p *translationAddParameter) GetLang() app.Lang2 {
	return p.Lang2
}

func (p *translationAddParameter) GetTranslated() string {
	return p.Translated
}

type CustomTranslationRepository interface {
	Add(ctx context.Context, param TranslationAddParameter) (TranslationID, error)

	FindByText(ctx context.Context, text string, lang app.Lang2) ([]Translation, error)

	FindByFirstLetter(ctx context.Context, firstLetter string, lang app.Lang2) ([]Translation, error)

	Contain(ctx context.Context, text string, lang app.Lang2) (bool, error)
}
