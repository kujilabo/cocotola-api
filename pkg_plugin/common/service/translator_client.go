//go:generate mockery --output mock --name TranslatorClient
package service

import (
	"context"
	"errors"

	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
)

var ErrTranslationNotFound = errors.New("translation not found")
var ErrTranslationAlreadyExists = errors.New("custsomtranslation already exists")

type TranslatorClient interface {
	DictionaryLookup(ctx context.Context, fromLang, toLang appD.Lang2, text string) ([]domain.Translation, error)
	DictionaryLookupWithPos(ctx context.Context, fromLang, toLang appD.Lang2, text string, pos domain.WordPos) (domain.Translation, error)
	FindTranslationsByFirstLetter(ctx context.Context, lang appD.Lang2, firstLetter string) ([]domain.Translation, error)
	FindTranslationByTextAndPos(ctx context.Context, lang appD.Lang2, text string, pos domain.WordPos) (domain.Translation, error)
	FindTranslationsByText(ctx context.Context, lang appD.Lang2, text string) ([]domain.Translation, error)
	AddTranslation(ctx context.Context, param TranslationAddParameter) error
	UpdateTranslation(ctx context.Context, lang appD.Lang2, text string, pos domain.WordPos, param TranslationUpdateParameter) error
	RemoveTranslation(ctx context.Context, lang appD.Lang2, text string, pos domain.WordPos) error
}

type TranslationAddParameter interface {
	GetText() string
	GetPos() domain.WordPos
	GetLang() appD.Lang2
	GetTranslated() string
}

type translationAddParameter struct {
	Text       string `validate:"required"`
	Pos        domain.WordPos
	Lang2      appD.Lang2
	Translated string
}

func NewTransalationAddParameter(text string, pos domain.WordPos, lang appD.Lang2, translated string) (TranslationAddParameter, error) {
	m := &translationAddParameter{
		Text:       text,
		Pos:        pos,
		Lang2:      lang,
		Translated: translated,
	}

	return m, libD.Validator.Struct(m)
}

func (p *translationAddParameter) GetText() string {
	return p.Text
}

func (p *translationAddParameter) GetPos() domain.WordPos {
	return p.Pos
}

func (p *translationAddParameter) GetLang() appD.Lang2 {
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

	return m, libD.Validator.Struct(m)
}

func (p *translationUpdateParameter) GetTranslated() string {
	return p.Translated
}
