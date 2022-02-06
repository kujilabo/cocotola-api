package domain

import (
	"time"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

type TatoebaSentence interface {
	GetSentenceNumber() int
	GetLang() app.Lang3
	GetText() string
	GetAuthor() string
	GetUpdatedAt() time.Time
}

type tatoebaSentence struct {
	SentenceNumber int
	Lang           app.Lang3
	Text           string
	Author         string
	UpdatedAt      time.Time
}

func NewTatoebaSentence(sentenceNumber int, lang app.Lang3, text, author string, updatedAt time.Time) (TatoebaSentence, error) {
	m := &tatoebaSentence{
		SentenceNumber: sentenceNumber,
		Lang:           lang,
		Text:           text,
		Author:         author,
		UpdatedAt:      updatedAt,
	}

	return m, libD.Validator.Struct(m)
}

func (m *tatoebaSentence) GetSentenceNumber() int {
	return m.SentenceNumber
}

func (m *tatoebaSentence) GetLang() app.Lang3 {
	return m.Lang
}

func (m *tatoebaSentence) GetText() string {
	return m.Text
}

func (m *tatoebaSentence) GetAuthor() string {
	return m.Author
}

func (m *tatoebaSentence) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

type TatoebaSentencePair interface {
	GetSrc() TatoebaSentence
	GetDst() TatoebaSentence
}

type tatoebaSentencePair struct {
	Src TatoebaSentence
	Dst TatoebaSentence
}

func NewTatoebaSentencePair(src, dst TatoebaSentence) (TatoebaSentencePair, error) {
	m := &tatoebaSentencePair{
		Src: src,
		Dst: dst,
	}

	return m, libD.Validator.Struct(m)
}

func (m *tatoebaSentencePair) GetSrc() TatoebaSentence {
	return m.Src
}

func (m *tatoebaSentencePair) GetDst() TatoebaSentence {
	return m.Dst
}
