package domain

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator"

	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

var ErrTatoebaSentenceAlreadyExists = errors.New("tatoebaSentence already exists")

type TatoebaSentenceAddParameter interface {
	GetSentenceNumber() int
	GetLang() appD.Lang3
	GetText() string
	GetAuthor() string
	GetUpdatedAt() time.Time
}

type tatoebaSentenceAddParameter struct {
	SentenceNumber int `validate:"required"`
	Lang           appD.Lang3
	Text           string `validate:"required"`
	Author         string `validate:"required"`
	UpdatedAt      time.Time
}

func NewTatoebaSentenceAddParameter(sentenceNumber int, lang appD.Lang3, text, author string, updatedAt time.Time) (TatoebaSentenceAddParameter, error) {
	m := &tatoebaSentenceAddParameter{
		SentenceNumber: sentenceNumber,
		Lang:           lang,
		Text:           text,
		Author:         author,
		UpdatedAt:      updatedAt,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (p *tatoebaSentenceAddParameter) GetSentenceNumber() int {
	return p.SentenceNumber
}

func (p *tatoebaSentenceAddParameter) GetLang() appD.Lang3 {
	return p.Lang
}

func (p *tatoebaSentenceAddParameter) GetText() string {
	return p.Text
}

func (p *tatoebaSentenceAddParameter) GetAuthor() string {
	return p.Author
}

func (p *tatoebaSentenceAddParameter) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

type TatoebaSentenceRepository interface {
	Add(ctx context.Context, param TatoebaSentenceAddParameter) error
}
