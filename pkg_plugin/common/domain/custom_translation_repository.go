package domain

import (
	"context"
	"errors"

	"github.com/go-playground/validator"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

// var ErrCustomTranslationNotFound = errors.New("azure translation not found")
// var ErrCustomTranslationAlreadyExists = errors.New("azure translation already exists")

var ErrTranslationAlreadyExists = errors.New("custsomtranslation already exists")

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

func NewTransalationAddParameter(text string, pos WordPos, lang app.Lang2, translated string) (TranslationAddParameter, error) {
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

type TranslationUpdateParameter interface {
	GetTranslated() string
}

type translationUpdateParameter struct {
	Translated string `validate:"required"`
}

func NewTransaltionUpdateParameter(translated string) (TranslationUpdateParameter, error) {
	m := &translationUpdateParameter{
		Translated: translated,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (p *translationUpdateParameter) GetTranslated() string {
	return p.Translated
}

type CustomTranslationRepository interface {
	Add(ctx context.Context, param TranslationAddParameter) (TranslationID, error)

	Update(ctx context.Context, lang app.Lang2, text string, pos WordPos, param TranslationUpdateParameter) error

	Remove(ctx context.Context, lang app.Lang2, text string, pos WordPos) error

	FindByText(ctx context.Context, lang app.Lang2, text string) ([]Translation, error)

	FindByTextAndPos(ctx context.Context, lang app.Lang2, text string, pos WordPos) (Translation, error)

	FindByFirstLetter(ctx context.Context, lang app.Lang2, firstLetter string) ([]Translation, error)

	Contain(ctx context.Context, lang app.Lang2, text string) (bool, error)
}
